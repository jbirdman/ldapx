// Package ldapx provides a higher level abstraction around github.com/go-ldap/ldap/v3
// that makes use of connection pooling.
package ldapx

import (
	"crypto/tls"
	"log"
	"net/url"

	"github.com/go-baa/pool"
	"github.com/go-ldap/ldap/v3"
	"github.com/jbirdman/ldapurl"
)

// Conn represents a connection to an LDAP server.
type Conn struct {
	ldapURL      *ldapurl.LdapURL // LDAP URL
	url          string
	pool         *pool.Pool //
	bindDN       string
	bindPassword string
	schema       *LDAPSchema
	tlsConfig    *tls.Config
	txConn       *ldap.Conn
}

// Client represents a client that can execute LDAP operations.
type Client interface {
	Execute(f func(*Conn) (interface{}, error)) (interface{}, error)
	ExecuteAs(dn string, password string, f func(*ldap.Conn) (interface{}, error)) (interface{}, error)
	Add(*ldap.AddRequest) error
	Del(*ldap.DelRequest) error
	CheckBind(dn string, password string) error
	Modify(*ldap.ModifyRequest) error
	Compare(dn string, attribute string, value string) (bool, error)
	PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	SearchWithPaging(request *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error)
	Lookup(dn string) (*Entry, error)
	LookupOrNew(dn string) (*Entry, error)
	QuickSearch(dn string, filter string, attributes []string) (*ldap.SearchResult, error)
	FindEntry(dn string, filter string, attributes []string) (*Entry, error)
	RootDSE() (*RootDSE, error)
	Schema() (*LDAPSchema, error)
	UpdateEntry(string, EntryUpdateFunc) error
	Update(entry *Entry) error
}

var _ Client = &Conn{}

// OpenURLSimple opens a connection to an LDAP server using the provided URL.
func OpenURLSimple(ldapURL, binddn, bindpw string, insecureSkipVerify bool) (*Conn, error) {
	// Get the host portion of the URL.
	host, err := urlHost(ldapURL)
	if err != nil {
		return nil, err
	}

	// If the host is an IP address, we need to disable hostname verification.
	return OpenURL(ldapURL, binddn, bindpw, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: insecureSkipVerify, //nolint: gosec
	})
}

// urlHost returns the host portion of a URL.
func urlHost(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

// OpenURL opens a connection to an LDAP server using the provided URL.
func OpenURL(url string, bindDN string, bindPassword string, tlsConfig *tls.Config) (*Conn, error) {
	// Parse the URL
	ldapURL, err := ldapurl.Parse(url)
	if err != nil {
		return nil, err
	}

	// Set up the connection pool.
	pl, err := setupConnectionPool(url, bindDN, bindPassword, tlsConfig)
	if err != nil {
		return nil, err
	}

	// Create the connection.
	conn := &Conn{
		ldapURL:      ldapURL,
		url:          url,
		pool:         pl,
		bindDN:       bindDN,
		bindPassword: bindPassword,
		tlsConfig:    tlsConfig,
	}

	// Get the schema.
	schema, err := conn.Schema()
	if err != nil {
		return nil, err
	}

	conn.schema = schema

	return conn, err
}

// setupConnectionPool sets up the connection pool.
func setupConnectionPool(url string, bindDN string, bindPassword string, tlsConfig *tls.Config) (*pool.Pool, error) {
	//
	pl, err := pool.New(1, 10, func() interface{} {
		// Dial the LDAP server.
		conn, err := dialURL(url, tlsConfig)
		if err != nil {
			log.Fatalf("create client connection error: %v\n", err)
		}

		// Bind to the LDAP server.
		// If the bind password is empty, attemp an unauthenticated bind
		if bindPassword != "" {
			err = conn.Bind(bindDN, bindPassword)
		} else {
			err = conn.UnauthenticatedBind(bindDN)
		}
		if err != nil {
			conn.Close()
			log.Fatalf("create client connection bind error: %v\n", err)
		}

		return conn
	})
	if err != nil {
		return nil, err
	}

	// Ping the connection.
	pl.Ping = func(conn interface{}) bool {

		return pingConnection(conn.(*ldap.Conn))
	}

	// Close the connection.
	pl.Close = func(conn interface{}) {
		conn.(*ldap.Conn).Close()
	}

	return pl, nil
}

func pingConnection(conn *ldap.Conn) bool {
	_, err := conn.Search(ldap.NewSearchRequest("", ldap.ScopeBaseObject, ldap.DerefAlways, 0, 0, false, "(objectclass=*)", nil, nil))
	if err != nil {
		return false
	}
	return true
}

// dialURL dials the LDAP server.
func dialURL(url string, tlsConfig *tls.Config) (*ldap.Conn, error) {
	return ldap.DialURL(url, ldap.DialWithTLSConfig(tlsConfig))
}

// getConn gets a connection from the pool.
func getConn(pool *pool.Pool) (*ldap.Conn, error) {
	lc, err := pool.Get()
	if err != nil {
		return nil, err
	}

	return lc.(*ldap.Conn), err
}

// putConn puts a connection back into the pool.
func putConn(pool *pool.Pool, lc *ldap.Conn) {
	pool.Put(lc)
}

// get gets a connection from the pool.
func (c *Conn) get() (*ldap.Conn, error) {
	if c.txConn != nil {
		return c.txConn, nil
	}
	return getConn(c.pool)
}

// put puts a connection back into the pool.
func (c *Conn) put(lc *ldap.Conn) {
	if c.txConn != nil && c.txConn == lc {
		return
	}
	putConn(c.pool, lc)
}

func (c *Conn) NewTx(conn *ldap.Conn) *Conn {
	return &Conn{
		txConn: conn,
	}
}

func (c *Conn) Execute(f func(*Conn) (interface{}, error)) (interface{}, error) {
	conn, err := c.get()
	if err != nil {
		return nil, err
	}
	defer c.put(conn)

	return f(c.NewTx(conn))
}

func (c *Conn) ExecuteLdap(f func(*ldap.Conn) (interface{}, error)) (interface{}, error) {
	conn, err := c.get()
	if err != nil {
		return nil, err
	}
	defer c.put(conn)

	return f(conn)
}

// ExecuteAs executes a function with a connection from the pool as a different user.
func (c *Conn) ExecuteAs(dn string, password string, f func(*ldap.Conn) (interface{}, error)) (interface{}, error) {
	conn, err := c.get()
	if err != nil {
		return nil, err
	}
	defer c.put(conn)

	err = conn.Bind(dn, password)
	if err != nil {
		return nil, err
	}
	defer func(c *Conn, conn *ldap.Conn) {
		_ = c.rebind(conn)
	}(c, conn)

	return f(conn)
}

// Search searches the LDAP server.
func (c *Conn) Search(request *ldap.SearchRequest) (*ldap.SearchResult, error) {
	result, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return conn.Search(request)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.SearchResult), nil
}

// SearchWithPaging searches the LDAP server with paging.
func (c *Conn) SearchWithPaging(request *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	result, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return conn.SearchWithPaging(request, pagingSize)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.SearchResult), nil
}

// Add adds an entry to the LDAP server.
func (c *Conn) Add(request *ldap.AddRequest) error {
	_, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Add(request)
	})
	return err
}

// Del deletes an entry from the LDAP server.
func (c *Conn) Del(request *ldap.DelRequest) error {
	_, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Del(request)
	})
	return err
}

// Modify modifies an entry on the LDAP server.
func (c *Conn) Modify(request *ldap.ModifyRequest) error {
	_, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Modify(request)
	})
	return err
}

// PasswordModify modifies a user's password on the LDAP server.
func (c *Conn) PasswordModify(request *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	result, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return conn.PasswordModify(request)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.PasswordModifyResult), nil
}

// Compare compares an attribute value on the LDAP server.
func (c *Conn) Compare(dn string, attribute string, value string) (bool, error) {
	result, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		return conn.Compare(dn, attribute, value)
	})
	return result.(bool), err
}

// CheckBind checks the bind credentials on the LDAP server.
func (c *Conn) CheckBind(dn string, password string) error {
	_, err := c.ExecuteLdap(func(conn *ldap.Conn) (interface{}, error) {
		err := conn.Bind(dn, password)
		defer func(c *Conn, conn *ldap.Conn) {
			_ = c.rebind(conn)
		}(c, conn)
		return nil, err
	})
	return err
}

// rebind rebinds to the LDAP server.
func (c *Conn) rebind(conn *ldap.Conn) error {
	return conn.Bind(c.bindDN, c.bindPassword)
}

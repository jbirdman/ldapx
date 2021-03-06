// Package provides a higher level abstration around gopkg.in/ldap.v2
// that makes use of connection pooling.
package ldapx

import (
	"crypto/tls"
	"git.jcu.edu.au/go/ldapurl"
	"github.com/go-baa/pool"
	"gopkg.in/ldap.v2"
	"log"
)

type Conn struct {
	ldapURL      *ldapurl.LdapURL
	pool         *pool.Pool
	bindDN       string
	bindPassword string
	tlsConfig    *tls.Config
}

type conn struct {
	conn         *ldap.Conn
	bindDN       string
	bindPassword string
}

type Client interface {
	Execute(f func(*ldap.Conn) (interface{}, error)) (interface{}, error)
	ExecuteAs(dn string, password string, f func(*ldap.Conn) (interface{}, error)) (interface{}, error)
	Add(*ldap.AddRequest) error
	Del(*ldap.DelRequest) error
	CheckBind(dn string, password string) error
	Modify(*ldap.ModifyRequest) error
	Compare(dn string, attribute string, value string) (bool, error)
	PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error)
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	SearchWithPaging(request *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error)
}

var _ Client = &Conn{}

func OpenURL(url string, bindDN string, bindPassword string, tlsConfig *tls.Config) (*Conn, error) {
	ldapUrl, err := ldapurl.Parse(url)
	if err != nil {
		return nil, err
	}

	pl, err := setupConnectionPool(ldapUrl, bindDN, bindPassword, tlsConfig)

	return &Conn{
		ldapURL:      ldapUrl,
		pool:         pl,
		bindDN:       bindDN,
		bindPassword: bindPassword,
		tlsConfig:    tlsConfig,
	}, err
}

func setupConnectionPool(ldapUrl *ldapurl.LdapURL, bindDN string, bindPassword string, tlsConfig *tls.Config) (*pool.Pool, error) {
	pl, err := pool.New(1, 10, func() interface{} {
		conn, err := dialUrl(ldapUrl, tlsConfig)
		if err != nil {
			log.Fatalf("create client connection error: %v\n", err)
		}

		err = conn.Bind(bindDN, bindPassword)
		if err != nil {
			conn.Close()
			log.Fatalf("create client connection bind error: %v\n", err)
		}

		return conn
	})
	if err != nil {
		return nil, err
	}

	pl.Ping = func(conn interface{}) bool {
		return true
	}

	pl.Close = func(conn interface{}) {
		conn.(*ldap.Conn).Close()
	}

	return pl, nil
}

func dialUrl(ldapURL *ldapurl.LdapURL, tlsConfig *tls.Config) (*ldap.Conn, error) {
	hostname := ldapURL.BuildHostnamePortString()
	var l *ldap.Conn
	var err error

	if ldapURL.IsTLS() {
		l, err = ldap.DialTLS("tcp", hostname, tlsConfig)
	} else {
		l, err = ldap.Dial("tcp", hostname)
	}

	return l, err
}

func (c *Conn) getAs(dn string, password string) (*conn, error) {
	panic("Unimplemeted")
}

func (c *Conn) get() (*ldap.Conn, error) {
	lc, err := c.pool.Get()
	if err != nil {
		return nil, err
	}

	return lc.(*ldap.Conn), err
}

func (c *Conn) mustGet() *ldap.Conn {
	lc, err := c.get()
	if err != nil {
		panic(err)
	}
	return lc
}

func (c *Conn) put(lc *ldap.Conn) {
	c.pool.Put(lc)
}

func (c *Conn) Execute(f func(*ldap.Conn) (interface{}, error)) (interface{}, error) {
	conn, err := c.get()
	if err != nil {
		return nil, err
	}
	defer c.put(conn)

	return f(conn)
}

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
	//noinspection GoUnhandledErrorResult
	defer c.rebind(conn)

	return f(conn)
}

func (c *Conn) Search(request *ldap.SearchRequest) (*ldap.SearchResult, error) {
	result, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return conn.Search(request)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.SearchResult), nil
}

func (c *Conn) SearchWithPaging(request *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	result, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return conn.SearchWithPaging(request, pagingSize)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.SearchResult), nil
}

func (c *Conn) Add(request *ldap.AddRequest) error {
	_, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Add(request)
	})
	return err
}

func (c *Conn) Del(request *ldap.DelRequest) error {
	_, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Del(request)
	})
	return err
}

func (c *Conn) Modify(request *ldap.ModifyRequest) error {
	_, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return nil, conn.Modify(request)
	})
	return err
}

func (c *Conn) PasswordModify(request *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	result, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return conn.PasswordModify(request)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ldap.PasswordModifyResult), nil
}

func (c *Conn) Compare(dn string, attribute string, value string) (bool, error) {
	result, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		return conn.Compare(dn, attribute, value)
	})
	return result.(bool), err
}

func (c *Conn) CheckBind(dn string, password string) error {
	_, err := c.Execute(func(conn *ldap.Conn) (interface{}, error) {
		err := conn.Bind(dn, password)
		//noinspection GoUnhandledErrorResult
		defer c.rebind(conn)
		return nil, err
	})
	return err
}

func (c *Conn) rebind(conn *ldap.Conn) error {
	return conn.Bind(c.bindDN, c.bindPassword)
}

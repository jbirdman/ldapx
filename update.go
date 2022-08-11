package ldapx

type EntryUpdateFunc func(*Entry) (*Entry, error)

func (c *Conn) UpdateEntry(dn string, f EntryUpdateFunc) error {
	entry, err := c.LookupOrNew(dn)
	if err != nil {
		return err
	}
	entry, err = f(entry)
	if err != nil {
		return err
	}

	if entry.Changed() {
		return entry.Update(c)
	}
	return nil
}

func (c *Conn) Update(entry *Entry) error {
	return entry.Update(c)
}

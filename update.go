package ldapx

// EntryUpdateFunc is a function that can be used to update an entry.
type EntryUpdateFunc func(*Entry) (*Entry, error)

// UpdateEntry updates the entry with the given DN using the given function.
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

// Update updates the given entry.
func (c *Conn) Update(entry *Entry) error {
	return entry.Update(c)
}

package osdlite

func (o OSDLite) initPragmas() error {
	query := `
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = 1;	
	`
	if _, err := o.fs.NewQuery(query).Execute(); err != nil {
		return err
	}
	return nil
}

func (o OSDLite) initOptions() {
	o.fs.DB().SetMaxOpenConns(1)
}

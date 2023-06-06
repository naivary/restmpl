package osdlite

func (o OSDLite) initSchema() error {
	schemas := []func() error{
		o.initBucketSchema,
		o.initObjectSchema,
		o.initNameUniquenessWithinBucket,
	}
	for _, schema := range schemas {
		if err := schema(); err != nil {
			return err
		}
	}
	return nil
}

func (o OSDLite) initBucketSchema() error {
	q := o.store.CreateTable("buckets", map[string]string{
		"id":         "TEXT PRIMARY KEY",
		"created_at": "INTEGER",
		"name":       "TEXT UNIQUE",
		"owner":      "TEXT",
		"basepath":   "TEXT",
	})
	if _, err := q.Execute(); err != nil {
		return err
	}
	return nil
}

func (o OSDLite) initObjectSchema() error {
	q := o.store.CreateTable("objects", map[string]string{
		"id":            "TEXT PRIMARY KEY",
		"created_at":    "INTEGER",
		"last_modified": "INTEGER",
		"name":          "TEXT",
		"owner":         "TEXT",
		"tags":          "TEXT",
		"version":       "INTEGER",
		"bucket_id":     "TEXT REFERENCES buckets(id) ON DELETE CASCADE",
		"payload":       "BLOB",
	})
	if _, err := q.Execute(); err != nil {
		return err
	}
	return nil
}

func (o OSDLite) initNameUniquenessWithinBucket() error {
	q := o.store.CreateUniqueIndex("objects", "composite_uniqueness", "name", "bucket_id")
	if _, err := q.Execute(); err != nil {
		return err
	}
	return nil
}

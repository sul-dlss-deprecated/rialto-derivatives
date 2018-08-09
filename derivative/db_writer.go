package derivative

// DbWriter writes a derivative document
type DbWriter interface {
	RemoveAll() error
	Add(docs []interface{}) error
}

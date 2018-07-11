package derivative

import "github.com/vanng822/go-solr/solr"

// Writer writes a derivative document
type Writer interface {
	RemoveAll() error
	Add(docs []solr.Document) error
}

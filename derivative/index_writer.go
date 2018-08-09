package derivative

import "github.com/vanng822/go-solr/solr"

// IndexWriter writes one or more Solr documents to a search index
type IndexWriter interface {
	RemoveAll() error
	Add(docs []solr.Document) error
}

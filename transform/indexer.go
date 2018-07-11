package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

type Indexer struct{}

func (m *Indexer) Map(resources []models.Resource) []solr.Document {
	docs := make([]solr.Document, len(resources))
	for i, v := range resources {
		docs[i] = m.mapOne(v)
	}
	return docs
}

func (m *Indexer) mapOne(resource models.Resource) solr.Document {
	doc := make(solr.Document)
	doc.Set("type_ssi", resource.ResourceTypes()[0].String())
	if resource.Titles() != nil {
		doc.Set("title_ssi", resource.Titles()[0].String())
	}
	return doc
}

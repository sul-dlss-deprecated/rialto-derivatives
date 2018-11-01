package transform

import (
	"fmt"

	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

// CompositeIndexer delegates to subindexers to transform resources to solr Documents
type CompositeIndexer struct {
	conceptIndexer      *ConceptIndexer
	grantIndexer        *GrantIndexer
	organizationIndexer *OrganizationIndexer
	personIndexer       *PersonIndexer
	projectIndexer      *ProjectIndexer
	publicationIndexer  *PublicationIndexer
}

// NewCompositeIndexer creates a new CompositeIndexer instance
func NewCompositeIndexer(repository repository.Repository) *CompositeIndexer {
	return &CompositeIndexer{
		conceptIndexer:      &ConceptIndexer{},
		grantIndexer:        &GrantIndexer{},
		organizationIndexer: &OrganizationIndexer{},
		personIndexer:       NewPersonIndexer(repository),
		projectIndexer:      &ProjectIndexer{},
		publicationIndexer:  &PublicationIndexer{},
	}
}

// Map transforms a collection of resources into a collection of Solr Documents
func (m *CompositeIndexer) Map(resources []models.Resource) []solr.Document {
	docs := make([]solr.Document, len(resources))
	for i, v := range resources {
		if v == nil {
			continue
		}
		docs[i] = m.mapOne(v)
	}
	return docs
}

// mapOne sets the id and type and then delegates to the type specific indexer
func (m *CompositeIndexer) mapOne(resource models.Resource) solr.Document {
	doc := make(solr.Document)
	doc.Set("id", resource.Subject())
	switch v := resource.(type) {
	case *models.Publication:
		doc = m.publicationIndexer.Index(v, doc)
	case *models.Person:
		doc = m.personIndexer.Index(v, doc)
	case *models.Organization:
		doc = m.organizationIndexer.Index(v, doc)
	case *models.Grant:
		doc = m.grantIndexer.Index(v, doc)
	case *models.Project:
		doc = m.projectIndexer.Index(v, doc)
	case *models.Concept:
		doc = m.conceptIndexer.Index(v, doc)
	default:
		panic(fmt.Errorf("No indexer for %v", v))
	}

	return doc
}

package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/repository"
	"github.com/vanng822/go-solr/solr"
)

// CompositeIndexer delegates to subindexers to transform resources to solr Documents
type CompositeIndexer struct {
	conceptIndexer      Indexer
	defaultIndexer      Indexer
	grantIndexer        Indexer
	organizationIndexer Indexer
	personIndexer       Indexer
	projectIndexer      Indexer
	publicationIndexer  Indexer
	typeIndexer         Indexer
}

// NewCompositeIndexer creates a new CompositeIndexer instance
func NewCompositeIndexer(repository repository.Repository) *CompositeIndexer {
	return &CompositeIndexer{
		conceptIndexer:      &ConceptIndexer{},
		defaultIndexer:      &DefaultIndexer{},
		grantIndexer:        &GrantIndexer{},
		organizationIndexer: &OrganizationIndexer{},
		personIndexer:       NewPersonIndexer(repository),
		projectIndexer:      &ProjectIndexer{},
		publicationIndexer:  &PublicationIndexer{},
		typeIndexer:         &TypeIndexer{},
	}
}

// Map transforms a collection of resources into a collection of Solr Documents
func (m *CompositeIndexer) Map(resources []models.Resource) []solr.Document {
	docs := make([]solr.Document, len(resources))
	for i, v := range resources {
		docs[i] = m.mapOne(v)
	}
	return docs
}

// mapOne sets the id and type and then delegates to the type specific indexer
func (m *CompositeIndexer) mapOne(resource models.Resource) solr.Document {
	doc := make(solr.Document)
	doc.Set("id", resource.Subject())

	doc = m.typeIndexer.Index(resource, doc)

	var indexer Indexer
	if resource.IsPublication() {
		indexer = m.publicationIndexer
	} else if resource.IsPerson() {
		indexer = m.personIndexer
	} else if resource.IsOrganization() {
		indexer = m.organizationIndexer
	} else if resource.IsGrant() {
		indexer = m.grantIndexer
	} else if resource.IsProject() {
		indexer = m.projectIndexer
	} else if resource.IsConcept() {
		indexer = m.conceptIndexer
	} else {
		indexer = m.defaultIndexer
	}
	return indexer.Index(resource, doc)
}

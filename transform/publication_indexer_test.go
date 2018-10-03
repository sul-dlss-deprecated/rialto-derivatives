package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

func TestPublicationResourceToDoc(t *testing.T) {
	indexer := &PublicationIndexer{}
	doi := "https://doi.org/10.1109/5.771073"
	author1 := &models.Author{URI: "http://example.com/person1", Label: "Harry Potter"}
	author2 := &models.Author{URI: "http://example.com/person2", Label: "Hermione Granger"}

	authors := []*models.Author{author1, author2}
	resource := &models.Publication{DOI: &doi, Title: "Hello world!", Authors: authors}
	in := make(solr.Document)
	in.Set("id", "http://example.com/record1")
	doc := indexer.Index(resource, in)

	assert.Equal(t, "Hello world!", doc.Get("title_tesi"))
	assert.Equal(t, "http://example.com/record1", doc.Get("id"))
	assert.Equal(t, "https://doi.org/10.1109/5.771073", doc.Get("doi_ssim"))
	assert.Equal(t, []string{"http://example.com/person1", "http://example.com/person2"}, doc.Get("authors_ssim"))
	assert.Equal(t, []string{"Harry Potter", "Hermione Granger"}, doc.Get("author_labels_tsim"))
}

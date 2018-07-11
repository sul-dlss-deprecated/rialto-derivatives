package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestResourceToDoc(t *testing.T) {
	indexer := &Indexer{}
	data := make(map[string]interface{})
	data["rdf:type"] = "http://purl.org/ontology/bibo/Document"
	data["dc:title"] = "Hello world!"

	resource := models.NewResource(data)
	resourceList := []models.Resource{resource}
	docs := indexer.Map(resourceList)

	assert.Equal(t, "Hello world!", docs[0].Get("title_ssi"))
}

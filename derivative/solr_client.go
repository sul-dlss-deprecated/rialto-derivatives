package derivative

import (
	"fmt"
	"net/url"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/sul-dlss-labs/rialto-derivatives/transform"
	"github.com/vanng822/go-solr/solr"
)

// SolrClient represents the functions on the Solr index
type SolrClient struct {
	si      *solr.SolrInterface
	indexer *transform.CompositeIndexer
}

// NewSolrClient returns a new SolrClient instance
func NewSolrClient(host string, collection string, indexer *transform.CompositeIndexer) *SolrClient {
	si, _ := solr.NewSolrInterface(host, collection)
	return &SolrClient{
		si:      si,
		indexer: indexer,
	}
}

// RemoveResourcesOfType clears the index of all the data with the matching type
func (d *SolrClient) RemoveResourcesOfType(resourceType string) error {
	query := fmt.Sprintf("%s:%s", transform.SolrTypeField, resourceType)
	data := map[string]interface{}{"query": query}
	_, err := d.si.Delete(data, nil)
	return err
}

// RemoveAll clears the index of all the data
func (d *SolrClient) RemoveAll() error {
	_, err := d.si.DeleteAll()
	return err
}

// Add puts a bunch of documents in Solr
func (d *SolrClient) Add(resourceList []models.Resource) error {
	docs := d.indexer.Map(resourceList)
	values := &url.Values{
		"softCommit": []string{"true"},
	}
	_, err := d.si.Add(docs, 1000, values)
	return err
}

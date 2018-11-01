package derivative

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/transform"
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
	query := fmt.Sprintf("%s:%s", "type_ssim", resourceType)
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
	log.Printf("[SOLR] Add invoked with %v docs", len(resourceList))
	docs := d.indexer.Map(resourceList)
	values := &url.Values{
		"softCommit": []string{"true"},
	}
	start := time.Now()
	_, err := d.si.Add(docs, 1000, values)
	elapsed := time.Since(start).Round(time.Microsecond * 100)
	log.Printf("[SOLR] Added %v docs in %s", len(resourceList), elapsed)
	return err
}

package derivative

import (
	"fmt"
	"net/url"

	"github.com/vanng822/go-solr/solr"
)

// SolrClient represents the functions on the Solr index
type SolrClient struct {
	si *solr.SolrInterface
}

const typeField = "type_ssi"

// NewSolrClient returns a new SolrClient instance
func NewSolrClient(host string, collection string) *SolrClient {
	si, _ := solr.NewSolrInterface(host, collection)
	return &SolrClient{
		si: si,
	}
}

func (d *SolrClient) removeResourcesOfType(resourceType string) error {
	query := fmt.Sprintf("%s:%s", typeField, resourceType)
	data := map[string]interface{}{"query": query}
	_, err := d.si.Delete(data, nil)
	return err
}

// Add puts a bunch of documents in Solr
func (d *SolrClient) Add(docs []solr.Document) error {
	values := &url.Values{
		"softCommit": []string{"true"},
	}
	_, err := d.si.Add(docs, 1000, values)
	return err
}

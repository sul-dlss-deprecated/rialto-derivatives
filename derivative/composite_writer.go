package derivative

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// CompositeWriter writes to both solr and postgres
type CompositeWriter struct {
	pg   Writer
	solr Writer
}

// NewCompositeWriter returns a new CompositeWriter instance
func NewCompositeWriter(pg Writer, solr Writer) *CompositeWriter {
	return &CompositeWriter{
		pg:   pg,
		solr: solr,
	}
}

// RemoveAll clears the index of all the data
func (d *CompositeWriter) RemoveAll() error {
	if err := d.pg.RemoveAll(); err != nil {
		return err
	}

	return d.solr.RemoveAll()
}

// Add puts a bunch of documents in PostgreSql & Solr
func (d *CompositeWriter) Add(resourceList []models.Resource) error {
	if err := d.pg.Add(resourceList); err != nil {
		return err
	}

	return d.solr.Add(resourceList)
}

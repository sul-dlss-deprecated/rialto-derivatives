package transform

import (
	"fmt"
	"time"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// ConceptSerializer transforms concept resource types into JSON Documents
type ConceptSerializer struct {
}

// NewConceptSerializer makes a new instance of the ConceptSerializer
func NewConceptSerializer() *ConceptSerializer {
	return &ConceptSerializer{}
}

// Serialize returns the Concept resource as a JSON string.
func (m *ConceptSerializer) Serialize(resource *models.Concept) string {
	return "{}"
}

// SQLForInsert returns the sql and the values to insert
func (m *ConceptSerializer) SQLForInsert(resource *models.Concept) (string, []interface{}) {
	table := "concepts"
	name := resource.Name()
	data := m.Serialize(resource)
	subject := resource.Subject()
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE %v.uri=$1`, table, table)
	vals := []interface{}{subject, name, data, time.Now(), time.Now()}
	return sql, vals
}

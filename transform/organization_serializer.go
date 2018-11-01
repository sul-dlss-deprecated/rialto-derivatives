package transform

import (
	"fmt"
	"time"

	"github.com/sul-dlss/rialto-derivatives/models"
)

// OrganizationSerializer transforms organization resource types into JSON Documents
type OrganizationSerializer struct {
}

// Serialize returns the Organization resource as a JSON string..
// Must include the following properties:
//
//   name (string)
//   type (URI) the most specific type (e.g. Department or University)
func (m *OrganizationSerializer) Serialize(resource *models.Organization) string {
	return fmt.Sprintf(`{"type": "%s"}`, resource.Subtype)
}

// SQLForInsert returns the sql and the values to insert
func (m *OrganizationSerializer) SQLForInsert(resource *models.Organization) (string, []interface{}) {
	table := "organizations"
	name := resource.Name
	data := m.Serialize(resource)
	subject := resource.Subject()
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE %v.uri=$1`, table, table)
	vals := []interface{}{subject, name, data, time.Now(), time.Now()}
	return sql, vals
}

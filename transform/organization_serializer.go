package transform

import (
	"fmt"

	"github.com/sul-dlss-labs/rialto-derivatives/models"
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
	return fmt.Sprintf(`{"name": "%s", "type": "%s"}`, resource.Name, resource.Subtype)
}

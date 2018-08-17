package transform

import (
	"fmt"

	"github.com/knakk/rdf"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

// OrganizationSerializer transforms organization resource types into JSON Documents
type OrganizationSerializer struct {
}

const agent = "http://xmlns.com/foaf/0.1/Agent"
const organization = "http://xmlns.com/foaf/0.1/Organization"

// Serialize returns the Organization resource as a JSON string..
// Must include the following properties:
//
//   name (string)
//   type (URI) the most specific type (e.g. Department or University)
func (m *OrganizationSerializer) Serialize(resource models.Resource) string {
	name := resource.ValueOf("orgName")[0]
	return fmt.Sprintf(`{"name": "%s", "type": "%s"}`, name, m.bestType(resource.ValueOf("type")))
}

// Return the most specific type of resource (e.g. not Agent or Organization)
func (m *OrganizationSerializer) bestType(types []rdf.Term) string {
	if len(types) == 0 {
		return ""
	}
	bestType := types[0].String()
	for _, thisType := range types {
		try := thisType.String()
		if try != agent && (try != organization || bestType == agent) {
			bestType = try
		}
	}
	return bestType
}

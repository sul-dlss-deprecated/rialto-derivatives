package models

// Resource represents the data we get from Neptune
type Resource struct {
	data map[string]interface{}
}

const typePredicate = "rdf:type"
const titlePredicate = "dc:title"

var publicationTypes = []string{
	"http://purl.org/ontology/bibo/Document",
}

// NewResource creates a new instance of the resource
func NewResource(data map[string]interface{}) Resource {
	return Resource{data: data}
}

// IsPublication returns true if the type is a publiction
func (r *Resource) IsPublication() bool {
	a := r.ResourceType()
	for _, b := range publicationTypes {
		if b == a {
			return true
		}
	}
	return false
}

// ResourceType returns the type assertion
func (r *Resource) ResourceType() string {
	return r.data[typePredicate].(string)
}

// Title returns the title assertion
func (r *Resource) Title() string {
	return r.data[titlePredicate].(string)
}

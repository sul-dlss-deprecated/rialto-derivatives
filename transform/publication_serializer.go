package transform

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sul-dlss/rialto-derivatives/models"
)

// PublicationSerializer transforms publication resource types into JSON Documents
type PublicationSerializer struct {
}

type publication struct {
	Title       string    `json:"title"`
	CreatedYear *int      `json:"created_year"`
	Concepts    *[]string `json:"concepts"`
}

// NewPublicationSerializer makes a new instance of the PersonSerializer
func NewPublicationSerializer() *PublicationSerializer {
	return &PublicationSerializer{}
}

// Serialize returns the Publication resource as a JSON string.
func (m *PublicationSerializer) Serialize(resource *models.Publication) string {
	p := &publication{
		Title:    resource.Title,
		Concepts: m.retrieveURIs(resource.Concepts),
	}

	if resource.CreatedYear != 0 {
		p.CreatedYear = &resource.CreatedYear
	}

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// SQLForInsert returns the sql and the values to insert
func (m *PublicationSerializer) SQLForInsert(resource *models.Publication) (string, []interface{}) {
	table := "publications"
	data := m.Serialize(resource)
	subject := resource.Subject()
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (uri) DO UPDATE SET metadata=$2, updated_at=$4 WHERE %v.uri=$1`, table, table)
	vals := []interface{}{subject, data, time.Now(), time.Now()}
	return sql, vals
}

// ShouldAdd returns True if this publication should be added.
// Publications should not be added if they will not be included in a report.
func (m *PublicationSerializer) ShouldAdd(resource *models.Publication) bool {
	if resource.CreatedYear < 2000 {
		return false
	}
	return true
}

func (m *PublicationSerializer) retrieveURIs(resources []*models.Concept) *[]string {
	uris := make([]string, len(resources))
	for n, resource := range resources {
		uris[n] = resource.URI
	}
	return &uris
}

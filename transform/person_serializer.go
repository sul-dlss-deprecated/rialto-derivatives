package transform

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-derivatives/repository"
)

// PersonSerializer transforms person resource types into JSON Documents
type PersonSerializer struct {
	repo repository.Repository
}

type person struct {
	Departments       *[]string `json:"departments"`
	DepartmentLabels  *[]string `json:"department_labels"`
	Schools           *[]string `json:"schools"`
	SchoolLabels      *[]string `json:"school_labels"`
	Institutions      *[]string `json:"institutions"`
	InstitutionLabels *[]string `json:"institution_labels"`
	Institutes        *[]string `json:"institutes"`
	InstituteLabels   *[]string `json:"institute_labels"`
	Countries         *[]string `json:"country_labels"`
}

// NewPersonSerializer makes a new instance of the PersonSerializer
func NewPersonSerializer(repo repository.Repository) *PersonSerializer {
	return &PersonSerializer{repo: repo}
}

// Serialize returns the Person resource as a JSON string.
func (m *PersonSerializer) Serialize(resource *models.Person) string {
	p := &person{
		Departments:       m.retrieveURIs(resource.DepartmentOrgs),
		DepartmentLabels:  m.retrieveLabels(resource.DepartmentOrgs),
		Institutions:      m.retrieveURIs(resource.InstitutionOrgs),
		InstitutionLabels: m.retrieveLabels(resource.InstitutionOrgs),
		Institutes:        m.retrieveURIs(resource.InstituteOrgs),
		InstituteLabels:   m.retrieveLabels(resource.InstituteOrgs),
		Schools:           m.retrieveURIs(resource.SchoolOrgs),
		SchoolLabels:      m.retrieveLabels(resource.SchoolOrgs),
		Countries:         m.retrieveLabels(resource.Countries),
	}

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// SQLForInsert returns the sql and the values to insert
func (m *PersonSerializer) SQLForInsert(resource *models.Person) (string, []interface{}) {
	table := "people"
	name := resource.Name()
	data := m.Serialize(resource)
	subject := resource.Subject()
	sql := fmt.Sprintf(`INSERT INTO "%v" ("uri", "name", "metadata", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uri) DO UPDATE SET name=$2, metadata=$3, updated_at=$5 WHERE %v.uri=$1`, table, table)
	vals := []interface{}{subject, name, data, time.Now(), time.Now()}
	return sql, vals
}

func (m *PersonSerializer) retrieveURIs(resources []*models.Labeled) *[]string {
	uris := make([]string, len(resources))
	for n, resource := range resources {
		uris[n] = resource.URI
	}
	return &uris
}

func (m *PersonSerializer) retrieveLabels(resources []*models.Labeled) *[]string {
	labels := make([]string, len(resources))
	for n, resource := range resources {
		labels[n] = resource.Label
	}
	return &labels
}

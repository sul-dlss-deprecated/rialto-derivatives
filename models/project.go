package models

import "github.com/knakk/rdf"

// Project represents research projects that are funded by grants or institutions,
// worked on by agents, and may produce multiple works.
type Project struct {
	URI              string
	Title            string
	AlternativeTitle string
	StartDate        string
	EndDate          string
}

// NewProject instantiates a project from sparql results
func NewProject(data map[string]rdf.Term) *Project {
	project := &Project{
		URI:   data["id"].String(),
		Title: data["title"].String(),
	}

	if alternativetitle := data["alternativetitle"]; alternativetitle != nil {
		project.AlternativeTitle = alternativetitle.String()
	}

	if startdate := data["startdate"]; startdate != nil {
		project.StartDate = startdate.String()
	}

	if enddate := data["enddate"]; enddate != nil {
		project.EndDate = enddate.String()
	}
	return project
}

// Subject returns the resources Subject URI
func (c *Project) Subject() string {
	return c.URI
}

package models

import (
	"strconv"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
)

// Publication is a representation of articles, research outputs, datasets, etc.
// If feasible, there should be a link to manifestations of that Work (i.e. DOI).
type Publication struct {
	URI         string
	Subtype     *string
	Title       string
	DOI         *string
	Abstract    *string
	Identifier  string
	Publisher   *string
	Description *string
	Created     string
	CreatedYear int
	Authors     []*Author
	Grants      []*Labeled
	Concepts    []*Concept
}

// Author is a representation of a person that authored a publication.
type Author Labeled

// NewPublication instantiates a publication from sparql results
func NewPublication(data map[string]rdf.Term) *Publication {
	pub := &Publication{
		URI:        data["id"].String(),
		Title:      data["title"].String(),
		Authors:    []*Author{},
		Grants:     []*Labeled{},
		Created:    data["created"].String(),
		Identifier: data["identifier"].String(),
	}

	if subtype := data["subtype"]; subtype != nil {
		subtypeStr := subtype.String()
		pub.Subtype = &subtypeStr
	}

	if doi := data["doi"]; doi != nil {
		doiStr := doi.String()
		pub.DOI = &doiStr
	}

	if abstract := data["abstract"]; abstract != nil {
		abstractStr := abstract.String()
		pub.Abstract = &abstractStr
	}

	if publisher := data["publisher"]; publisher != nil {
		publisherStr := publisher.String()
		pub.Publisher = &publisherStr
	}

	if description := data["description"]; description != nil {
		descriptionStr := description.String()
		pub.Description = &descriptionStr
	}

	if date := data["created"]; date != nil {
		yearStr := date.String()[0:4]
		yearInt, err := strconv.Atoi(yearStr)
		if err == nil {
			pub.CreatedYear = yearInt
		}
	}
	return pub
}

// Name returns the resources Name
func (c *Publication) Name() string {
	return c.Title
}

// SetAuthorInfo allow author relationships to be passed in
func (c *Publication) SetAuthorInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		uri := solution["id"].String()
		label := solution["label"].String()
		c.Authors = append(c.Authors, &Author{URI: uri, Label: label})
	}
}

// SetGrantInfo allow grant relationships to be passed in
func (c *Publication) SetGrantInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		uri := solution["id"].String()
		label := solution["label"].String()
		c.Grants = append(c.Grants, &Labeled{URI: uri, Label: label})
	}
}

// SetConceptInfo allow concept relationships to be passed in
func (c *Publication) SetConceptInfo(results *sparql.Results) {
	solutions := results.Solutions()
	for _, solution := range solutions {
		c.Concepts = append(c.Concepts, NewConcept(solution))
	}
}

// Subject returns the resources Subject URI
func (c *Publication) Subject() string {
	return c.URI
}

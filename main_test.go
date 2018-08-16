package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/rialto-derivatives/actions"
	"github.com/sul-dlss-labs/rialto-derivatives/message"
	"github.com/sul-dlss-labs/rialto-derivatives/models"
)

func TestMain(m *testing.M) {
	beforeAll()
	os.Exit(m.Run())
}

func beforeAll() {
	err := buildServiceRegistry().DbWriter.RemoveAll()
	if err != nil {
		panic(err)
	}
}

func TestIntegration(t *testing.T) {
	client := buildServiceRegistry().DbWriter
	organization := models.Organization{
		URI: "http://example.org/organization/1",
		Metadata: models.OrganizationMetadata{
			Department: "Biochemistry",
			Name:       "School of the Sciences",
			Type:       "School",
			Country:    "USA",
		},
	}
	err := client.AddOrganization(&organization)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "Biochemistry", organization.Metadata.Department)

	initialCreateDate := organization.CreatedAt
	initialUpdateDate := organization.UpdatedAt

	updatedOrg := models.Organization{
		URI: "http://example.org/organization/1",
		Metadata: models.OrganizationMetadata{
			Department: "BioX",
			Type:       "Institute",
		},
	}

	err = client.UpdateOrganization(&updatedOrg)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, initialCreateDate, updatedOrg.CreatedAt)
	assert.NotEqual(t, initialUpdateDate, updatedOrg.UpdatedAt)
	assert.Equal(t, "BioX", updatedOrg.Metadata.Department)
	assert.Equal(t, "Institute", updatedOrg.Metadata.Type)
}

func TestTouchAction(t *testing.T) {
	msg := &message.Message{
		Action: "touch",
	}
	action := actionForMessage(msg, nil)

	assert.IsType(t, &actions.TouchAction{}, action)
}

func TestRebuildAction(t *testing.T) {
	msg := &message.Message{
		Action: "rebuild",
	}
	action := actionForMessage(msg, nil)

	assert.IsType(t, &actions.RebuildAction{}, action)
}

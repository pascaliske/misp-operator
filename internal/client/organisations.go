package client

import (
	"fmt"
	"net/http"
)

type Organisation struct {
	Id           string `json:"id"`
	Uuid         string `json:"uuid"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type organisationEnvelope struct {
	Organisation Organisation `json:"Organisation"`
}

type OrganisationCreatePayload struct {
	Uuid        string `json:"uuid,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (o Organisation) String() string {
	return fmt.Sprintf("organisation %q", o.Name)
}

func (o Organisation) CreatePayload() OrganisationCreatePayload {
	return OrganisationCreatePayload{
		Uuid:        o.Uuid,
		Name:        o.Name,
		Type:        o.Type,
		Description: o.Description,
	}
}

func (c *ApiClient) GetOrganisations() (*[]Organisation, error) {
	return getResourceList(c, "/organisations", func(e organisationEnvelope) Organisation {
		return e.Organisation
	})
}

func (c *ApiClient) GetOrganisationByName(name string) (*Organisation, error) {
	// get all organisations
	organisations, err := c.GetOrganisations()
	if err != nil {
		return nil, err
	}

	// find first by name
	organisation, found := findFirst(*organisations, func(o Organisation) bool {
		return o.Name == name
	})
	if !found {
		return nil, &ApiError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("organisation %q not found", name)}
	}

	return organisation, nil
}

func (c *ApiClient) CreateOrganisation(organisation *Organisation) (*Organisation, error) {
	return create(c, "/admin/organisations/add", *organisation, func(e organisationEnvelope) Organisation {
		return e.Organisation
	})
}

func (c *ApiClient) EnsureOrganisation(organisation *Organisation) (*Organisation, error) {
	return ensure(*organisation, func() (*Organisation, error) {
		return c.GetOrganisationByName(organisation.Name)
	}, func() (*Organisation, error) {
		return c.CreateOrganisation(organisation)
	})
}

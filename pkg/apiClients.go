package gospeckle

import (
	"context"
	"net/http"
)

const apiClientBasePath = "clients"

// APIClient is the request response when fetching apiClients
type APIClient struct {
	Metadata         Metadata
	Role             string `json:"role,omitempty"`
	DocumentName     string `json:"documentName,omitempty"`
	DocumentType     string `json:"documentType,omitempty"`
	DocumentLocation string `json:"documentLocation,omitempty"`
	DocumentGUID     string `json:"documentGuid,omitempty"`
	StreamID         string `json:"streamId,omitempty"`
	Online           bool   `json:"online,omitempty"`
	// Owner            Account `json:"owner,omitempty"`
}

// APIClientRequest is the request payload used to create and update apiClients
type APIClientRequest struct {
	Metadata         *RequestMetadata
	Role             string `json:"role,omitempty"`
	DocumentName     string `json:"documentName,omitempty"`
	DocumentType     string `json:"documentType,omitempty"`
	DocumentLocation string `json:"documentLocation,omitempty"`
	DocumentGUID     string `json:"documentGuid,omitempty"`
	StreamID         string `json:"streamId,omitempty"`
	Online           bool   `json:"online,omitempty"`
}

// APIClientService is the service that communicates with the APIClients API
type APIClientService struct {
	client *Client
}

// List retrieves a list of APIClients
func (s *APIClientService) List(ctx context.Context) ([]APIClient, *http.Response, error) {
	resource := new([]APIClient)

	req, err := s.client.NewRequest(ctx, http.MethodGet, apiClientBasePath, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Get retrieves a specific apiClient indexed by it's ID
func (s *APIClientService) Get(ctx context.Context, id string) (APIClient, *http.Response, error) {
	resource := new(APIClient)

	req, err := s.client.NewRequest(ctx, http.MethodGet, apiClientBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Create will create a new apiClient.
func (s *APIClientService) Create(ctx context.Context, new APIClientRequest) (APIClient, *http.Response, error) {
	resource := APIClient{} // new(APIClient)

	req, err := s.client.NewRequest(ctx, http.MethodPost, apiClientBasePath, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// Update will update apiClient indexed by it's ID.
func (s *APIClientService) Update(ctx context.Context, id string, update APIClientRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, apiClientBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a specific apiClient indexed by it's ID
func (s *APIClientService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiClientBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateComment will create a new comment for a specific APIClient indexed by id.
func (s *APIClientService) CreateComment(ctx context.Context, id string, new CommentRequest) (Comment, *http.Response, error) {
	resource := Comment{} // new(APIClient)

	req, err := s.client.NewRequest(ctx, http.MethodPost, "comments/"+apiClientBasePath+"/"+id, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// GetComments will get comments for a specific APIClient indexed by ID.
func (s *APIClientService) GetComments(ctx context.Context, id string) ([]Comment, *http.Response, error) {
	resource := new([]Comment)

	req, err := s.client.NewRequest(ctx, http.MethodGet, "comments/"+apiClientBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

package gospeckle

import (
	"context"
	"fmt"
	"net/http"
)

const objectBasePath = "objects"

// Object is the request response when fetching objects
type Object struct {
	Metadata      Metadata
	ApplicationID string                 `json:"applicationId,omitempty"`
	Hash          string                 `json:"hash,omitempty"`
	GeometryHash  string                 `json:"geometryHash,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Type          string                 `json:"type"`
	PartOf        []string               `json:"partOf,omitempty"`
	Parent        []string               `json:"parent,omitempty"`
	Children      []string               `json:"children,omitempty"`
	Ancestors     []string               `json:"ancestors,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// ObjectRequest is the request payload used to create and update objects
type ObjectRequest struct {
	Metadata      *RequestMetadata
	ApplicationID string                 `json:"applicationId,omitempty"`
	Hash          string                 `json:"hash,omitempty"`
	GeometryHash  string                 `json:"geometryHash,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Type          string                 `json:"type"`
	PartOf        []string               `json:"partOf,omitempty"`
	Parent        []string               `json:"parent,omitempty"`
	Children      []string               `json:"children,omitempty"`
	Ancestors     []string               `json:"ancestors,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// ObjectGetBulkQuery is a map of key value pairs used to generate
// a http query of the form `?key1=value1&key2=value2` etc...
type ObjectGetBulkQuery map[string]string

// ObjectGetBulkIDs is an array of object IDs to limit a
// GetBulk query to.
type ObjectGetBulkIDs []string

// ObjectService is the service that communicates with the Objects API
type ObjectService struct {
	client *Client
}

// List retrieves a list of Objects
func (s *ObjectService) List(ctx context.Context) ([]Object, *http.Response, error) {
	resource := new([]Object)

	req, err := s.client.NewRequest(ctx, http.MethodGet, objectBasePath, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Get retrieves a specific object indexed by it's ID
func (s *ObjectService) Get(ctx context.Context, id string) (Object, *http.Response, error) {
	resource := new(Object)

	req, err := s.client.NewRequest(ctx, http.MethodGet, objectBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Create will create a new object.
func (s *ObjectService) Create(ctx context.Context, new ObjectRequest) (Object, *http.Response, error) {
	resource := Object{} // new(Object)

	req, err := s.client.NewRequest(ctx, http.MethodPost, objectBasePath, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// Update will update object indexed by it's ID.
func (s *ObjectService) Update(ctx context.Context, id string, update ObjectRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, objectBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a specific object indexed by it's ID
func (s *ObjectService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, objectBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateComment will create a new comment for a specific Object indexed by id.
func (s *ObjectService) CreateComment(ctx context.Context, id string, new CommentRequest) (Comment, *http.Response, error) {
	resource := Comment{} // new(Object)

	req, err := s.client.NewRequest(ctx, http.MethodPost, "comments/"+objectBasePath+"/"+id, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// GetComments will get comments for a specific Object indexed by ID.
func (s *ObjectService) GetComments(ctx context.Context, id string) ([]Comment, *http.Response, error) {
	resource := new([]Comment)

	req, err := s.client.NewRequest(ctx, http.MethodGet, "comments/"+objectBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// GetBulk will search for objects in a given range of IDs using a query.
func (s *ObjectService) GetBulk(ctx context.Context, idList ObjectGetBulkIDs, query ObjectGetBulkQuery) ([]Object, *http.Response, error) {
	resource := []Object{} // new(Object)
	queryPath := "?"

	for k, v := range query {
		queryPath += k + "=" + v + "&"
	}

	if string(queryPath[len(queryPath)-1]) == "&" {
		queryPath = queryPath[len(queryPath)-1:]
	}

	fmt.Println(queryPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, objectBasePath+"/"+queryPath, idList)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

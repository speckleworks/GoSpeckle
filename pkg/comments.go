package gospeckle

import (
	"context"
	"net/http"
)

const commentBasePath = "comments"

// CommentResource is the parent resource of the comment
type CommentResource struct {
	ResourceType string `json:"resourceType,omitempty"`
	ResourceID   string `json:"resourceId,omitempty"`
}

// View contains the information for a camera view (can be expanded to hold other scene settings)
type View struct {
}

// Comment is the request response when fetching comments
type Comment struct {
	Metadata       Metadata
	Resource       CommentResource `json:"resource,omitempty"`
	Flagged        bool            `json:"flagged,omitempty"`
	Closed         bool            `json:"closed,omitempty"`
	AssignedTo     []string        `json:"assignedTo,omitempty"`
	Labels         []string        `json:"labels,omitempty"`
	Text           string          `json:"text,omitempty"`
	OtherResources []string        `json:"otherResources,omitempty"`
	View           View            `json:"view,omitempty"`
	Screenshot     string          `json:"screenshot,omitempty"`
}

// CommentRequest is the request payload used to create and update comments
type CommentRequest struct {
	Metadata       *RequestMetadata
	Resource       CommentResource `json:"resource,omitempty"`
	Flagged        bool            `json:"flagged,omitempty"`
	Closed         bool            `json:"closed,omitempty"`
	AssignedTo     []string        `json:"assignedTo,omitempty"`
	Labels         []string        `json:"labels,omitempty"`
	Text           string          `json:"text,omitempty"`
	OtherResources []string        `json:"otherResources,omitempty"`
	View           View            `json:"view,omitempty"`
	Screenshot     string          `json:"screenshot,omitempty"`
}

// CommentService is the service that communicates with the Comments API
type CommentService struct {
	client *Client
}

// Get retrieves a specific comment indexed by it's ID
func (s *CommentService) Get(ctx context.Context, id string) (Comment, *http.Response, error) {
	resource := new(Comment)

	req, err := s.client.NewRequest(ctx, http.MethodGet, commentBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Update will update comment indexed by it's ID.
func (s *CommentService) Update(ctx context.Context, id string, update CommentRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, commentBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a specific comment indexed by it's ID
func (s *CommentService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, commentBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

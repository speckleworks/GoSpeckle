package gospeckle

import (
	"context"
	"net/http"
)

const projectBasePath = "projects"

// ProjectPermissions represents project specific user permission data.
type ProjectPermissions struct {
	CanRead  []string `json:"canRead"`
	CanWrite []string `json:"canWrite"`
}

// Project is the request response when fetching projects
type Project struct {
	Metadata
	Name        string             `json:"name"`
	Tags        []string           `json:"tags,omitempty"`
	Streams     []string           `json:"streams,omitempty"`
	Permissions ProjectPermissions `json:"permissions,omitempty"`
}

// ProjectStreamResponse is the response payload from adding/removing streams to
// a project.
type ProjectStreamResponse struct {
	Project Project `json:"project"`
	Stream  string  `json:"stream"`
}

// ProjectRequest is the request payload used to create and update projects
type ProjectRequest struct {
	RequestMetadata
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Streams []string `json:"streams"`
}

// ProjectService is the service that communicates with the Projects API
type ProjectService struct {
	client *Client
}

// List retrieves a list of Projects
func (s *ProjectService) List(ctx context.Context) ([]Project, *http.Response, error) {
	resource := new([]Project)

	req, err := s.client.NewRequest(ctx, http.MethodGet, projectBasePath, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Get retrieves a specific project indexed by it's ID
func (s *ProjectService) Get(ctx context.Context, id string) (Project, *http.Response, error) {
	resource := new(Project)

	req, err := s.client.NewRequest(ctx, http.MethodGet, projectBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Create will create a new project.
func (s *ProjectService) Create(ctx context.Context, new ProjectRequest) (Project, *http.Response, error) {
	resource := Project{} // new(Project)

	req, err := s.client.NewRequest(ctx, http.MethodPost, projectBasePath, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// Update will update project indexed by it's ID.
func (s *ProjectService) Update(ctx context.Context, id string, update ProjectRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, projectBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a specific project indexed by it's ID
func (s *ProjectService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, projectBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateComment will create a new comment for a specific Project indexed by id.
func (s *ProjectService) CreateComment(ctx context.Context, id string, new CommentRequest) (Comment, *http.Response, error) {
	resource := Comment{} // new(Project)

	req, err := s.client.NewRequest(ctx, http.MethodPost, "comments/"+projectBasePath+"/"+id, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// GetComments will get comments for a specific Project indexed by ID.
func (s *ProjectService) GetComments(ctx context.Context, id string) ([]Comment, *http.Response, error) {
	resource := new([]Comment)

	req, err := s.client.NewRequest(ctx, http.MethodGet, "comments/"+projectBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// AddStream adds a stream to a project by indexing them both by ID.
func (s *ProjectService) AddStream(ctx context.Context, id string, streamID string) (ProjectStreamResponse, *http.Response, error) {
	resource := new(ProjectStreamResponse)

	req, err := s.client.NewRequest(ctx, http.MethodPut, projectBasePath+"/"+id+"/addstream/"+streamID, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// RemoveStream adds a stream to a project by indexing them both by ID.
func (s *ProjectService) RemoveStream(ctx context.Context, id string, streamID string) (ProjectStreamResponse, *http.Response, error) {
	resource := new(ProjectStreamResponse)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, projectBasePath+"/"+id+"/removestream/"+streamID, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// AddUser adds a user to a project by indexing them both by ID.
func (s *ProjectService) AddUser(ctx context.Context, id string, userID string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, projectBasePath+"/"+id+"/adduser/"+userID, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpgradeUser adds a user to a project by indexing them both by ID.
func (s *ProjectService) UpgradeUser(ctx context.Context, id string, userID string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, projectBasePath+"/"+id+"/upgradeuser/"+userID, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DowngradeUser adds a user to a project by indexing them both by ID.
func (s *ProjectService) DowngradeUser(ctx context.Context, id string, userID string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, projectBasePath+"/"+id+"/downgradeuser/"+userID, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RemoveUser adds a user to a project by indexing them both by ID.
func (s *ProjectService) RemoveUser(ctx context.Context, id string, userID string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, projectBasePath+"/"+id+"/removeuser/"+userID, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

package gospeckle

import (
	"context"
	"encoding/json"
	"net/http"
)

const streamBasePath = "streams"

// StreamPermissions represents stream specific user permission data.
type StreamPermissions struct {
	CanRead  []string `json:"canRead"`
	CanWrite []string `json:"canWrite"`
}

// LayerProperties are a set of visualisation properties
// attached to each layer.
type LayerProperties struct {
	Color     interface{} `json:"color,omitemtpy"`
	Visible   bool        `json:"visible,omitempty"`
	PointSize float64     `json:"pointSize,omitempty"`
	LineWidth float64     `json:"lineWidth,omitempty"`
	Shininess float64     `json:"shininess,omitempty"`
	Smooth    bool        `json:"smooth,omitempty"`
	ShowEdges bool        `json:"showEdges,omitempty"`
	WireFrame bool        `json:"wireframe,omitempty"`
}

// Layer is an object attached to a stream that objects also belonging
// to the stream can be attributed to.
type Layer struct {
	GUID        string          `json:"guid,omitempty"`
	Name        string          `json:"name,omitempty"`
	OrderIndex  int             `json:"orderIndex,omitempty"`
	StartIndex  int             `json:"startIndex,omitempty"`
	ObjectCount int             `json:"objectCount,omitempty"`
	Topology    string          `json:"topology,omitempty"`
	Properties  LayerProperties `json:"properties,omitempty"`
}

// Stream is the request response when fetching streams
type Stream struct {
	Metadata
	StreamID      string   `json:"streamId,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	CommitMessage string   `json:"commit,omitempty"`
	Objects       []Object `json:"objects,omitempty"`
	Layers        []Layer  `json:"layers,omitempty"`
	Parent        []string `json:"parents,omitempty"`
	Children      []string `json:"children,omitempty"`
}

// StreamStreamResponse is the response payload from adding/removing streams to
// a stream.
type StreamStreamResponse struct {
	Stream Stream `json:"stream"`
}

// StreamRequest is the request payload used to create and update streams
type StreamRequest struct {
	RequestMetadata
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	CommitMessage string    `json:"commit,omitempty"`
	Objects       []*Object `json:"objects,omitempty"`
	Layers        []*Layer  `json:"layers,omitempty"`
	Parent        []string  `json:"parents,omitempty"`
	Children      []string  `json:"children,omitempty"`
}

// StreamCloneRequest is the request payload used to clone a stream
type StreamCloneRequest struct {
	Name string `json:"name,omitemtpy"`
}

// StreamCloneResponse is the response from a Stream Clone request
type StreamCloneResponse struct {
	Clone  Stream `json:"clone"`
	Parent Stream `json:"parent"`
}

// StreamDiff is a diff object returned from a diff query.
type StreamDiff struct {
	Common []interface{} `json:"common,omitempty"`
	InA    []interface{} `json:"inA,omitempty"`
	InB    []interface{} `json:"inB,omitempty"`
}

// StreamDiffResponse is the response from a Stream Diff request
type StreamDiffResponse struct {
	Objects StreamDiff `json:"objects,omitempty"`
	Layers  StreamDiff `json:"layers,omitempty"`
}

// StreamService is the service that communicates with the Streams API
type StreamService struct {
	client *Client
}

// List retrieves a list of Streams
func (s *StreamService) List(ctx context.Context) ([]Stream, *http.Response, error) {
	resource := new([]Stream)

	req, err := s.client.NewRequest(ctx, http.MethodGet, streamBasePath, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Get retrieves a specific stream indexed by it's ID
func (s *StreamService) Get(ctx context.Context, id string) (Stream, *http.Response, error) {
	resource := new(Stream)

	req, err := s.client.NewRequest(ctx, http.MethodGet, streamBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Create will create a new stream.
func (s *StreamService) Create(ctx context.Context, new StreamRequest) (Stream, *http.Response, error) {
	resource := Stream{} // new(Stream)

	req, err := s.client.NewRequest(ctx, http.MethodPost, streamBasePath, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// Update will update stream indexed by it's ID.
func (s *StreamService) Update(ctx context.Context, id string, update StreamRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, streamBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a specific stream indexed by it's ID
func (s *StreamService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, streamBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateComment will create a new comment for a specific Stream indexed by id.
func (s *StreamService) CreateComment(ctx context.Context, id string, new CommentRequest) (Comment, *http.Response, error) {
	resource := Comment{} // new(Stream)

	req, err := s.client.NewRequest(ctx, http.MethodPost, "comments/"+streamBasePath+"/"+id, new)
	if err != nil {
		return resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, &resource, false)
	if err != nil {
		return resource, nil, err
	}

	return resource, resp, nil
}

// GetComments will get comments for a specific Stream indexed by ID.
func (s *StreamService) GetComments(ctx context.Context, id string) ([]Comment, *http.Response, error) {
	resource := new([]Comment)

	req, err := s.client.NewRequest(ctx, http.MethodGet, "comments/"+streamBasePath+"/"+id, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, false)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// Clone will create a new stream with a given code name.
func (s *StreamService) Clone(ctx context.Context, streamID string, cloneName string) (StreamCloneResponse, *http.Response, error) {
	resource := new(StreamCloneResponse)
	payload := StreamCloneRequest{Name: cloneName}

	req, err := s.client.NewRequest(ctx, http.MethodPost, streamBasePath+"/"+streamID+"/clone", payload)
	if err != nil {
		return *resource, nil, err
	}

	resp, data, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return *resource, nil, err
	}

	encoder, err := json.Marshal(data)
	if err != nil {
		return *resource, resp, err
	}
	err = json.Unmarshal(encoder, &resource)
	if err != nil {
		return *resource, resp, err
	}

	return *resource, resp, nil
}

// Diff will return a diff between two streams identified by their IDs.
func (s *StreamService) Diff(ctx context.Context, streamA string, streamB string) (StreamDiffResponse, *http.Response, error) {
	resource := new(StreamDiffResponse)

	req, err := s.client.NewRequest(ctx, http.MethodGet, streamBasePath+"/"+streamA+"/diff/"+streamB, nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, data, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return *resource, nil, err
	}

	encoder, err := json.Marshal(data)
	if err != nil {
		return *resource, resp, err
	}
	err = json.Unmarshal(encoder, &resource)
	if err != nil {
		return *resource, resp, err
	}

	return *resource, resp, nil
}

// ListObjects retrieves a list of objects in the Stream.
func (s *StreamService) ListObjects(ctx context.Context, streamID string) ([]map[string]interface{}, *http.Response, error) {
	resource := new([]map[string]interface{})

	req, err := s.client.NewRequest(ctx, http.MethodGet, streamBasePath+"/"+streamID+"/objects", nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

// ListClients retrieves a list of clients suscribed to the Stream.
func (s *StreamService) ListClients(ctx context.Context, streamID string) ([]APIClient, *http.Response, error) {
	resource := new([]APIClient)

	req, err := s.client.NewRequest(ctx, http.MethodGet, streamBasePath+"/"+streamID+"/clients", nil)
	if err != nil {
		return *resource, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, resource, true)
	if err != nil {
		return *resource, nil, err
	}

	return *resource, resp, nil
}

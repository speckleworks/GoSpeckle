package gospeckle

import "time"

// Metadata is the common metadata all resources (bar accounts) possess when returned
// from the Speckle Server.
type Metadata struct {
	ID                string    `json:"_id,omitempty"`
	Private           bool      `json:"private"`
	CanRead           []string  `json:"canRead"`
	CanWrite          []string  `json:"canWrite"`
	Owner             string    `json:"owner"`
	AnonymousComments bool      `json:"anonymousComments"`
	Comments          []string  `json:"comments"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
	Version           int       `json:"__v"`
}

// RequestMetadata is the common metadata all resources (bar accounts) should posess
// when creating/updating resources to the Speckle Server.
type RequestMetadata struct {
	Private           bool     `json:"private,omitempty"`
	CanRead           []string `json:"canRead,omitempty"`
	CanWrite          []string `json:"canWrite,omitempty"`
	AnonymousComments bool     `json:"anonymousComments,omitempty"`
}

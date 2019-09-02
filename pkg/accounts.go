package gospeckle

import (
	"context"
	"net/http"
)

const accountBasePath = "accounts"

// Account is the struct that defines and account object
type Account struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Role     string `json:"role,omitempty"`
	Archived bool   `json:"archived,omitempty"`
	Email    string `json:"email"`
	Company  string `json:"company,omitempty"`
}

// AccountService is the service that communicates with the Accounts API
type AccountService struct {
	client *Client
}

// AccountSearchRequest is the request payload used to search for users/accounts
type AccountSearchRequest struct {
	SearchString string `json:"searchString"`
}

// AccountUpdateRequest is the request payload used to update account information
type AccountUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Email   string `json:"email,omitempty"`
	Company string `json:"company,omitempty"`
}

// AccountRoleUpdateRequest is the request payload used to update the role of an account
type AccountRoleUpdateRequest struct {
	Role string `json:"role"`
}

// Search is a request made to the server to search for accounts using the AccountSearchRequest
func (s *AccountService) Search(ctx context.Context, search string) ([]Account, *http.Response, error) {
	accounts := new([]Account)

	searchRequest := AccountSearchRequest{SearchString: search}
	req, err := s.client.NewRequest(ctx, http.MethodPost, accountBasePath+"/search", searchRequest)
	if err != nil {
		return *accounts, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, accounts, true)

	if err != nil {
		return *accounts, nil, err
	}

	return *accounts, resp, nil
}

// Get retrieves a specific account indexed by it's ID
func (s *AccountService) Get(ctx context.Context, id string) (Account, *http.Response, error) {
	account := new(Account)

	req, err := s.client.NewRequest(ctx, http.MethodGet, accountBasePath+"/"+id, nil)
	if err != nil {
		return *account, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, account, false)
	if err != nil {
		return *account, nil, err
	}

	return *account, resp, nil
}

// Me retrieves the current user's account as identified by using the authenticaton token
func (s *AccountService) Me(ctx context.Context) (Account, *http.Response, error) {
	account := new(Account)

	req, err := s.client.NewRequest(ctx, http.MethodGet, accountBasePath, nil)
	if err != nil {
		return *account, nil, err
	}

	resp, _, err := s.client.Do(ctx, req, account, false)
	if err != nil {
		return *account, nil, err
	}

	return *account, resp, nil
}

// Update will update the current user's account as identified by using the authentication token.
func (s *AccountService) Update(ctx context.Context, update AccountUpdateRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, accountBasePath, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetRole will update the role of an account indexed by it's ID. Note: only admins can use this command.
func (s *AccountService) SetRole(ctx context.Context, id string, update AccountRoleUpdateRequest) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPut, accountBasePath+"/"+id, update)
	if err != nil {
		return nil, err
	}

	resp, _, err := s.client.Do(ctx, req, nil, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

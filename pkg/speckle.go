package gospeckle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultHost             = "hestia.speckle.works"
	defaultHTTPScheme       = "https"
	defaultWebsocketsScheme = "wss"
	libraryVersion          = "0.1.0"
	userAgent               = "gospeckle/" + libraryVersion
	mediaType               = "application/json"
)

// ResponseData is a Speckle Server response. This wraps the standard http.Response returned from DigitalOcean.
type ResponseData struct {
	// *http.Response
	Success   bool                     `json:"success"`
	Message   string                   `json:"message"`
	Resource  map[string]interface{}   `json:"resource"`
	Resources []map[string]interface{} `json:"resources"`
	Clone     map[string]interface{}   `json:"clone"`
	Parent    map[string]interface{}   `json:"parent"`
	Objects   map[string]interface{}   `json:"objects"`
	Layers    map[string]interface{}   `json:"layers"`
}

// ErrorResponse is either a Speckle Server response resulting from an incorrect API call or an error unhandled by the server.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
	// Error message
	Message string `json:"message"`
}

// Client managed communication with a SpeckleServer.
type Client struct {
	client        *http.Client
	APIURL        *url.URL    `json:"api_url"`
	WebsocketsURL *url.URL    `json:"websockets_url"`
	Token         string      `json:"token"`
	User          *ClientUser `json:"user"`

	// Services used for communicating with the API
	Account   AccountService
	APIClient APIClientService
	Comment   CommentService
	Stream    StreamService
	Project   ProjectService
	Object    ObjectService
}

// NewClient returns a new Speckle server API client.
func NewClient(httpClient *http.Client, apiURL *url.URL, websocketsURL *url.URL, apiVersion string, authToken string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	if apiURL == nil {
		apiURL = &url.URL{
			Scheme: defaultHTTPScheme,
			Host:   defaultHost,
		}
	}

	if apiVersion == "" {
		apiVersion = "v1"
	}

	if apiURL.Path == "" {
		apiURL.Path = "api/" + apiVersion + "/"
	}

	if websocketsURL == nil {
		websocketsURL = &url.URL{
			Scheme: defaultWebsocketsScheme,
			Host:   defaultHost,
		}
	}

	c := &Client{client: httpClient, APIURL: apiURL, WebsocketsURL: websocketsURL, Token: authToken}

	c.Account = AccountService{client: c}
	c.APIClient = APIClientService{client: c}
	c.Comment = CommentService{client: c}
	c.Project = ProjectService{client: c}
	c.Stream = StreamService{client: c}
	c.Object = ObjectService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// APIURL of the Client. Relative URLs should not be preceded by a forward slash. If provided, the value pointed
// to by the body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	url := c.APIURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("Authorization", c.Token)
	return req, nil
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *ResponseData {
	response := new(ResponseData)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	// var data map[string][]map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		// TODO: Better error handling!
		fmt.Println(err)
		return response
	}
	// response.Data = data
	return response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}, isList bool) (*http.Response, *ResponseData, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	responseData := newResponse(resp)

	defer resp.Body.Close()

	err = CheckResponse(responseData, resp)
	if err != nil {
		return resp, nil, err
	}

	var encoder bytes.Buffer
	enc := json.NewEncoder(&encoder)
	dec := json.NewDecoder(&encoder)

	if isList == true {
		enc.Encode(responseData.Resources)
	} else {
		enc.Encode(responseData.Resource)
	}

	dec.Decode(&v)
	// &v = response.Data

	return resp, responseData, err
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// CheckResponse determines whether an error has occured and returns the error
// if such is the case.
func CheckResponse(respData *ResponseData, r *http.Response) error {
	if respData.Success == false {
		return &ErrorResponse{
			Message:  respData.Message,
			Response: r,
		}
	}

	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

// AuthPayload is the JSON payload sent to the Speckle Server to Authenticate an existing user.
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ClientUser is the parsed response from an login call.
type ClientUser struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Role      string    `json:"role"`
	Private   bool      `json:"private"`
	Archived  bool      `json:"archived"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	APIToken  string    `json:"apitoken"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Token     string    `json:"token"`
}

// Login is the function used to login an existing user with an email and a password.
func (c *Client) Login(ctx context.Context, email string, password string, persistentToken bool) error {
	authPayload := AuthPayload{
		Email:    email,
		Password: password,
	}
	req, err := c.NewRequest(ctx, http.MethodPost, "accounts/login", authPayload)
	if err != nil {
		return err
	}

	clientUser := new(ClientUser)

	_, _, err = c.Do(ctx, req, clientUser, false)

	c.User = clientUser

	if persistentToken {
		c.Token = clientUser.APIToken
	} else {
	c.Token = clientUser.Token
	}

	return err
}

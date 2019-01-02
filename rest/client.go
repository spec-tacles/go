package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client represents a REST client
type Client struct {
	Token       string
	HTTP        *http.Client
	Buckets     *sync.Map
	GlobalReset time.Time
	APIVersion  string
	URLHost     string
	URLScheme   string
}

// NewClient makes a new client
func NewClient(token string) *Client {
	return &Client{
		Token:       token,
		HTTP:        http.DefaultClient,
		Buckets:     &sync.Map{},
		GlobalReset: time.Time{},
		APIVersion:  "6",
		URLHost:     "discordapp.com",
		URLScheme:   "https",
	}
}

// GloballyLimited returns whether we're currently globally ratelimited
func (c *Client) GloballyLimited() bool {
	return time.Now().Before(c.GlobalReset)
}

// DoJSON is a utility method for making JSON requests to the Discord API
func (c *Client) DoJSON(method, url string, body io.Reader, respBody interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(respBody)
	if err != nil {
		return err
	}

	res.Body.Close()
	return nil
}

// Do makes a raw HTTP ratelimited request to the Discord API
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	route := MakeRoute(req.URL.Path)
	req.URL.Path = "/api/v" + c.APIVersion + req.URL.Path

	bucket, _ := c.Buckets.Load(route)
	if bucket == nil {
		bucket = NewBucket(c, route)
		c.Buckets.Store(route, bucket)
	}

	if req.URL.Host == "" {
		req.URL.Host = c.URLHost
	}

	if req.URL.Scheme == "" {
		req.URL.Scheme = c.URLScheme
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "DiscordBot (https://github.com/spec-tacles/spectacles, v1)")
	}

	if req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bot "+c.Token)
	}

	return bucket.(*Bucket).Do(req)
}

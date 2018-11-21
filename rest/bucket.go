package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Bucket represents a ratelimit bucket
type Bucket struct {
	mux *sync.Mutex

	Client *Client
	Route  string

	Remaining int64
	Reset     time.Time
	Limit     int64
}

type ratelimited struct {
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
	Global     bool   `json:"global"`
}

// NewBucket makes a new bucket
func NewBucket(client *Client, route string) *Bucket {
	return &Bucket{
		mux:       &sync.Mutex{},
		Client:    client,
		Route:     route,
		Remaining: 1,
		Reset:     time.Time{},
		Limit:     1,
	}
}

// Make a request in this bucket
func (b *Bucket) Make(req *http.Request) (res *http.Response, err error) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.Client.GloballyLimited() {
		wait := time.Until(b.Client.GlobalReset)
		<-time.After(wait)
	}

	if b.Remaining <= 0 {
		wait := time.Until(b.Reset)
		<-time.After(wait)
		b.Remaining = b.Limit
	}

	res, err = b.Client.HTTP.Do(req)
	if err != nil {
		return
	}

	limit := res.Header.Get("x-ratelimit-limit")
	if limit != "" {
		b.Limit, err = strconv.ParseInt(limit, 10, 32)
		if err != nil {
			return
		}
	}

	remaining := res.Header.Get("x-ratelimit-remaining")
	if remaining != "" {
		b.Remaining, err = strconv.ParseInt(remaining, 10, 32)
	}

	if err != nil {
		return
	}

	switch {
	case res.StatusCode == http.StatusTooManyRequests:
		err = b.handle429(res)
	case res.StatusCode >= 500 && res.StatusCode < 600:
		return b.handle500(req) // retry on 500
	default:
		err = b.handleNormal(res)
	}
	return
}

func (b *Bucket) handleNormal(res *http.Response) (err error) {
	// handle normal bucket resets
	reset := res.Header.Get("x-ratelimit-reset")
	if reset == "" {
		return nil
	}

	var resetTime int64
	resetTime, err = strconv.ParseInt(reset, 10, 64)
	if err != nil {
		return
	}

	sent, err := time.Parse(res.Header.Get("date"), time.RFC1123)
	if err != nil {
		sent = time.Now()
	}

	// calculate clock difference when determining reset time based on header timestamp
	diff := time.Now().Sub(sent)
	b.Reset = time.Unix(resetTime, 0).Add(diff)
	return nil
}

func (b *Bucket) handle429(res *http.Response) (err error) {
	var (
		body  = ratelimited{}
		bytes []byte
	)

	bytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	res.Body.Close()
	json.Unmarshal(bytes, &body)

	// handle 429 resets
	reset := time.Now().Add(time.Duration(body.RetryAfter))
	if body.Global {
		b.Client.GlobalReset = reset
	} else {
		b.Reset = reset
	}

	return
}

func (b *Bucket) handle500(req *http.Request) (*http.Response, error) {
	<-time.After(5 * time.Second)
	return b.Make(req)
}

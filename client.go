package deluge

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	url        string
	password   string
	httpClient *http.Client

	cookies      []*http.Cookie
	cookiesMutex sync.Mutex

	lastRequest      time.Time
	lastRequestMutex sync.Mutex
}

func New(url, password string) *Client {
	return &Client{
		url:      url,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

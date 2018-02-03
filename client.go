package deluge

import (
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	url        string
	httpClient *http.Client
	cookies    []*http.Cookie

	requestCount      int
	requestCountMutex sync.Mutex
}

func New(url, password string) (*Client, error) {
	c := Client{
		url: url,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	if err := c.login(password); err != nil {
		return nil, errors.Wrap(err, "failed login")
	}

	return &c, nil
}

func (c *Client) login(password string) error {
	var result bool
	err := c.rpc(Login, []string{password}, &result)
	if err != nil {
		return errors.Wrap(err, "failed rpc")
	}

	if !result {
		return errors.New("rpc returned false")
	}

	return nil
}

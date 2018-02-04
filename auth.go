package deluge

import (
	"time"

	"github.com/pkg/errors"
)

const minAuthFrequency = 10 * time.Minute

func (c *Client) loginIfExpired() error {
	c.lastRequestMutex.Lock()
	expired := time.Since(c.lastRequest) > minAuthFrequency
	c.lastRequestMutex.Unlock()

	if !expired {
		return nil
	}

	err := c.login()
	return errors.Wrap(err, "failed to login")
}

func (c *Client) login() error {
	var result bool
	err := c.rpc(Login, []string{c.password}, &result)
	if err != nil {
		return errors.Wrap(err, "failed rpc")
	}

	if !result {
		return errors.New("rpc returned false")
	}

	return nil
}

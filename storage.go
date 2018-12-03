package deluge

import "github.com/pkg/errors"

func (c *Client) MoveStorage(id, dest string) error {
	if err := c.loginIfExpired(); err != nil {
		return errors.Wrap(err, "failed login if expired")
	}

	ids := []string{id}

	err := c.rpc(MoveStorage, []interface{}{ids, dest}, nil)
	return errors.Wrap(err, "failed rpc")
}

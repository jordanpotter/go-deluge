package deluge

import (
	"path/filepath"

	"github.com/pkg/errors"
)

type Torrent struct {
	ID        string
	Name      string `json:"name"`
	Directory string `json:"save_path"`

	State    string  `json:"state"`
	Paused   bool    `json:"paused"`
	Progress float32 `json:"progress"`
	Finished bool    `json:"is_finished"`

	Size  uint64  `json:"total_size"`
	Ratio float32 `json:"ratio"`

	Tracker    string `json:"tracker"`
	Peers      int    `json:"num_peers"`
	TotalPeers int    `json:"total_peers"`
	TotalSeeds int    `json:"total_seeds"`

	Created Time `json:"time_added"`
}

func (t *Torrent) Path() string {
	return filepath.Join(t.Directory, t.Name)
}

func (c *Client) AddTorrent(url, incompletePath, completePath string) (string, error) {
	if err := c.loginIfExpired(); err != nil {
		return "", errors.Wrap(err, "failed login if expired")
	}

	options := map[string]interface{}{
		"download_location":   incompletePath,
		"move_on_completed":   true,
		"move_completed_path": completePath,
	}

	var result string
	err := c.rpc(AddTorrent, []interface{}{url, options}, &result)
	return result, errors.Wrap(err, "failed rpc")
}

func (c *Client) Torrent(id string) (*Torrent, error) {
	if err := c.loginIfExpired(); err != nil {
		return nil, errors.Wrap(err, "failed login if expired")
	}

	keys := []string{}

	var result Torrent
	err := c.rpc(GetTorrent, []interface{}{id, keys}, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed rpc")
	}

	return &result, nil
}

func (c *Client) Torrents() ([]*Torrent, error) {
	if err := c.loginIfExpired(); err != nil {
		return nil, errors.Wrap(err, "failed login if expired")
	}

	filter := map[string]string{}
	keys := []string{}

	var result map[string]*Torrent
	err := c.rpc(GetTorrents, []interface{}{filter, keys}, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed rpc")
	}

	var torrents []*Torrent
	for id, torrent := range result {
		torrent.ID = id
		torrents = append(torrents, torrent)
	}

	return torrents, nil
}

func (c *Client) RemoveTorrent(id string, removeData bool) (bool, error) {
	if err := c.loginIfExpired(); err != nil {
		return false, errors.Wrap(err, "failed login if expired")
	}

	var result bool
	err := c.rpc(RemoveTorrent, []interface{}{id, removeData}, &result)
	return result, errors.Wrap(err, "failed rpc")
}

func (c *Client) PauseTorrent(id string) error {
	if err := c.loginIfExpired(); err != nil {
		return errors.Wrap(err, "failed login if expired")
	}

	err := c.rpc(PauseTorrent, []interface{}{id}, nil)
	return errors.Wrap(err, "failed rpc")
}

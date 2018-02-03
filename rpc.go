package deluge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	Login       = "auth.login"
	AddTorrent  = "core.add_torrent_url"
	GetTorrent  = "core.get_torrent_status"
	GetTorrents = "core.get_torrents_status"
)

type requestData struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type responseData struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}

func (c *Client) rpc(method string, params interface{}, dest interface{}) error {
	reqData := requestData{
		ID:     c.requestID(),
		Method: method,
		Params: params,
	}

	reqBody, err := json.Marshal(&reqData)
	if err != nil {
		return errors.Wrap(err, "failed to serialize request data")
	}

	url := c.url + "/json"
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}

	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("received status code %d", resp.StatusCode)
	}

	c.cookies = resp.Cookies()

	respData := responseData{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return errors.Wrap(err, "failed to parse response body")
	}

	if respData.Err != nil {
		return errors.New(fmt.Sprint(respData.Err))
	}

	if respData.ID != reqData.ID {
		return errors.New("mismatched request/response ids")
	}

	if err = json.Unmarshal(respData.Result, dest); err != nil {
		return errors.Wrap(err, "failed to parse result")
	}

	return nil
}

func (c *Client) requestID() int {
	c.requestCountMutex.Lock()
	defer c.requestCountMutex.Unlock()

	c.requestCount++
	return c.requestCount
}

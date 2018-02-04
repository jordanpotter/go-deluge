package deluge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	Login       = "auth.login"
	AddTorrent  = "core.add_torrent_url"
	GetTorrent  = "core.get_torrent_status"
	GetTorrents = "core.get_torrents_status"
)

type requestData struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type responseData struct {
	ID     string          `json:"id"`
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

	c.cookiesMutex.Lock()
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	c.cookiesMutex.Unlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("received status code %d", resp.StatusCode)
	}

	c.cookiesMutex.Lock()
	c.cookies = resp.Cookies()
	c.cookiesMutex.Unlock()

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

	c.lastRequestMutex.Lock()
	c.lastRequest = time.Now()
	c.lastRequestMutex.Unlock()

	return nil
}

func (c *Client) requestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

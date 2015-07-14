package slackhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Message to send to Slack's Incoming WebHook API.
//
// See https://api.slack.com/incoming-webhooks
type Message struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// Poster interface is the methods of http.Client required by Client to ease
// testing.
type Poster interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// Client for Slack's Incoming WebHook API.
type Client struct {
	url        string
	HTTPClient Poster
}

// New Slack Incoming WebHook Client using http.DefaultClient for its Poster.
func New(url string) *Client {
	return &Client{url: url, HTTPClient: http.DefaultClient}
}

// Simple text message.
func (c *Client) Simple(msg string) error {
	return c.Send(&Message{Text: msg})
}

// Send a Message.
func (c *Client) Send(msg *Message) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Post(c.url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard response body to reuse connection
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

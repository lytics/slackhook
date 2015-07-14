package slackhook_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/lytics/slackhook"
)

type fakeposter struct {
	body []byte
}

func (p *fakeposter) Post(_, _ string, body io.Reader) (*http.Response, error) {
	p.body, _ = ioutil.ReadAll(body)
	return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(nil)), StatusCode: 200}, nil
}

func TestSimple(t *testing.T) {
	txt := "test"
	s := New("")
	fake := &fakeposter{}
	s.HTTPClient = fake
	if err := s.Simple(txt); err != nil {
		t.Fatal(err)
	}
	msg := Message{}
	if err := json.Unmarshal(fake.body, &msg); err != nil {
		t.Fatal(err)
	}
	if msg.Text != txt {
		t.Errorf("Expected text=%q but found %q", txt, msg.Text)
	}
	if msg.Channel != "" {
		t.Errorf("Expected channel to be empty but found %q", msg.Channel)
	}
	if msg.IconURL != "" {
		t.Errorf("Expected channel to be empty but found %q", msg.IconURL)
	}
	if msg.IconEmoji != "" {
		t.Errorf("Expected channel to be empty but found %q", msg.IconEmoji)
	}
}

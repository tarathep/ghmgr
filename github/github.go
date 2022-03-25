package github

import (
	"io"
	"net/http"

	"github.com/tarathep/ghmgr/login"
)

type GitHub struct {
	Auth login.Auth
}

func (github GitHub) Request(method, url string, body io.Reader) (error, int, []byte) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err, 0, nil
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+github.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, 0, nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, 0, nil
	}

	return nil, resp.StatusCode, bodyBytes
}

func (GitHub) GetMessage(bodyBytes []byte) string {
	return string(bodyBytes)
}

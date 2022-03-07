package github

import (
	"io"
	"log"
	"net/http"

	"github.com/tarathep/githuby/login"
)

type GitHub struct {
	Auth login.Auth
}

func (github GitHub) Request(method, url string, body io.Reader) (int, []byte) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+github.Auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode, bodyBytes
}

func (GitHub) GetMessage(bodyBytes []byte) string {
	return string(bodyBytes)
}

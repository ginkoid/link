package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type redir struct {
	To string `json:"to"`
}

type redirList struct {
	mux    sync.RWMutex
	redirs map[string]redir
}

type server struct {
	redirList redirList
}

func (s *server) fetchRedirs() error {
	res, getErr := http.Get(fmt.Sprintf(
		"https://api.github.com/gists/%s?client_id=%s&client_secret=%s",
		os.Getenv("APP_GITHUB_GIST_ID"),
		os.Getenv("APP_GITHUB_CLIENT_ID"),
		os.Getenv("APP_GITHUB_CLIENT_SECRET"),
	))
	if getErr != nil {
		return getErr
	}
	defer res.Body.Close()
	content, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	var contentJSON map[string]interface{}
	contentJSONErr := json.Unmarshal(content, &contentJSON)
	if contentJSONErr != nil {
		return contentJSONErr
	}
	fileContent := contentJSON["files"].(map[string]interface{})["link"].(map[string]interface{})["content"].(string)
	s.redirList.mux.Lock()
	fileJSONErr := json.Unmarshal([]byte(fileContent), &s.redirList.redirs)
	if fileJSONErr != nil {
		return fileJSONErr
	}
	s.redirList.mux.Unlock()
	return nil
}

func (s *server) startRedirFetch() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.fetchRedirs()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("strict-transport-security", "max-age=31536000; includeSubDomains; preload")
	headers.Set("content-security-policy", "default-src 'none'; sandbox")
	headers.Set("referrer-policy", "no-referrer")
	headers.Set("x-content-type-options", "nosniff")
	headers.Set("x-frame-options", "SAMEORIGIN")
	headers.Set("x-xss-protection", "1; mode=block")
	requestPath := path.Clean(r.URL.EscapedPath())
	s.redirList.mux.RLock()
	selectedRedir, ok := s.redirList.redirs[requestPath]
	s.redirList.mux.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404"))
		return
	}
	http.Redirect(w, r, selectedRedir.To, http.StatusFound)
}

func main() {
	s := server{}
	fetchErr := s.fetchRedirs()
	if fetchErr != nil {
		panic(fetchErr)
	}
	s.startRedirFetch()
	panic(http.ListenAndServe(":80", &s))
}

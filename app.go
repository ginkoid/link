package main

import (
  "os"
  "sync"
  "time"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type redir struct {
  To string `json:"to"`
}

type redirList struct {
  Mux    sync.RWMutex
  Redirs map[string]redir
}

var appRedirList redirList

func fetchRedirs() error {
  res, getErr := http.Get("https://api.github.com/gists/" + os.Getenv("APP_GIST_ID"))
  if getErr != nil {
    return getErr
  }
  defer res.Body.Close()
  content, readErr := ioutil.ReadAll(res.Body)
  if readErr != nil {
    return readErr
  }
  var contentJSON interface{}
  contentJSONErr := json.Unmarshal(content, &contentJSON)
  if contentJSONErr != nil {
    return contentJSONErr
  }
  fileContent := contentJSON.(map[string]interface{})["files"].(map[string]interface{})["link"].(map[string]interface{})["content"].(string)
  appRedirList.Mux.Lock()
  fileJSONErr := json.Unmarshal([]byte(fileContent), &appRedirList.Redirs)
  if fileJSONErr != nil {
    return fileJSONErr
  }
  appRedirList.Mux.Unlock()
  return nil
}

func startRedirFetch() {
  ticker := time.NewTicker(10 * time.Second)
  quit := make(chan struct{})
  go func() {
    for {
      select {
      case <-ticker.C:
        fetchRedirs()
      case <-quit:
        ticker.Stop()
        return
      }
    }
  }()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
  appRedirList.Mux.RLock()
  selectedRedir, ok := appRedirList.Redirs[r.URL.Path]
  appRedirList.Mux.RUnlock()
  if !ok {
    w.Write([]byte("404"))
    return
  }
  http.Redirect(w, r, selectedRedir.To, 302)
}

func main() {
  fetchErr := fetchRedirs()
  if fetchErr != nil {
    panic(fetchErr)
  }
  startRedirFetch()
  http.HandleFunc("/", handleRequest)
  panic(http.ListenAndServe(":80", nil))
}

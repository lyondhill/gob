package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path)
    hah, _ := ioutil.ReadAll(r.Body);
    fmt.Print(string(hah))
}

func main() {
  http.HandleFunc("/", handler)  
  http.ListenAndServe(":8080", nil)
}
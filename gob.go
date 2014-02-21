package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "regexp"
  "io"
  "os"

// memory checking
  "runtime"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
        b, err := json.Marshal(r)
        if err != nil {
                s = ""
                return
        }
        s = string(b)
        return
}

func fileSend(w http.ResponseWriter, r *http.Request) {
  if len(r.Header["Token"]) > 0 {
    re := regexp.MustCompile("(.*)/(file|data)")

    token := r.Header["Token"][0]
    
    key, err := ioutil.ReadFile(fmt.Sprintf("%s%s/key", config.Root, re.FindStringSubmatch(r.URL.Path)[1]))
    if err != nil { http.Error(w, err.Error(), 500) }

    if token == string(key) {

      http.ServeFile(w, r, fmt.Sprintf("%s%s/data", config.Root, re.FindStringSubmatch(r.URL.Path)[1]))
    
    } else {
      http.Error(w, Response{"success": false, "error": "I didnt get a Token"}.String(), 401)
    }
  } else {
    http.Error(w, Response{"success": false, "error": "I didnt get a Token"}.String(), 401)
  }
}

func fileInfo(w http.ResponseWriter, r *http.Request) {
  if len(r.Header["Token"]) > 0 {
    token := r.Header["Token"][0]
    key, err := ioutil.ReadFile(fmt.Sprintf("%s%s/key", config.Root, r.URL.Path))
    if err != nil { http.Error(w, err.Error(), 500) }


    if token == string(key) {
      f, err := os.Stat(fmt.Sprintf("%s%s/data", config.Root, r.URL.Path))
      if err != nil {
        http.Error(w, err.Error(), 500)
      }
      fmt.Fprint(w, Response{"size:": f.Size()})
    } else {
      http.Error(w, Response{"success": false, "error": "I didnt get a Token"}.String(), 401)
    }
    
  } else {
    http.Error(w, Response{"success": false, "error": "I didnt get a Token"}.String(), 401)
  }

  // fmt.Println("name:", f.Name())
  // fmt.Println("size:", f.Size(), "bytes")
  // fmt.Println("mode:", f.Mode())
  // fmt.Println("time:", f.ModTime())
}

func fileRecieve(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  if len(r.Header["Token"]) > 0 {

    err := os.MkdirAll(fmt.Sprintf("%s%s/", config.Root, r.URL.Path), 0777)
    if err != nil {
      http.Error(w, err.Error(), 500)
    }

    token := r.Header["Token"][0]
    key, err := os.Create(fmt.Sprintf("%s%s/key", config.Root, r.URL.Path))
    if err != nil { panic(err) }
    key.Write([]byte(token))
    if err := key.Close(); err != nil {
      panic(err)
    }

    
    fileOut, err := os.Create(fmt.Sprintf("%s%s/data", config.Root, r.URL.Path))
    if err != nil { panic(err) }
    defer func() {
      if err := fileOut.Close(); err != nil {
        panic(err)
      }
    }()

    buf := make([]byte, 1024)
    total := 0
    for {
      // read a chunk
      n, err := r.Body.Read(buf)
      total += n
      if err != nil && err != io.EOF { panic(err) }
      if n == 0 { break }

      // write a chunk
     if _, err := fileOut.Write(buf[:n]); err != nil {
        panic(err)
      }
    }

    fmt.Fprint(w, Response{"success": true, "length": total})
    // fmt.Fprint(w, "thanks for all the bytes (%d)!\n", total)
    
  } else {
    http.Error(w, Response{"success": false, "error": "I didnt get a Token"}.String(), 401)
  }


}

func handler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case "GET":
    file, _ := regexp.MatchString("/file$", r.URL.Path)
    data, _ := regexp.MatchString("/data$", r.URL.Path)
    if file || data {
      fileSend(w, r)
    } else {
      fileInfo(w, r)
    }
  case "POST":
    fileRecieve(w, r)
  }

// MEMORY STUFF
  memstats := new(runtime.MemStats)
  runtime.ReadMemStats(memstats)
  runtime.GC()
  fmt.Println("memstats before GC: bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)
  fmt.Println("**********")
  // memstats = new(runtime.MemStats)
  runtime.ReadMemStats(memstats)
  fmt.Println("memstats after GC:  bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)
}


func main() {
  http.HandleFunc("/", handler)
  fmt.Printf(":%d\n", config.Port)
  http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

// CONFIG STUFF
type Config struct {
  Port int
  Root string
}

var config Config

func init() {
  if len(os.Args) > 1 {
    println(os.Args[1])
    file, e := ioutil.ReadFile(os.Args[1])
    if e != nil {
        println("File error: %v\n", e)
        os.Exit(1)
    }
    println(string(file))
    if e := json.Unmarshal(file, &config); e != nil {
      fmt.Printf("Json error: %v\n", e)
    }
    println(config.Port)
  } else {
    config.Port = 8080
    config.Root = "./storage"
  }
}
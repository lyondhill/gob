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
  // "runtime"
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


// fi, err := os.Open("input.txt")
//     if err != nil { panic(err) }
//     // close fi on exit and check for its returned error
//     defer func() {
//         if err := fi.Close(); err != nil {
//             panic(err)
//         }
//     }()

//     // open output file

//     // make a buffer to keep chunks that are read
//     buf := make([]byte, 1024)
//     for {
//         // read a chunk
//         n, err := fi.Read(buf)
//         if err != nil && err != io.EOF { panic(err) }
//         if n == 0 { break }

//         // write a chunk
//         if _, err := fo.Write(buf[:n]); err != nil {
//             panic(err)
//         }
//     }

func fileSend(w http.ResponseWriter, path string) {
  w.Header().Set("Content-Type", "application/octet-stream")
 
  buf := make([]byte, 1024)
  re := regexp.MustCompile("(.*)/file")
  fmt.Printf("%s\n", re.FindStringSubmatch(path)[1])

  fileIn, err := os.Open(fmt.Sprintf("%s%s", config.Root, path))
  if err != nil { panic(err) }
  defer func() {
    if err := fileIn.Close(); err != nil {
      panic(err)
    }
  }()

  for {
    // read a chunk
    n, err := fileIn.Read(buf)
    if err != nil && err != io.EOF { panic(err) }
    if n == 0 { break }

    // write a chunk
    if _, err := w.Write(buf[:n]); err != nil {
      panic(err)
    }
  }

  // fmt.Fprint(w, "thanks for all the bytes (%d)!\n", total)
}

func fileInfo(w http.ResponseWriter, path string) {
  f, err := os.Stat(fmt.Sprintf("%s%s/data", config.Root, path))
  if err != nil {
    http.Error(w, err.Error(), 500)
  }
  fmt.Fprint(w, Response{"size:": f.Size()})
  // fmt.Println("name:", f.Name())
  // fmt.Println("size:", f.Size(), "bytes")
  // fmt.Println("mode:", f.Mode())
  // fmt.Println("time:", f.ModTime())
}

func fileRecieve(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  buf := make([]byte, 1024)
  total := 0

  err := os.MkdirAll(fmt.Sprintf("%s%s/", config.Root, r.URL.Path), 0777)
  if err != nil {
    http.Error(w, err.Error(), 500)
  }
  fileOut, err := os.Create(fmt.Sprintf("%s%s/data", config.Root, r.URL.Path))
  if err != nil { panic(err) }
  defer func() {
    if err := fileOut.Close(); err != nil {
      panic(err)
    }
  }()

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
}

func handler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case "GET":
    file, _ := regexp.MatchString("/file$", r.URL.Path)
    data, _ := regexp.MatchString("/data$", r.URL.Path)
    if file || data {
      fileSend(w, r.URL.Path)
    } else {
      fileInfo(w, r.URL.Path)
    }
  case "POST":
    fileRecieve(w, r)
  }

// MEMORY STUFF
  // memstats := new(runtime.MemStats)
  // runtime.ReadMemStats(memstats)
  // runtime.GC()
  // fmt.Println("memstats before GC: bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)
  // fmt.Println("**********")
  // // memstats = new(runtime.MemStats)
  // runtime.ReadMemStats(memstats)
  // fmt.Println("memstats after GC:  bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)
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
package main

import (
  "fmt"
  "net/http"
  // "io/ioutil"
  "io"
  "runtime"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // hah, _ := ioutil.ReadAll(r.Body);
    buf := make([]byte, 1024)
    total := 0
    for {
        // read a chunk
        n, err := r.Body.Read(buf)
        total += n
        if err != nil && err != io.EOF { panic(err) }
        if n == 0 { break }

        // write a chunk
        // if _, err := w.Write(buf[:n]); err != nil {
        //     panic(err)
        // }
    }
    fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path)
    fmt.Fprintf(w, "thanks for all the bytes (%d)!\n", total)


    memstats := new(runtime.MemStats)
    runtime.ReadMemStats(memstats)
    runtime.GC()
    fmt.Println("memstats before GC: bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)
    fmt.Println("**********")
    // memstats = new(runtime.MemStats)
    runtime.ReadMemStats(memstats)
    fmt.Println("memstats after GC:  bytes = ", memstats.HeapAlloc, " footprint = ", memstats.Sys)

    fmt.Println("len(hah)")
}

func main() {
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
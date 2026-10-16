package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync/atomic"
)

var total atomic.Uint64
var successful atomic.Uint64
var failed atomic.Uint64

func main() {
    version := os.Getenv("APP_VERSION")
    if version == "" { version = "local" }
    hash := os.Getenv("POD_HASH")
    if hash == "" { hash = "local" }
    failurePercent, err := strconv.Atoi(os.Getenv("FAILURE_PERCENT"))
    if err != nil || failurePercent < 0 || failurePercent > 100 { failurePercent = 0 }

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")
    })
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
        fmt.Fprintln(w, "# HELP lab_requests_total HTTP requests served by the sample app")
        fmt.Fprintln(w, "# TYPE lab_requests_total counter")
        fmt.Fprintf(w, "lab_requests_total{version=%q,pod_hash=%q,code=\"200\"} %d\n", version, hash, successful.Load())
        fmt.Fprintf(w, "lab_requests_total{version=%q,pod_hash=%q,code=\"500\"} %d\n", version, hash, failed.Load())
    })
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        n := total.Add(1)
        if int(n%100) < failurePercent {
            failed.Add(1)
            http.Error(w, "simulated v2 failure", http.StatusInternalServerError)
            return
        }
        successful.Add(1)
        fmt.Fprintf(w, "orders-api version=%s pod_hash=%s request=%d\n", version, hash, n)
    })
    log.Printf("listening on :8080; version=%s failurePercent=%d", version, failurePercent)
    log.Fatal(http.ListenAndServe(":8080", mux))
}

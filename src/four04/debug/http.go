package debug

import (
  "github.com/kellegous/pork"
  "net/http"
  "runtime/pprof"
)

func Setup(r pork.Router) {
  r.RespondWithFunc("/debug/goroutine", func(w pork.ResponseWriter, r *http.Request) {
    p := pprof.Lookup("goroutine")
    w.Header().Set("Content-Type", "text/plain;charset=utf-8")
    p.WriteTo(w, 2)
  })
}

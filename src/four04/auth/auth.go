package auth

import (
  "code.google.com/p/goauth2/oauth"
  "fmt"
  "github.com/kellegous/pork"
  "net/http"
)

type Context interface {
  ClientId() string
  ClientSecret() string
}

func urlFor(r *http.Request) string {
  sc := "http"
  if r.TLS != nil {
    sc = "https"
  }
  return fmt.Sprintf("%s://%s/auth/z", sc, r.Host)
}

func configFromRequest(ctx Context, r *http.Request) *oauth.Config {
  return &oauth.Config{
    ClientId:     ctx.ClientId(),
    ClientSecret: ctx.ClientSecret(),
    Scope:        "user:email,gist",
    AuthURL:      "https://github.com/login/oauth/authorize",
    TokenURL:     "https://github.com/login/oauth/access_token",
    RedirectURL:  urlFor(r),
  }
}

func Setup(r pork.Router, ctx Context) {
  r.RespondWithFunc("/auth/a", func(w pork.ResponseWriter, r *http.Request) {
    http.Redirect(w, r,
      configFromRequest(ctx, r).AuthCodeURL(""),
      http.StatusTemporaryRedirect)
  })

  r.RespondWithFunc("/auth/z", func(w pork.ResponseWriter, r *http.Request) {
  })
}

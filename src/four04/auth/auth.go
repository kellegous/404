package auth

import (
  "code.google.com/p/goauth2/oauth"
  "encoding/json"
  "fmt"
  "github.com/kellegous/pork"
  "net/http"
  "time"
)

type Context interface {
  ClientId() string
  ClientSecret() string
}

type ghUser struct {
  Email     string    `json:"email"`
  Name      string    `json:"name"`
  Company   string    `json:"company"`
  Location  string    `json:"location"`
  Blog      string    `json:"blog"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
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

func fetchUser(tx *oauth.Transport, user *ghUser) error {
  c := http.Client{
    Transport: tx,
  }

  res, err := c.Get("https://api.github.com/user")
  if err != nil {
    return err
  }
  defer res.Body.Close()

  return json.NewDecoder(res.Body).Decode(user)
}

func Setup(r pork.Router, ctx Context) {
  r.RespondWithFunc("/auth/a", func(w pork.ResponseWriter, r *http.Request) {
    http.Redirect(w, r,
      configFromRequest(ctx, r).AuthCodeURL(""),
      http.StatusTemporaryRedirect)
  })

  r.RespondWithFunc("/auth/z", func(w pork.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    if code == "" {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    tx := oauth.Transport{
      Config: configFromRequest(ctx, r),
    }

    _, err := tx.Exchange(code)
    if err != nil {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    var user ghUser
    if err := fetchUser(&tx, &user); err != nil {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    fmt.Fprintf(w, "%v\n", user)
  })
}

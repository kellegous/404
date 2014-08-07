package auth

import (
  "code.google.com/p/goauth2/oauth"
  "encoding/json"
  "fmt"
  "four04/config"
  "four04/store"
  "github.com/kellegous/pork"
  "net/http"
  "time"
)

type ghUser struct {
  Id        int       `json:"id"`
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

func configFromRequest(cfg *config.Config, r *http.Request) *oauth.Config {
  return &oauth.Config{
    ClientId:     cfg.OAuth.ClientId,
    ClientSecret: cfg.OAuth.ClientSecret,
    Scope:        "user:email,gist",
    AuthURL:      "https://github.com/login/oauth/authorize",
    TokenURL:     "https://github.com/login/oauth/access_token",
    RedirectURL:  urlFor(r),
  }
}

func fetchGhUser(tx *oauth.Transport, user *ghUser) error {
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

func createSessionFrom(gh *ghUser, t *oauth.Token) (*store.Session, error) {
  // user := &store.User{
  //   Id:        gh.Id,
  //   Name:      gh.Name,
  //   Email:     gh.Email,
  //   Company:   gh.Company,
  //   Location:  gh.Location,
  //   Blog:      gh.Blog,
  //   CreatedAt: gh.CreatedAt,
  //   UpdatedAt: gh.UpdatedAt,
  //   Token:     t,
  // }

  return nil, nil
}

func Setup(r pork.Router, cfg *config.Config) {
  r.RespondWithFunc("/auth/a", func(w pork.ResponseWriter, r *http.Request) {
    http.Redirect(w, r,
      configFromRequest(cfg, r).AuthCodeURL(""),
      http.StatusTemporaryRedirect)
  })

  r.RespondWithFunc("/auth/z", func(w pork.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    if code == "" {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    tx := oauth.Transport{
      Config: configFromRequest(cfg, r),
    }

    _, err := tx.Exchange(code)
    if err != nil {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    var user ghUser
    if err := fetchGhUser(&tx, &user); err != nil {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    // TODO(knoton):
    // 1 - Validate user
    // 2 - Create or update the user
    // 3 - Create a session for the user
    // 4 - Add the session key as a cookie value
    fmt.Fprintf(w, "%v\n", user)
  })
}

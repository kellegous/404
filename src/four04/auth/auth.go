package auth

import (
  "bytes"
  "code.google.com/p/goauth2/oauth"
  "encoding/json"
  "fmt"
  "four04/config"
  "four04/context"
  "four04/secure"
  "four04/store"
  "github.com/kellegous/base62"
  "github.com/kellegous/pork"
  "io"
  "net/http"
  "time"
)

const (
  AuthCookieName = "s"
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

func createSessionFrom(ctx *context.Context, gh *ghUser, t *oauth.Token) (*store.Session, error) {
  user := &store.User{
    Id:        uint64(gh.Id),
    Name:      gh.Name,
    Email:     gh.Email,
    Company:   gh.Company,
    Location:  gh.Location,
    Blog:      gh.Blog,
    CreatedAt: gh.CreatedAt,
    UpdatedAt: gh.UpdatedAt,
    Token:     t,
  }

  if err := user.Save(ctx); err != nil {
    return nil, err
  }

  sess, err := store.NewSession(user.Id)
  if err != nil {
    return nil, err
  }

  if err := sess.Save(ctx); err != nil {
    return nil, err
  }

  return sess, nil
}

func setAuthCookie(w http.ResponseWriter, cfg *config.Config, sess *store.Session) error {
  env, err := secure.Sign(sess.Key, cfg.HmacKey)
  if err != nil {
    return err
  }

  var buf bytes.Buffer
  e := base62.NewEncoder(&buf)
  if _, err := e.Write(env); err != nil {
    return err
  }
  if err := e.Close(); err != nil {
    return err
  }

  // TODO(knorton): Should also be a secure cookie.
  http.SetCookie(w, &http.Cookie{
    Name:     AuthCookieName,
    Value:    buf.String(),
    Path:     "/",
    MaxAge:   24 * 60 * 60,
    HttpOnly: true,
  })

  return nil
}

func SessionIdFromRequest(ctx *context.Context, r *http.Request) ([]byte, error) {
  c, err := r.Cookie(AuthCookieName)
  if err != nil || c.Value == "" {
    return nil, store.ErrNotFound
  }

  var buf bytes.Buffer
  if _, err := io.Copy(&buf, base62.NewDecoder(bytes.NewBufferString(c.Value))); err != nil {
    return nil, err
  }

  sid, _, err := secure.Verify(buf.Bytes(), ctx.Cfg.HmacKey)
  if err != nil {
    return nil, err
  }

  return sid, nil
}

func SessionFromRequest(ctx *context.Context, r *http.Request) (*store.Session, error) {
  sid, err := SessionIdFromRequest(ctx, r)
  if err != nil {
    return nil, err
  }

  return store.FindSession(ctx, sid)
}

func UserFromRequest(ctx *context.Context, r *http.Request) (*store.Session, *store.User, error) {
  sess, err := SessionFromRequest(ctx, r)
  if err != nil {
    return nil, nil, err
  }

  user, err := sess.User(ctx)
  if err != nil {
    return nil, nil, err
  }

  return sess, user, nil
}

func Setup(r pork.Router, ctx *context.Context) {
  r.RespondWithFunc("/auth/a", func(w pork.ResponseWriter, r *http.Request) {
    http.Redirect(w, r,
      configFromRequest(ctx.Cfg, r).AuthCodeURL(""),
      http.StatusTemporaryRedirect)
  })

  r.RespondWithFunc("/auth/z", func(w pork.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    if code == "" {
      http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
      return
    }

    tx := oauth.Transport{
      Config: configFromRequest(ctx.Cfg, r),
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

    sess, err := createSessionFrom(ctx, &user, tx.Token)
    if err != nil {
      panic(err)
    }

    if err := setAuthCookie(w, ctx.Cfg, sess); err != nil {
      panic(err)
    }
  })

  r.RespondWithFunc("/auth/exit", func(w pork.ResponseWriter, r *http.Request) {
    sid, err := SessionIdFromRequest(ctx, r)
    if err != nil {
      panic(err)
    }

    if sid == nil {
      return
    }

    if err := store.DeleteSession(ctx, sid); err != nil {
      panic(err)
    }

    http.SetCookie(w, &http.Cookie{
      Name:     AuthCookieName,
      Value:    "",
      Path:     "/",
      MaxAge:   0,
      HttpOnly: true,
    })
  })
}

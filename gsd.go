package gorilla_session_django

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/yylt/gorilla-session-django/internal/github.com/nlpodyssey/gopickle/pickle"
)

var _ sessions.Store = &Store{}

var (
	bufpool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	defaultSessionCfg = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   43200,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	}
)

type Store struct {
	cfg            *Gsd_config
	memcli         Memcacher
	sessionOptions *sessions.Options
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		newkey  string
		err     error
		session = sessions.NewSession(s, name)
	)

	opts := *s.sessionOptions
	session.Options = &opts
	session.IsNew = true

	cookie, err := r.Cookie(s.cfg.CookieKey)
	if err != nil {
		return session, err
	}
	if s.cfg.PrefixMemcache != "" {
		newkey = fmt.Sprint("%s%s", s.cfg.PrefixMemcache, cookie.Value)
	} else {
		newkey = cookie.Value
	}
	data, err := s.memcli.Get(newkey)
	if err != nil {
		return session, err
	}
	buf := bufpool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufpool.Put(buf)
	buf.WriteString(data)
	err = s.valueCheck(buf)
	return session, err
}

func (s *Store) valueCheck(reader io.Reader) error {
	var (
		value interface{}
		err   error
	)
	if !s.cfg.JsonSerializer {
		value, err = pickle.NewUnpickler(reader).Load()
	} else {
		err = json.NewDecoder(reader).Decode(&value)
	}
	if err != nil {
		return err
	}
	if s.cfg.Auth != nil {
		return s.cfg.Auth(value)
	}
	return nil
}

// Save not implementation now
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return nil
}

func NewSessionStore(memcacher Memcacher, cfg *Gsd_config, sessioncfg *sessions.Options) sessions.Store {
	if sessioncfg != nil {
		defaultSessionCfg = sessioncfg
	}

	store := &Store{
		memcli:         memcacher,
		cfg:            cfg,
		sessionOptions: defaultSessionCfg,
	}
	return store
}

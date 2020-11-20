package gorilla_session_django

import (
	"encoding/json"
	"fmt"
	"github.com/yylt/gorilla-session-django/internal/github.com/nlpodyssey/gopickle/types"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yylt/gorilla-session-django/internal/github.com/nlpodyssey/gopickle/pickle"
)

var _ sessions.Store = &Store{}

var (
	defaultSessionCfg = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   43200,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	}
)

//(TODO) the session values support more type
type Store struct {
	cfg            *Gsd_config
	memcli         Memcacher
	sessionOptions *sessions.Options
	value          map[string]string
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
		newkey = fmt.Sprintf("%s%s", s.cfg.PrefixMemcache, cookie.Value)
	} else {
		newkey = cookie.Value
	}
	data, err := s.memcli.Get(newkey)
	if err != nil {
		return session, err
	}
	return session, s.valueCheck([]byte(data))
}

func (s *Store) Values() map[string]string {
	return s.value
}

func (s *Store) valueCheck(data []byte) error {
	var (
		value interface{}
		err   error
	)
	if !s.cfg.JsonSerializer {
		value, err = pickle.Loads(string(data))
	} else {
		err = json.Unmarshal(data, &value)
	}
	if err != nil {
		return err
	}
	dict, ok := value.(*types.Dict)
	if !ok {
		return fmt.Errorf("value is not dict")
	}
	s.value = make(map[string]string)
	dict.Item(func(k, v interface{}) {
		ks, ok := k.(string)
		if !ok {
			return
		}
		vs, ok := v.(string)
		if !ok {
			return
		}
		s.value[ks] = vs
	})

	if s.cfg.Auth != nil {
		return s.cfg.Auth(s.value)
	}
	return nil
}

// Save not implementation now
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return nil
}

func NewSessionStore(memcacher Memcacher, cfg *Gsd_config, sessioncfg *sessions.Options) *Store {
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

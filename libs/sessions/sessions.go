package sessions

import (
	"log"
	"net/http"

	"gossh/libs/sessions/sessions"
	"gossh/libs/sessions/sessions/context"

	"github.com/gin-gonic/gin"
)

const (
	DefaultKey  = "github.com/gin-contrib/sessions"
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
	Options(Options)
}

// Wraps thinly gorilla-sessions methods.
// Session stores the values and optional configuration for a sessions.
type Session interface {
	// Get returns the sessions value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the sessions value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the sessions value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the sessions.
	Clear()
	// AddFlash adds a flash message to the sessions.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the sessions.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a sessions.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error
}

func Sessions(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

func SessionsMany(names []string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionsData := make(map[string]Session, len(names))
		for _, name := range names {
			sessionsData[name] = &session{name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(DefaultKey, sessionsData)
		defer context.Clear(c.Request)
		c.Next()
	}
}

type session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options Options) {
	s.Session().Options = options.ToGorillaOptions()
}

func (s *session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *session) Written() bool {
	return s.written
}

// shortcut to get sessions
func Default(c *gin.Context) Session {
	return c.MustGet(DefaultKey).(Session)
}

// shortcut to get sessions with given name
func DefaultMany(c *gin.Context, name string) Session {
	return c.MustGet(DefaultKey).(map[string]Session)[name]
}

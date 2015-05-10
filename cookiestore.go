package cookiestore

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

type CookieStore struct {
	Sessions   map[string]*Session
	CookieName string
}

type Session struct {
	c            *CookieStore
	ID           string
	Values       map[string]interface{}
	LastAccessed time.Time
}

func New(cookiename string) *CookieStore {
	return &CookieStore{Sessions: make(map[string]*Session), CookieName: cookiename}
}

func (c *CookieStore) NewSession() *Session {
	bID := make([]byte, 32)
	rand.Read(bID)
	ID := base64.StdEncoding.EncodeToString(bID)
	session := &Session{c: c, ID: ID, Values: make(map[string]interface{}), LastAccessed: time.Now()}
	c.Sessions[ID] = session
	return session
}

func (c *CookieStore) GetSession(req *http.Request) *Session {
	cookie, err := req.Cookie(c.CookieName)
	if err == http.ErrNoCookie {
		return c.NewSession()
	}
	if s, ok := c.Sessions[cookie.Value]; ok {
		return s
	}
	return c.NewSession()
}

func (s *Session) Save(rw http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  s.c.CookieName,
		Value: s.ID,
	}
	rw.Header().Add("Set-Cookie", cookie.String())
}

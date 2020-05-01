package types

import (
	"fmt"
	"github.com/cornelk/hashmap"
	"regexp"
	mock "github.com/jordwest/mock-conn"
)

var (
	// Shortcut Registry for build-in shortcut connections.
	Shortcut    *shortcut

	// Concurency Estimate concurrency required.
	Concurrency = 100

	shortcutAddress = "shortcut:%d"
	shortcutRecognizer = regexp.MustCompile(`^shortcut:`)
)

type shortcut struct {
	ports *hashmap.HashMap
}

func InitShortcut() *shortcut {
	if Shortcut == nil {
		Shortcut = &shortcut{
			ports: hashmap.New(200), // Concurrency * 2
		}
	}
	return Shortcut
}

func (s *shortcut) Prepare(id int, n int) *ShortcutConn {
	address := fmt.Sprintf(shortcutAddress, id)
	conn, existed := s.ports.Get(address) // For specified id, GetOrInsert is not necessary.
	if !existed {
		newConn := NewShortcutConn(n)
		newConn.Address = address
		s.ports.Set(address, newConn)
		return newConn
	} else {
		return conn.(*ShortcutConn)
	}
}

func (s *shortcut) Validate(address string) bool {
	return shortcutRecognizer.MatchString(address)
}

func (s *shortcut) Dial(address string) ([]*mock.Conn, bool)  {
	conn, existed := s.ports.Get(address)
	if !existed {
		return nil, false
	} else {
		return conn.(*ShortcutConn).Conns, true
	}
}

func (s *shortcut) Invalidate(conn *ShortcutConn) {
	s.ports.Del(conn.Address)
}

type ShortcutConn struct {
	Conns   []*mock.Conn
	Client  interface{}
	Address string
}

func NewShortcutConn(n int) *ShortcutConn {
	conn := &ShortcutConn{ Conns: make([]*mock.Conn, n) }
	for i := 0; i < n; i++ {
		conn.Conns[i] = mock.NewConn()
	}
	return conn
}
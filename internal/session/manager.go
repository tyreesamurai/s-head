package session

import (
	"sync"
)

type SessionManager struct {
	sync.RWMutex
	Sessions map[string]*Session
}

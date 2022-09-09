package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	SessionID string
	Login     string
	Created   time.Time
}

type Provider interface {
	SessionInit(sid string) (*Session, error)
	SessionRead(sid string) (*Session, error)
	SessionDestroy(sid string) error
	SessionGC()
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	storage     map[string]*Session
	maxlifetime int64
}

func NewManager(cookieName string, maxlifetime int64) *Manager {
	manager := &Manager{cookieName: cookieName, maxlifetime: maxlifetime, storage: make(map[string]*Session)}
	go manager.SessionGC()
	return manager
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionInit(sid string) (*Session, error) {
	newSession := &Session{SessionID: sid}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if _, ok := manager.storage[sid]; ok {
		return nil, fmt.Errorf("Session with this sid already exists")
	}
	manager.storage[sid] = newSession
	return newSession, nil
}

func (manager *Manager) SessionRead(sid string) (Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if value, ok := manager.storage[sid]; !ok {
		return Session{}, fmt.Errorf("Session with this sid does not exist")
	} else {
		return *value, nil
	}
}

func (manager *Manager) sessionDestroy(sid string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.storage, sid)
}

func (manager *Manager) SessionGC() {
	for {
		time.Sleep(time.Second * time.Duration(manager.maxlifetime))
		manager.lock.Lock()
		defer manager.lock.Lock()
		for key, value := range manager.storage {
			if value.Created.Add(time.Second * time.Duration(manager.maxlifetime)).Before(time.Now()) {
				delete(manager.storage, key)
			}
		}
	}
}

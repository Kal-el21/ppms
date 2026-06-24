package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type OTPSession struct {
	UserID    uint64
	ExpiresAt time.Time
}

type OTPSessionStore struct {
	mu       sync.Mutex
	sessions map[string]OTPSession
}

func NewOTPSessionStore() *OTPSessionStore {
	store := &OTPSessionStore{
		sessions: make(map[string]OTPSession),
	}
	go store.cleanupLoop()
	return store
}

func (s *OTPSessionStore) Create(userID uint64, ttl time.Duration) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	s.mu.Lock()
	s.sessions[token] = OTPSession{UserID: userID, ExpiresAt: time.Now().Add(ttl)}
	s.mu.Unlock()

	return token, nil
}

func (s *OTPSessionStore) Consume(token string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[token]
	if !ok || time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return 0, false
	}

	delete(s.sessions, token) // consumed = single-use
	return session.UserID, true
}

func (s *OTPSessionStore) Peek(token string) (uint64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[token]
	if !ok || time.Now().After(session.ExpiresAt) {
		return 0, false
	}
	return session.UserID, true
}

func (s *OTPSessionStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for k, v := range s.sessions {
			if now.After(v.ExpiresAt) {
				delete(s.sessions, k)
			}
		}
		s.mu.Unlock()
	}
}

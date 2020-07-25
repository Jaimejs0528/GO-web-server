package session

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/user"
	"github.com/google/uuid"
)

// Session represent a user session
type Session struct {
	UserEmail  string
	LastAction time.Time
	mutex      *sync.Mutex
}

var sessionsDB map[string]Session

// MaxSeconds amount in seconds that sessions will durate
const MaxSeconds = 60
const expirationTime = time.Second * MaxSeconds

func init() {
	sessionsDB = make(map[string]Session)
}

// CreateSession creates a new session, returning the sessionID and error
func CreateSession(userEmail string) (string, error) {
	if u := user.GetUser(userEmail); u != nil {
		UUID, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		newSession := Session{userEmail, time.Now(), &sync.Mutex{}}
		sessionsDB[UUID.String()] = newSession
		return UUID.String(), nil
	}
	return "", fmt.Errorf("The user doesn't exist")
}

// GetSession return a session
func GetSession(uuidS string) *Session {
	if session, ok := sessionsDB[uuidS]; ok {
		return &session
	}
	return nil
}

// RenovateSessionTime renovates last action session to Now
func RenovateSessionTime(uuidS string) error {
	if session, ok := sessionsDB[uuidS]; ok {
		session.mutex.Lock()
		session.LastAction = time.Now()
		sessionsDB[uuidS] = session
		session.mutex.Unlock()
		return nil
	}
	return errors.New("Session doesn't exist")
}

// RemoveSession remove a session from DB, return true if success, false if doesn't exist that session
func RemoveSession(uuidS string) bool {
	if session, ok := sessionsDB[uuidS]; ok {
		session.mutex.Lock()
		delete(sessionsDB, uuidS)
		session.mutex.Unlock()
		return true
	}
	return false
}

// CleanExpiredSessions remove all sessions that already have expired
func CleanExpiredSessions() {
	for key, session := range sessionsDB {
		session.mutex.Lock()
		currentTime := time.Now()
		diffTime := session.LastAction.Add(expirationTime).Sub(currentTime)
		if diffTime < 0 {
			fmt.Println(key, session.LastAction, session.UserEmail)
			delete(sessionsDB, key)
		}
		session.mutex.Unlock()
	}
}

// GoCleanExpireSessions remove all sessions that already have expired (Use it with secondary goroutine)
func GoCleanExpireSessions() {
	for {
		CleanExpiredSessions()
		time.Sleep(time.Second * 5)
	}
}

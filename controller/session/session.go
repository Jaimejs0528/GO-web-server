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
}

var sessionsDB map[string]Session
var mutexSession sync.Mutex

// MaxSeconds amount in seconds that sessions will durate
const MaxSeconds = 60
const expirationTime = time.Second * MaxSeconds

func init() {
	sessionsDB = make(map[string]Session)
	mutexSession = sync.Mutex{}
}

// CreateSession creates a new session, returning the sessionID and error
func CreateSession(userEmail string) (string, error) {
	if u := user.GetUser(userEmail); u != nil {
		UUID, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		newSession := Session{userEmail, time.Now()}
		mutexSession.Lock()
		sessionsDB[UUID.String()] = newSession
		mutexSession.Unlock()
		return UUID.String(), nil
	}
	return "", fmt.Errorf("The user doesn't exist")
}

// GetSession return a session
func GetSession(uuidS string) *Session {
	mutexSession.Lock()
	defer mutexSession.Unlock()
	if session, ok := sessionsDB[uuidS]; ok {
		return &session
	}
	return nil
}

// RenovateSessionTime renovates last action session to Now
func RenovateSessionTime(uuidS string) error {
	mutexSession.Lock()
	defer mutexSession.Unlock()
	if session, ok := sessionsDB[uuidS]; ok {
		session.LastAction = time.Now()
		sessionsDB[uuidS] = session
		return nil
	}
	return errors.New("Session doesn't exist")
}

// RemoveSession remove a session from DB, return true if success, false if doesn't exist that session
func RemoveSession(uuidS string) bool {
	mutexSession.Lock()
	defer mutexSession.Unlock()
	if _, ok := sessionsDB[uuidS]; ok {
		delete(sessionsDB, uuidS)
		return true
	}
	return false
}

// CleanExpiredSessions remove all sessions that already have expired
func CleanExpiredSessions() {
	mutexSession.Lock()
	for key, session := range sessionsDB {
		currentTime := time.Now()
		diffTime := session.LastAction.Add(expirationTime).Sub(currentTime)
		if diffTime < 0 {
			fmt.Println(key, session.LastAction, session.UserEmail)
			delete(sessionsDB, key)
		}
	}
	mutexSession.Unlock()
}

// GoCleanExpireSessions remove all sessions that already have expired (Use it with secondary goroutine)
func GoCleanExpireSessions() {
	for {
		CleanExpiredSessions()
		time.Sleep(time.Second * 5)
	}
}

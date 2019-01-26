package cookies

import (
	"sync"
	"github.com/pkg/errors"
	"math/rand"
	"github.com/ByteArena/box2d"
)

type gameSession struct {
	sync.RWMutex
	ID       uint64
	userName string
	score    uint64 // TODO: score should be an atomic counter
	logged   bool
	state    state
	viewport *viewport
	box2dbody *box2d.B2Body
}

type viewport struct {
	x     float32
	y     float32
	xx    float32
	yy    float32
	angle float32
	turbo bool
}

var errUserWasLogged = errors.New("user already logged")
var errCannotSendScreenUpdates = errors.New("session not found")

func newGameSession(id uint64) *gameSession {
	return &gameSession{ID: id, state: &notLoggedState{}, score: 100}
}

func (s *gameSession) getViewport() (*viewport, error) {
	s.RLock()
	defer s.RUnlock()

	if !s.state.canSendScreenUpdates() || s.viewport == nil {
		return nil, errCannotSendScreenUpdates
	}

	return s.viewport, nil
}

func (s *gameSession) updateViewPort(x float32, y float32, xx float32, yy float32, a float32, t bool) {
	s.Lock()
	if s.state.canSendScreenUpdates() {
		s.viewport = &viewport{x: x, y: y, xx: xx, yy: yy, angle: a, turbo: t}
	}
	s.Unlock()
}

func (s *gameSession) login(username string) error {
	s.Lock()
	defer s.Unlock()

	if s.logged {
		return errUserWasLogged
	}
	s.state = &loggedState{}
	s.userName = username
	s.logged = true

	return nil
}

func (s *gameSession) startPlaying() error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.state.(*loggedState); !ok {
		return errors.New("not logged user wants to play")
	}
	s.state = &playingState{}

	return nil
}

func (s *gameSession) getScore() uint64 {
	s.RLock()
	sc := s.score
	s.RUnlock()

	return sc
}


func (s *gameSession) setBox2DBody(b *box2d.B2Body) {
	s.Lock()
	s.box2dbody = b
	s.Unlock()
}

func (s *gameSession) getBox2DBody() *box2d.B2Body {
	s.RLock()
	defer s.RUnlock()
	return s.box2dbody
}




type gameSessions struct {
	sync.RWMutex
	sessions map[uint64]*gameSession
}

func newGameSessions() *gameSessions {
	return &gameSessions{
		sessions: make(map[uint64]*gameSession),
	}
}

func (s *gameSessions) add() uint64 {
	ID := rand.Uint64() << 8// Javascript does not support number larger than 57 bits. Let's avoid problems.
	s.Lock()
	s.sessions[ID] = newGameSession(ID)
	s.Unlock()
	return ID
}

func (s *gameSessions) session(id uint64) *gameSession {
	s.RLock()
	defer s.RUnlock()
	return s.sessions[id]
}

func (s *gameSessions) StartPlaying(ID uint64) error {
	return s.session(ID).startPlaying()
}

func (s *gameSessions) each(fn func(s *gameSession) bool) {
	s.Lock()
	defer s.Unlock()
	for _, session := range s.sessions {
		if !fn(session) {
			return
		}
	}
}

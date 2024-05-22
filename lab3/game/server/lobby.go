package main

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type LobbyState int

const (
	LobbyWaiting LobbyState = iota
	LobbyStarted
)

type Lobby struct {
	Id string

	State LobbyState

	Players           []*User
	PlayersReady      []*User
	LastPlayerReadyAt time.Time

	Game *Game
}

func NewLobby() *Lobby {
	return &Lobby{
		Id:    uuid.New().String(),
		State: LobbyWaiting,
	}
}

func (l *Lobby) IsWaiting() bool {
	return l.State == LobbyWaiting
}

func (l *Lobby) HasPlayer(targetPlayer *User) bool {
	if l.Game.State == GameFinished {
		return false
	}

	for _, player := range l.Game.Players {
		if player.Id == targetPlayer.Id {
			return true
		}
	}

	return false
}

func (l *Lobby) AddPlayer(user *User) error {
	if !l.IsWaiting() {
		return errors.New("lobby is not waiting")
	}

	if l.HasPlayer(user) {
		return errors.New("player is already in lobby")
	}

	l.Players = append(l.Players, user)

	l.Game.AddPlayer(user)

	return nil
}

func (l *Lobby) IsPlayerReady(targetPlayer *User) bool {
	for _, player := range l.PlayersReady {
		if player.Id == targetPlayer.Id {
			return true
		}
	}
	return false
}

func (l *Lobby) MarkPlayerReady(player *User) error {
	if !l.IsWaiting() {
		return errors.New("lobby is not waiting")
	}

	if !l.HasPlayer(player) {
		return errors.New("player is not in lobby")
	}

	if l.IsPlayerReady(player) {
		return errors.New("player is already ready")
	}

	l.PlayersReady = append(l.PlayersReady, player)
	l.LastPlayerReadyAt = time.Now()
	return nil
}

func (l *Lobby) Start() {
	l.State = LobbyStarted
	l.Game.Start()
}

func (l *Lobby) Tick(_ time.Duration) {
	if l.State == LobbyWaiting {
		if len(l.PlayersReady) == len(l.Players) && len(l.Players) >= 3 && l.LastPlayerReadyAt.After(time.Now().Add(-5*time.Second)) {
			l.Start()
			return
		}

		return
	}
}

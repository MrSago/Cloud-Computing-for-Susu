package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/noartem/labs/4/2/cloud/3/game/proto"
	"strings"
	"sync"
	"time"
)

type GameState int

const (
	GameWaiting GameState = iota
	GameStarted
	GameInProgress
	GameFinished
)

type Game struct {
	mutex sync.RWMutex

	Id string

	State     GameState
	StartedAt time.Time

	Players     []*User
	PlayersLost []*User
	PlayerWon   *User

	CurrentPlayer               int
	CurrentPlayerRoundStarted   bool
	CurrentPlayerRoundStartedAt time.Time
	CurrentPlayerRoundTries     int

	UsedAnswers []string
	LastAnswer  string

	Logs []string
}

func NewGame() *Game {
	return &Game{
		Id:            uuid.New().String(),
		State:         GameWaiting,
		CurrentPlayer: -1,
		UsedAnswers:   make([]string, 0),
	}
}

func (g *Game) AddPlayer(player *User) {
	g.Players = append(g.Players, player)
}

func (g *Game) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.State = GameStarted
	g.StartedAt = time.Now()
	Shuffle(g.Players)

	g.log("Все игроки готовы. Игра начнется через 10 секунд")
}

type PlayerGameData struct {
	state proto.PlayerGameState
	logs  []string
}

func (g *Game) GetPlayerGameData(player *User) (*PlayerGameData, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if !g.hasPlayer(player) {
		return nil, errors.New("player not found")
	}

	return &PlayerGameData{
		state: g.getPlayerState(player),
		logs:  g.Logs,
	}, nil
}

func (g *Game) Answer(player *User, answer string) (bool, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if !g.hasPlayer(player) {
		return false, errors.New("player not found")
	}

	if !g.CurrentPlayerRoundStarted {
		return false, errors.New("round not started")
	}

	if g.Players[g.CurrentPlayer].Id != player.Id {
		return false, errors.New("wrong player")
	}

	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)

	if !g.isAnswerCorrect(answer) {
		g.CurrentPlayerRoundTries++

		g.log(fmt.Sprintf("Неподходящий ответ: %s", answer))

		return false, nil
	}

	g.LastAnswer = answer
	g.UsedAnswers = append(g.UsedAnswers, answer)

	g.log(fmt.Sprintf("Подходящий ответ: \"%s\"", answer))

	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
	g.CurrentPlayerRoundStarted = true
	g.CurrentPlayerRoundStartedAt = time.Now()
	g.CurrentPlayerRoundTries = 0

	g.log(fmt.Sprintf("Ход игрока \"%s\"", g.Players[g.CurrentPlayer].Name))

	return true, nil
}

func (g *Game) Tick(_ time.Duration) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.State == GameWaiting || g.State == GameFinished {
		return
	}

	if g.State == GameStarted {
		if time.Now().Add(-10 * time.Second).Before(g.StartedAt) {
			return
		}

		g.State = GameInProgress

		g.log("Игра начинается!")

		g.CurrentPlayer = 0
		g.CurrentPlayerRoundStarted = true
		g.CurrentPlayerRoundStartedAt = time.Now()

		g.log(fmt.Sprintf("Ход игрока \"%s\"", g.Players[g.CurrentPlayer].Name))
	}

	if time.Now().Add(-30*time.Second).After(g.CurrentPlayerRoundStartedAt) || g.CurrentPlayerRoundTries >= 3 {
		g.log(fmt.Sprintf("Игрок \"%s\" не ответил и выбывает из игры", g.Players[g.CurrentPlayer].Name))

		g.PlayersLost = append(g.PlayersLost, g.Players[g.CurrentPlayer])
		g.Players = append(g.Players[:g.CurrentPlayer], g.Players[g.CurrentPlayer+1:]...)

		g.CurrentPlayer = g.CurrentPlayer % len(g.Players)
		g.CurrentPlayerRoundStarted = true
		g.CurrentPlayerRoundStartedAt = time.Now()
		g.CurrentPlayerRoundTries = 0

		g.log(fmt.Sprintf("Ход игрока \"%s\"", g.Players[g.CurrentPlayer].Name))
	}

	if len(g.Players) == 1 {
		g.State = GameFinished
		g.PlayerWon = g.Players[0]
		g.CurrentPlayerRoundStarted = false

		g.log(fmt.Sprintf("Игра закончилась. Игрок \"%s\" выиграл!", g.PlayerWon.Name))
	}
}

func (g *Game) log(message string) {
	g.Logs = append(g.Logs, message)
}

func (g *Game) hasAnswer(answer string) bool {
	for _, usedAnswer := range g.UsedAnswers {
		if usedAnswer == answer {
			return true
		}
	}
	return false
}

func stringStartsWithStringLastRune(a string, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)

	if len(aRunes) == 0 || len(bRunes) == 0 {
		return true
	}

	return aRunes[0] == bRunes[len(bRunes)-1]
}

func (g *Game) isAnswerCorrect(answer string) bool {
	return len(answer) > 0 &&
		stringStartsWithStringLastRune(answer, g.LastAnswer) &&
		!g.hasAnswer(answer)
}

func (g *Game) hasPlayer(player *User) bool {
	for _, p := range g.Players {
		if p.Id == player.Id {
			return true
		}
	}

	for _, p := range g.PlayersLost {
		if p.Id == player.Id {
			return true
		}
	}

	return false
}

func (g *Game) getPlayerState(player *User) proto.PlayerGameState {
	if g.State == GameWaiting {
		return proto.PlayerGameState_NOT_READY
	}

	if g.PlayerWon != nil {
		if g.PlayerWon.Id == player.Id {
			return proto.PlayerGameState_WON
		} else {
			return proto.PlayerGameState_FINISHED
		}
	}

	if g.CurrentPlayer != -1 && g.Players[g.CurrentPlayer].Id == player.Id && g.CurrentPlayerRoundTries < 3 {
		return proto.PlayerGameState_YOUR_TURN
	}

	if g.State == GameInProgress {
		return proto.PlayerGameState_PLAYING
	}

	return proto.PlayerGameState_READY
}

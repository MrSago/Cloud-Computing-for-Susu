package main

import (
	"context"
	"errors"
	"github.com/noartem/labs/4/2/cloud/3/game/proto"
	"strings"
	"sync"
)

type GameServer struct {
	proto.UnimplementedGameServer

	mutex sync.RWMutex

	eventLoop *EventLoop

	users   map[string]*User
	games   map[string]*Game
	lobbies map[string]*Lobby
}

func NewGameServer() *GameServer {
	return &GameServer{
		users:   make(map[string]*User),
		games:   make(map[string]*Game),
		lobbies: make(map[string]*Lobby),
	}
}

func (s *GameServer) StartEventLoop() {
	s.eventLoop = NewEventLoop()
	go s.eventLoop.Start()
}

func (s *GameServer) Register(_ context.Context, in *proto.RegisterRequest) (*proto.RegisterReply, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := s.findOrCreateUser(in.Name)

	return &proto.RegisterReply{
		UserId: user.Id,
	}, nil
}

func (s *GameServer) JoinLobby(_ context.Context, in *proto.JoinLobbyRequest) (*proto.JoinLobbyReply, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	lobby := s.findOrCreateWaitingLobby(in.LobbyId)

	user := s.findUser(in.PlayerId)
	if user == nil {
		return nil, errors.New("user not found")
	}

	userLobby := s.findUserLobby(user)
	if userLobby != nil {
		return nil, errors.New("user already in lobby")
	}

	err := lobby.AddPlayer(user)
	if err != nil {
		return nil, err
	}

	return &proto.JoinLobbyReply{
		LobbyId: lobby.Id,
	}, nil
}

func (s *GameServer) Ready(_ context.Context, in *proto.ReadyRequest) (*proto.ReadyReply, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	lobby := s.findLobby(in.LobbyId)
	if lobby == nil {
		return nil, errors.New("lobby not found")
	}

	user := s.findUser(in.PlayerId)
	if user == nil {
		return nil, errors.New("user not found")
	}

	err := lobby.MarkPlayerReady(user)
	if err != nil {
		return nil, err
	}

	return &proto.ReadyReply{
		GameId: lobby.Game.Id,
	}, nil
}

func (s *GameServer) GetPlayerGameData(_ context.Context, in *proto.GetPlayerGameDataRequest) (*proto.GetPlayerGameDataReply, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := s.findUser(in.PlayerId)
	if user == nil {
		return nil, errors.New("user not found")
	}

	game := s.findGame(in.GameId)
	if game == nil {
		return nil, errors.New("game not found")
	}

	data, err := game.GetPlayerGameData(user)
	if err != nil {
		return nil, err
	}

	return &proto.GetPlayerGameDataReply{
		State: data.state,
		Logs:  strings.Join(data.logs, "\n"),
	}, nil
}

func (s *GameServer) AnswerGame(_ context.Context, in *proto.AnswerGameRequest) (*proto.AnswerGameReply, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := s.findUser(in.PlayerId)
	if user == nil {
		return nil, errors.New("user not found")
	}

	game := s.findGame(in.GameId)
	if game == nil {
		return nil, errors.New("game not found")
	}

	isCorrect, err := game.Answer(user, in.Answer)
	if err != nil {
		return nil, err
	}

	return &proto.AnswerGameReply{
		IsCorrect: isCorrect,
	}, nil
}

func (s *GameServer) findUser(id string) *User {
	return s.users[id]
}

func (s *GameServer) findUserByName(name string) *User {
	for _, user := range s.users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

func (s *GameServer) findOrCreateUser(name string) *User {
	user := s.findUserByName(name)
	if user != nil {
		return user
	}

	return s.createUser(name)
}

func (s *GameServer) createUser(name string) *User {
	user := NewUser(name)
	s.users[user.Id] = user

	return user
}

func (s *GameServer) findLobby(id string) *Lobby {
	return s.lobbies[id]
}

func (s *GameServer) createLobby() *Lobby {
	lobby := NewLobby()
	s.lobbies[lobby.Id] = lobby

	game := NewGame()
	lobby.Game = game
	s.games[game.Id] = game

	s.eventLoop.Add(lobby)
	s.eventLoop.Add(game)

	return lobby
}

func (s *GameServer) findOrCreateWaitingLobby(id string) *Lobby {
	lobby := s.findLobby(id)
	if lobby != nil && lobby.IsWaiting() {
		return lobby
	}

	for _, lobby := range s.lobbies {
		if lobby.IsWaiting() {
			return lobby
		}
	}

	return s.createLobby()
}

func (s *GameServer) findUserLobby(user *User) *Lobby {
	for _, lobby := range s.lobbies {
		if lobby.HasPlayer(user) {
			return lobby
		}
	}

	return nil
}

func (s *GameServer) findGame(id string) *Game {
	return s.games[id]
}

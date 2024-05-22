package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/noartem/labs/4/2/cloud/3/game/proto"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Application struct {
	client  proto.GameClient
	userId  string
	lobbyId string
	gameId  string
}

func NewApplication(client proto.GameClient) (*Application, error) {
	fmt.Print("> Ведите имя пользователя (Anon): ")
	r := bufio.NewReader(os.Stdin)
	username, err := r.ReadString('\n')
	if err != nil {
		return nil, errors.New("ошибка при вводе имени")
	}
	username = strings.TrimSpace(username)
	if username == "" {
		username = "Anon"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	registerReply, err := client.Register(ctx, &proto.RegisterRequest{
		Name: username,
	})
	if err != nil {
		return nil, err
	}

	return &Application{
		client:  client,
		userId:  registerReply.UserId,
		lobbyId: "",
		gameId:  "",
	}, nil
}

func (a *Application) Join(lobbyId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	joinResponse, err := a.client.JoinLobby(ctx, &proto.JoinLobbyRequest{
		PlayerId: a.userId,
		LobbyId:  lobbyId,
	})
	if err != nil {
		return err
	}

	a.lobbyId = joinResponse.LobbyId

	return nil
}

func (a *Application) Ready() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	readyResponse, err := a.client.Ready(ctx, &proto.ReadyRequest{
		PlayerId: a.userId,
		LobbyId:  a.lobbyId,
	})
	if err != nil {
		return err
	}

	a.gameId = readyResponse.GameId

	return nil
}

func (a *Application) getData() (*proto.GetPlayerGameDataReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	data, err := a.client.GetPlayerGameData(ctx, &proto.GetPlayerGameDataRequest{
		PlayerId: a.userId,
		GameId:   a.gameId,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readStringWithTimeout(timeout time.Duration) (string, error) {
	dataStream := make(chan string, 1)

	go func() {
		r := bufio.NewReader(os.Stdin)
		data, _ := r.ReadString('\n')
		dataStream <- data
		close(dataStream)
	}()

	select {
	case res := <-dataStream:
		res = strings.TrimSpace(res)
		return res, nil
	case <-time.After(timeout):
		return "", errors.New("Время ожидания ответа истекло")
	}
}

func (a *Application) answer() (bool, error) {
	tries := 0

	fmt.Println(". Ваш ход")
	for {
		fmt.Println("> Введите название города: ")
		answer, err := readStringWithTimeout(time.Second * 30)
		if err != nil {
			return false, nil
		}
		answer = strings.TrimSpace(answer)

		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		res, err := a.client.AnswerGame(ctx, &proto.AnswerGameRequest{
			PlayerId: a.userId,
			GameId:   a.gameId,
			Answer:   answer,
		})
		if err != nil {
			return false, err
		}

		tries++

		if res.IsCorrect {
			fmt.Println(". Хороший ответ")
			return true, nil
		}

		if tries >= 3 {
			fmt.Println(". Исчерпано максимальное количество попыток")
			return false, nil
		}

		fmt.Println(". Неправильный ответ. Попробуйте ещё раз")
	}
}

func (a *Application) Play() error {
	lastLogs := make([]string, 0)

	for {
		data, err := a.getData()
		if err != nil {
			return err
		}

		if data.Logs != "" {
			logs := strings.Split(strings.TrimSpace(data.Logs), "\n")
			if len(logs) > len(lastLogs) {
				for i := len(lastLogs); i < len(logs); i++ {
					fmt.Printf("@ %s\n", logs[i])
				}
				lastLogs = logs
			}
		}

		if data.State == proto.PlayerGameState_YOUR_TURN {
			answered, err := a.answer()
			if err != nil {
				return err
			}

			if !answered {
				fmt.Println(". Вы не ответили и проиграли")
				return nil
			}
		}

		if data.State == proto.PlayerGameState_WON {
			fmt.Println(". Вы победили")
			return nil
		}

		if data.State == proto.PlayerGameState_FINISHED {
			fmt.Println(". Вы проиграли")
			return nil
		}

		time.Sleep(time.Millisecond * 1500)
	}
}

func main() {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("> Ведите адрес сервера (localhost:50051): ")
	address, err := r.ReadString('\n')
	if err != nil {
		fmt.Printf("! Ошибка при вводе адреса: %v\n", err)
		os.Exit(1)
	}
	address = strings.TrimSpace(address)
	if address == "" {
		address = "localhost:50051"
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("! Ошибка при подключении к серверу: %v\n", err)
		os.Exit(1)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("! Ошибка при закрытии соединения: %v\n", err)
			os.Exit(1)
		}
	}(conn)

	c := proto.NewGameClient(conn)

	app, err := NewApplication(c)
	if err != nil {
		fmt.Printf("! Ошибка при создании приложения: %v\n", err)
		os.Exit(1)
	}

	for {
		fmt.Print("> ")
		command, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("! Ошибка при вводе команды: %v\n", err)
			continue
		}
		command = strings.TrimSpace(command)

		switch command {
		case "help":
			fmt.Println("\tjoin  - вход в лобби")
			fmt.Println("\tready - готовность к игре")
			fmt.Println("\tquit  - выход из игры")
			fmt.Println("\thelp  - помощь")
			break
		case "join":
			if app.lobbyId != "" {
				fmt.Println("! Вы уже подключились к лобби")
				break
			}

			fmt.Print("> ")
			lobbyId, err := r.ReadString('\n')
			if err != nil {
				fmt.Printf("! Ошибка при вводе id лобби: %v\n", err)
				continue
			}
			lobbyId = strings.TrimSpace(lobbyId)

			err = app.Join(lobbyId)
			if err != nil {
				fmt.Printf("! Ошибка при входе в лобби: %v\n", err)
				break
			}
			fmt.Println(". Вы успешно подключились к лобби " + app.lobbyId)
			break
		case "ready":
			if app.lobbyId == "" {
				fmt.Println("! Вы не подключились к лобби")
				break
			}

			if app.gameId != "" {
				fmt.Println("! Вы уже готовы к игре")
				break
			}

			err = app.Ready()
			if err != nil {
				fmt.Printf("! Ошибка при готовности к игре: %v\n", err)
				break
			}
			fmt.Println(". Вы успешно подключились к игре")

			fmt.Println(". Ожидаем начала игры")
			err := app.Play()
			app.lobbyId = ""
			app.gameId = ""
			if err != nil {
				fmt.Printf("! Ошибка при игре: %v\n", err.Error())
				break
			}
			break
		case "quit":
			fmt.Println(". До свидания")
			return
		case "":
			break
		default:
			fmt.Println(". Неизвестная команда")
			break
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/moul/tictactoe"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type QuestionMessage struct {
	GameID      string
	Action      string
	Game        string
	Players     int
	Board       map[string]string
	You         string
	PlayerIndex int
}

func main() {
	logrus.Warnf("%s << %v", os.Args[0], os.Args[1])

	var question QuestionMessage
	if err := json.Unmarshal([]byte(os.Args[1]), &question); err != nil {
		logrus.Fatalf("%s XX err: %v", err)
	}

	switch question.Action {
	case "init":
		logrus.Warnf("init: %v", question)
		fmt.Println("{\"name\":\"moul-tictactoe\"}")
	case "play-turn":
		ttt := tictactoe.NewTicTacToe()
		ttt.Set(0, 0, question.Board["0-0"])
		ttt.Set(0, 1, question.Board["0-1"])
		ttt.Set(0, 2, question.Board["0-2"])
		ttt.Set(1, 0, question.Board["1-0"])
		ttt.Set(1, 1, question.Board["1-1"])
		ttt.Set(1, 2, question.Board["1-2"])
		ttt.Set(2, 0, question.Board["2-0"])
		ttt.Set(2, 1, question.Board["2-1"])
		ttt.Set(2, 2, question.Board["2-2"])
		ttt.SetPlayer(question.You)

		logrus.Debugf("map: %q", ttt.ShowMap())

		next, err := ttt.Next()
		if err != nil {
			logrus.Fatalf("Error while getting the next piece: %v", err)
		}

		fmt.Printf("{\"play\":\"%d-%d\"}", next.Y, next.X)
	default:
		logrus.Fatalf("Unknown action: %q", question.Action)
	}
}

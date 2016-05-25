package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/moul/bolosseum/bots"
	"github.com/moul/tictactoe"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func NewTictactoeBot() *TictactoeBot {
	return &TictactoeBot{}
}

type TictactoeBot struct{}

func (b *TictactoeBot) Init(message bots.QuestionMessage) *bots.ReplyMessage {
	// FIXME: init ttt here
	return &bots.ReplyMessage{
		Name: "moul-tictactoe",
	}
}

func (b *TictactoeBot) PlayTurn(question bots.QuestionMessage) *bots.ReplyMessage {
	ttt := tictactoe.NewTicTacToe()

	board := question.Board.(map[string]interface{})

	ttt.Set(0, 0, board["0-0"].(string))
	ttt.Set(0, 1, board["0-1"].(string))
	ttt.Set(0, 2, board["0-2"].(string))
	ttt.Set(1, 0, board["1-0"].(string))
	ttt.Set(1, 1, board["1-1"].(string))
	ttt.Set(1, 2, board["1-2"].(string))
	ttt.Set(2, 0, board["2-0"].(string))
	ttt.Set(2, 1, board["2-1"].(string))
	ttt.Set(2, 2, board["2-2"].(string))
	ttt.SetPlayer(question.You.(string))

	logrus.Debugf("map: %q", ttt.ShowMap())

	next, err := ttt.Next()
	if err != nil {
		logrus.Fatalf("Error while getting the next piece: %v", err)
	}

	return &bots.ReplyMessage{
		Play: fmt.Sprintf("%d-%d", next.Y, next.X),
	}
}

func main() {
	if len(os.Args) == 1 {
		// web mode
		logrus.Warnf("You ran this program without argument, it will then start a web server")
		logrus.Warnf("usage: ")
		logrus.Warnf("- %s            # web mode", os.Args[0])
		logrus.Warnf("- %s some-json  # cli mode", os.Args[0])

		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			if message := c.Query("message"); message != "" {
				bot := NewTictactoeBot()

				logrus.Warnf("<< %s", message)
				var question bots.QuestionMessage
				if err := json.Unmarshal([]byte(message), &question); err != nil {
					c.JSON(500, gin.H{"Error": err})
					return
				}

				reply := &bots.ReplyMessage{}
				switch question.Action {
				case "init":
					reply = bot.Init(question)
				case "play-turn":
					reply = bot.PlayTurn(question)
				default:
					// FIXME: reply message error
					c.JSON(500, gin.H{"Error": fmt.Errorf("Unknown action: %q", question.Action)})
					return
				}

				c.JSON(200, reply)
			} else {
				c.String(404, "This server is a bot for bolosseum.")
			}
		})
		r.Run(":8080")
	} else {
		// cli mode
		logrus.Warnf("%s << %v", os.Args[0], os.Args[1])

		var question bots.QuestionMessage
		if err := json.Unmarshal([]byte(os.Args[1]), &question); err != nil {
			logrus.Fatalf("%s XX err: %v", err)
		}

		bot := NewTictactoeBot()

		reply := &bots.ReplyMessage{}
		switch question.Action {
		case "init":
			reply = bot.Init(question)
		case "play-turn":
			reply = bot.PlayTurn(question)
		default:
			// FIXME: reply message error
			logrus.Fatalf("Unknown action: %q", question.Action)
		}

		jsonString, err := json.Marshal(reply)
		if err != nil {
			logrus.Fatalf("Failed to marshal json: %v", err)
		}

		fmt.Println(string(jsonString))
	}
}

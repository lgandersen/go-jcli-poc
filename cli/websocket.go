package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func Dial(endpoint string) (chan struct{}, chan os.Signal, *websocket.Conn) {
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})
	return done, interrupt, ws
}

func AwaitDoneOrUserInterrupt(done chan struct{}, interrupt chan os.Signal, ws *websocket.Conn) {
	defer ws.Close()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Interrupted by user")
			TryGracefulWSDisconnectconnect(done, ws)
			return
		}
	}
}

func ListenForWSMessages(done chan struct{}, ws *websocket.Conn) {
	defer close(done)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if strings.HasPrefix(err.Error(), succesful_ws_exit) {
				fmt.Println(err.Error()[len(succesful_ws_exit):])
			} else {
				fmt.Println("websocket closed unexpectedly:", err.Error())
			}
			return
		}
		msg := string(message)
		if msg[:3] == "ok:" {
			// First message receieved when the ws is succesfully established.
			continue
		}
		if msg[:3] == "io:" {
			fmt.Print(msg[3:])
		}
	}
}

func TryGracefulWSDisconnectconnect(done chan struct{}, ws *websocket.Conn) {
	_ = ws.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}

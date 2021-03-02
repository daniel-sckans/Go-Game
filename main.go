package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	port := os.Getenv("PORT")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	currentBoxes := NewBoxes{
		X1: fmt.Sprintf("%f", r1.Float64()*.3),
		Y1: fmt.Sprintf("%f", r1.Float64()*.8),
		W1: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
		H1: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
		C:  fmt.Sprintf("%f", r1.Float64()),
		X2: fmt.Sprintf("%f", r1.Float64()*(.8-.5)+.5),
		Y2: fmt.Sprintf("%f", r1.Float64()*.8),
		W2: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
		H2: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
	}

	http.Handle("/", http.FileServer(http.Dir("./site")))

	log.Print("Running on: " + port)

	const (
		LEAVE = iota
		JOIN
		GREEN
		RED
	)
	type wsMsg struct {
		Change int
		Conn   *websocket.Conn
	}
	wsStateChange := make(chan wsMsg)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		defer func() {
			wsStateChange <- wsMsg{LEAVE, conn}
		}()
		wsStateChange <- wsMsg{JOIN, conn}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			if string(msg) == "loaded" {
				fmt.Println("User Connected")
			}
			if string(msg) == "green" {
				println("New Boxes")
				wsStateChange <- wsMsg{GREEN, conn}
			}
			if string(msg) == "red" {
				fmt.Println("Bad Choice")
				wsStateChange <- wsMsg{RED, conn}
			}
		}
	})

	go func() {
		connections := []*websocket.Conn{}
		for {
			select {
			case <-time.Tick(5 * time.Second):
				for _, conn := range connections {
					conn.WriteJSON(currentBoxes)
				}
			case wsc := <-wsStateChange:
				switch wsc.Change {
				case LEAVE:
					for i, conn := range connections {
						if conn == wsc.Conn {
							connections = append(connections[:i], connections[i+1:]...)
						}
					}
				case JOIN:
					connections = append(connections, wsc.Conn)
					wsc.Conn.WriteJSON(currentBoxes)
				case GREEN:
					wsc.Conn.WriteJSON(SingleMessage{Message: "+1"})
					currentBoxes = NewBoxes{
						X1: fmt.Sprintf("%f", r1.Float64()*.3),
						Y1: fmt.Sprintf("%f", r1.Float64()*.8),
						W1: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
						H1: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
						C:  fmt.Sprintf("%f", r1.Float64()),
						X2: fmt.Sprintf("%f", r1.Float64()*(.8-.5)+.5),
						Y2: fmt.Sprintf("%f", r1.Float64()*.8),
						W2: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
						H2: fmt.Sprintf("%f", r1.Float64()*(.2-.1)+.1),
					}
					for _, conn := range connections {
						conn.WriteJSON(currentBoxes)
					}
				case RED:
					wsc.Conn.WriteJSON(SingleMessage{Message: "-5"})
				}
			}
		}
	}()

	http.ListenAndServe(":"+port, nil)
}

type NewBoxes struct {
	X1 string `json:"x1"`
	Y1 string `json:"y1"`
	W1 string `json:"w1"`
	H1 string `json:"h1"`
	C  string `json:"c"`
	X2 string `json:"x2"`
	Y2 string `json:"y2"`
	W2 string `json:"w2"`
	H2 string `json:"h2"`
}

type SingleMessage struct {
	Message string `json:"message"`
}

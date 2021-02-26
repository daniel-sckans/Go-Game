package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
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

	fmt.Println("Running")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			ch := time.Tick(time.Millisecond)

			for range ch {
				_, msg, err := conn.ReadMessage()

				if err != nil {
					fmt.Println(err)
					return
				}

				if string(msg) == "loaded" {
					fmt.Println("User Connected")
					conn.WriteJSON(currentBoxes)
				}
				if string(msg) == "green" {
					println("New Boxes")
					err := conn.WriteJSON(SingleMessage{
						Message: "+1",
					})

					if err != nil {
						fmt.Println(err)
						return
					}

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
				} else if string(msg) == "red" {
					err := conn.WriteJSON(SingleMessage{
						Message: "-5",
					})

					if err != nil {
						fmt.Println(err)
						return
					}
				}
				conn.WriteJSON(currentBoxes)
			}
		}(conn)
	})

	http.ListenAndServe(":8080", nil)
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

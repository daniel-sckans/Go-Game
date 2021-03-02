package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
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

	mtx := &sync.Mutex{}
	connections := []*websocket.Conn{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		connections = append(connections, conn)
		defer func() {
			for i, c := range connections {
				if c == conn {
					connections = append(connections[:i], connections[i:]...)

					// OR REPLACE THE DEAD CONNECTION WITH THE LAST ON THE LIST
					// connections[i] = connections[len(connections)-1]
					// connections = connections[:len(connections)-1]
				}
			}
		}()

		for {
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
				// conn.WriteJSON(SingleMessage{
				// 	Message: "+1",
				// })
				mtx.Lock()
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
				mtx.Unlock()
			} else if string(msg) == "red" {
				// conn.WriteJSON(SingleMessage{
				// 	Message: "-5",
				// })
			}
		}
	})

	go func() {
		ch := time.Tick(time.Millisecond)

		for range ch {
			for _, conn := range connections {
				conn.WriteJSON(currentBoxes)
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

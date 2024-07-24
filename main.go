package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/ws", handleWsConnections)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	setCookie(w)

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, nil)
}

func setCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "username",
		Value:    "secret user",
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func handleWsConnections(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		fmt.Printf("Cookie name: %s, value: %s\n", cookie.Name, cookie.Value)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrade error: %v\n", err)
		return
	}
	defer ws.Close()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		var data = fmt.Sprintf("test response. Cookie server has received from browser: %v", cookies)
		err = ws.WriteMessage(messageType, []byte(data))
		if err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}
}

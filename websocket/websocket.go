package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't Receive")
			break
		}

		fmt.Println("Received back from client:" + reply)
		msg := "Received:" + reply

		fmt.Println("Sending to client:" + msg)
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}

}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("chat.html")
	t.Execute(w, t)
}

func main() {
	//http.Handle("/", websocket.Handler(Echo))
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Handler(websocket.Handler(Echo))
	router.HandleFunc("/chat", http.HandlerFunc(chatHandler))
	if err := http.ListenAndServe(":1234", router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type msg struct {
	msg string
}

func main() {
	http.HandleFunc("/echo", wsHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	panic(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Fora da origem", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Erro ao abrir conexao", http.StatusBadRequest)
	}

	conn.SetCloseHandler(closeHandler)

	go echo(conn)
}

func closeHandler(code int, text string) error {
	fmt.Println("Cliente desconectou")
	return errors.New("")
}

func echo(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			conn.Close()
			fmt.Println("Erro ao ler json.", err)
			return
		}

		fmt.Printf("Mensagem: %#v\n", m)

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}

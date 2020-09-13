package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = newHub()

func main() {

	flag.Parse()

	go hub.run()

	http.HandleFunc("/", serveHome)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/create", createRoomHandler)
	http.HandleFunc("/join/{room-id}", joinRoomHandler)


	fmt.Println("starting...")
	http.ListenAndServe(*addr, nil)

}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomId := createRoom(hub, w, r)
	fmt.Println("Created room: " + roomId)

	// build response
	d := map[string]string{"room-id": roomId}
	response, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	//http.Redirect(w, r, "http://" + "localhost" + *addr + "/join" + "/"+roomId, http.StatusCreated)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["room_id"]
	if !roomExists(hub, roomId) {
		fmt.Println("room does not exist")
		http.Redirect(w,r,"http://"+r.Host+r.URL.String()+"/"+roomId, http.StatusNotFound)
	}

	// join the room
	joinRoom(hub, roomId)

	http.Redirect(w,r,"http://"+r.Host+r.URL.String()+"/"+vars["room_id"], http.StatusOK)
}
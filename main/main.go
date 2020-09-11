package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var hub = newHub()

func main() {


	//hub := newHub()
	go hub.run()

	var dir string
	flag.StringVar(&dir, "dir", ".", "directory to serve files from.  defaults to current directory")

	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})



	r := mux.NewRouter()

	//http.Handle("/", r)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))




	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("starting...")
	r.HandleFunc("/create", createRoomHandler).Methods("POST").Schemes("http")
	r.HandleFunc("/join/{room_id}", joinRoomHandler).Methods("GET").Schemes("http")

	go srv.ListenAndServe()
	http.ListenAndServe(*addr, r)
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomId := createRoom(hub, w, r)
	http.Redirect(w, r, "http://"+r.Host+r.URL.String()+"/"+roomId, http.StatusMovedPermanently)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO check if room exists
	// TODO retrieve client
	//if _, ok := hub.rooms[vars["room_id"]]; !ok {
	//	http.Redirect(w,r,"http://"+r.Host+r.URL.String()+"/"+vars["room_id"], http.StatusNotFound)
	//}

	// otherwise join the room
	http.Redirect(w,r,"http://"+r.Host+r.URL.String()+"/"+vars["room_id"], http.StatusNotFound)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))


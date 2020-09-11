package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := newHub()
	go hub.run()

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))


	r := mux.NewRouter()
	fmt.Printf("starting...")

	r.HandleFunc("/create", create_room_handler).Methods("POST").Schemes("http")
	r.HandleFunc("/join/{room_id}", join_room_handler).Methods("GET").Schemes("http")


	http.ListenAndServe(":80", r)
}


func create_room_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// create new room
	// TODO add dynamoDB library

	// load room information + messages into memory


	// start websocket session


	fmt.Fprintf(w, "You hit the create page motherPUCKER")
}

func join_room_handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	// check if room exists
	// TODO add dynamoDB library
	// start websocket session


	fmt.Fprintf(w, "Joining room: %v\n", vars["room_id"])
}

// echo request back to client
func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO take out for web deployment
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
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


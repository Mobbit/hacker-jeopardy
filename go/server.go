package  main


import (
	"net/http"
	"log"
	"time"
)

import (
	"jeopardy/statemachine"
)

import "code.google.com/p/go.net/websocket"

const assetDir    = "./static"
const q_fn        = "./questions/questions.json"


type LogHandler struct {
	hdlr http.Handler
}

type ResponseWriter struct {
	rw http.ResponseWriter
	length int
	status int
}

func (r *ResponseWriter) Header() http.Header{
	return r.rw.Header()
}
func (r *ResponseWriter) Write(b []byte) (i int, e error) {
	if r.status == 0 {
		r.status = 200
	}
	i,e = r.rw.Write(b)

	r.length += i
	return
}
func (r *ResponseWriter) WriteHeader(i int) {
	r.status = i
	r.rw.WriteHeader(i)
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	rw := ResponseWriter{rw: w}
	l.hdlr.ServeHTTP(&rw, r)

	//127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
	// or so says wikipedia. grmpf.

	log.Printf("%s %s %s [%s] \"%s %s %s\", %d %d",
		r.RemoteAddr, "-", "-", time.Now().Format(time.RFC3339),
		r.Method, r.RequestURI, r.Proto, rw.status, rw.length)
}

func MakeLogging (h http.Handler)http.Handler {
	return &LogHandler{h}
}

var ws statemachine.WebsocketHandler
func WSHandle(con *websocket.Conn) {
	ws = statemachine.WebsocketHandler{}
	ws.SetSocket(con)
	game.Admin.Prompt("Websocket connected")
	game.UI = &ws

	con.Write([]byte("{\"hello\":\"hello\"}"))
	for {
		b := []byte("bla")
		con.Read(b)
	}
}

var game *statemachine.Game

func main (){
	admin := statemachine.Admin{} 
	game  = statemachine.NewGame(q_fn, &admin) 
	http.Handle("/", MakeLogging(http.FileServer(http.Dir(assetDir))))
	http.Handle("/ws/", websocket.Handler(WSHandle))
	http.ListenAndServe(":9090", nil)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/negroni"
)

var (
	config   configuration
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
)

type configuration struct {
	Debug         bool   `default:"true"`
	Scheme        string `default:"http"`
	ListenAddress string `default:":8080"`
	PrivateKey    string `default:"ssl/server.key"`
	Certificate   string `default:"ssl/server.pem"`
}

type httpErr struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func handleErr(w http.ResponseWriter, err error, status int) {
	msg, err := json.Marshal(&httpErr{
		Msg:  err.Error(),
		Code: status,
	})
	if err != nil {
		msg = []byte(err.Error())
	}
	http.Error(w, string(msg), status)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleErr(w, err, http.StatusInternalServerError)
		return
	}
	defer c.Close()
	// fetch and send the slots from cache
	// ...
	for {
		// subscribe to cache change events and forward to client if somethnig cahnged
	}
	// time.AfterFunc(20*time.Second, func() {
	// 	err = c.WriteMessage(websocket.TextMessage, []byte(`hii from me...`))
	// 	if err != nil {
	// 		handleErr(w, err, http.StatusInternalServerError)
	// 	}
	// })
	// // fetch the slots and write to client
	// err = c.WriteMessage(websocket.TextMessage, []byte(`mssadfg`))
	// if err != nil {
	// 	handleErr(w, err, http.StatusInternalServerError)
	// }
	// for {
	// 	mt, msg, err := c.ReadMessage()
	// 	if err != nil {
	// 		handleErr(w, err, http.StatusInternalServerError)
	// 		break
	// 	}
	// 	if mt != websocket.TextMessage {
	// 		handleErr(w, errors.New("Only text message are supported"), http.StatusNotImplemented)
	// 		break
	// 	}
	// 	var v []byte
	// 	json.Unmarshal(msg, &v)
	// 	err = c.WriteMessage(mt, []byte(msg))
	// 	if err != nil {
	// 		handleErr(w, err, http.StatusInternalServerError)
	// 		break
	// 	}
	// }
}

func main() {

	// Default values
	err := envconfig.Process("SOCKETCAM", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	if config.Debug {
		log.Printf("==> SCHEME: %v", config.Scheme)
		log.Printf("==> ADDRESS: %v", config.ListenAddress)
		log.Printf("==> PRIVATEKEY: %v", config.PrivateKey)
		log.Printf("==> CERTIFICATE: %v", config.Certificate)
	}

	router := newRouter()
	n := negroni.Classic()

	n.UseHandler(router)
	if config.Scheme == "https" {
		log.Fatal(http.ListenAndServeTLS(config.ListenAddress, config.Certificate, config.PrivateKey, n))

	} else {
		log.Fatal(http.ListenAndServe(config.ListenAddress, n))

	}
}

// NewRouter is the constructor for all my routes
func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("GET").
		Path("/ws").
		Name("Communication Channel").
		HandlerFunc(serveWs)

	router.
		Methods("GET").
		PathPrefix("/").
		Name("Static").
		Handler(http.FileServer(http.Dir("./htdocs")))
	return router
}

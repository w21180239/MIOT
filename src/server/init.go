package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	formatter *render.Render
)

func init() {
	formatter = render.New(render.Options{
		Directory:     "views",
		IndentJSON:    true,
		UnEscapeHTML:  true,
		IsDevelopment: true,
	})
}

// NewServer create a new Negroni Server
func NewServer() *negroni.Negroni {
	n := negroni.Classic()
	router := mux.NewRouter()
	initRouter(router)
	n.UseHandler(router)

	return n
}

func initRouter(router *mux.Router) {
	router.HandleFunc("/miot/adapter", handlerProtocol()).Methods(http.MethodPost)
}

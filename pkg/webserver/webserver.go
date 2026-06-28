package webserver

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type method string

var POST method = "POST"
var GET method = "GET"

type WebServer struct {
	Router        chi.Router
	Handlers      []Handler
	WebServerPort string
}

type Handler struct {
	Method  method
	Handler http.HandlerFunc
	Path    string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make([]Handler, 0),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, m method, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, Handler{
		Method:  m,
		Handler: handler,
		Path:    path,
	})
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Timeout(60 * time.Second))
	s.Router.Use(middleware.Logger)
	for _, wsh := range s.Handlers {
		// otelhttp gera um span por request, nomeado com a rota, e lê o
		// traceparent dos headers para continuar traces distribuidos.
		handler := otelhttp.NewHandler(wsh.Handler, wsh.Path)
		switch wsh.Method {
		case POST:
			s.Router.Method(http.MethodPost, wsh.Path, handler)
		case GET:
			s.Router.Method(http.MethodGet, wsh.Path, handler)
		}
	}
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityalstkp/gov8react/internal/constants"
	"github.com/adityalstkp/gov8react/internal/handler"
	"github.com/adityalstkp/gov8react/internal/usecase"
	"github.com/adityalstkp/gov8react/internal/utilities"
	"github.com/go-chi/chi/v5"
)

var httpAddr string
var withHydration bool

func init() {
	flag.StringVar(&httpAddr, "http_addr", "0.0.0.0:3000", "http listen address")
	flag.BoolVar(&withHydration, "with_hydration", false, "render with hydration")
	flag.Parse()
}

func main() {
	v8Ctx, err := newV8Ctx()
	if err != nil {
		log.Panicln("error init v8 context", err)
	}

	reactHtmlPath := fmt.Sprintf("%s/react.tmpl", constants.BASE_TEMPLATE_DIR)
	tmpl, err := utilities.CreateTemplate(reactHtmlPath, "react.html")
	if err != nil {
		log.Panicln("error create template", err)
	}

	introUsecase := usecase.NewIntroUsecase()
	reactHandler := handler.NewReactHandler(handler.ReactHandlerOpts{
		V8Ctx:         v8Ctx,
		Tmpl:          tmpl,
		WithHydration: withHydration,
		IntroUsecase:  introUsecase,
	})
	introHandler := handler.NewIntroHandler(handler.IntroHandlerOpts{IntroUsecase: introUsecase})

	server := http.Server{Addr: httpAddr,
		Handler: NewHandler(httpHandler{
			reactHandler: reactHandler,
			introHandler: introHandler}),
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancelCtx := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				cancelCtx()
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Println("Listen on", httpAddr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

type httpHandler struct {
	reactHandler handler.ReactHandlerRouter
	introHandler handler.IntroHandlerRouter
}

func NewHandler(h httpHandler) http.Handler {
	r := chi.NewRouter()

	staticFs := http.FileServer(http.Dir(constants.BASE_ARTIFACTS_DIR))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFs))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/intro", h.introHandler.Greet)
	})
	r.Get("/*", h.reactHandler.RenderReact)

	return r
}

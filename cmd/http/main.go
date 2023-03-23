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
	"github.com/adityalstkp/gov8react/internal/utilities"
	"github.com/go-chi/chi/v5"
)

var httpAddr string

func init() {
	flag.StringVar(&httpAddr, "http_addr", "0.0.0.0:3000", "http listen address")
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
	reactHandler := handler.NewReactHandler(handler.ReactHandlerOpts{
		V8Ctx: v8Ctx,
		Tmpl:  tmpl,
	})

	server := http.Server{Addr: httpAddr, Handler: NewHandler(httpHandler{reactHandler: reactHandler})}
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
}

func NewHandler(h httpHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/*", h.reactHandler.RenderReact)

	return r
}

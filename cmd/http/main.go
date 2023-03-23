package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/adityalstkp/gov8react/internal/constants"
	"github.com/adityalstkp/gov8react/internal/handler"
	"github.com/go-chi/chi/v5"
	v8 "rogchap.com/v8go"
)

var httpAddr string

func init() {
	flag.StringVar(&httpAddr, "http_addr", "0.0.0.0:3000", "http listen address")
	flag.Parse()
}

func main() {
	iso := v8.NewIsolate()

	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v", info.Args())
		return nil
	})

	procEnv := v8.NewObjectTemplate(iso)
	process := v8.NewObjectTemplate(iso)
	process.Set("env", procEnv)
	process.Set("version", "gov8")

	global := v8.NewObjectTemplate(iso)
	global.Set("print", printfn)
	global.Set("process", process)
	global.Set("GO_APP", v8.NewObjectTemplate(iso))

	v8Ctx := v8.NewContext(iso, global)

	var polyfills = []struct {
		name string
	}{{name: "text_encoder"}, {name: "buffer"}}

	for _, p := range polyfills {
		pName := fmt.Sprintf("%s/polyfills.%s.js", constants.BASE_ARTIFACTS_DIR, p.name)
		pB, err := ioutil.ReadFile(pName)
		if err != nil {
			log.Panicf("error read js polyfill %s", err.Error())
		}

		pO := fmt.Sprintf("polyfills_%s.js", p.name)
		_, err = v8Ctx.RunScript(string(pB), pO)
		if err != nil {
			log.Panicln("error in bundling polyfill", err)
		}
	}

	appBundlePath := fmt.Sprintf("%s/main.js", constants.BASE_ARTIFACTS_DIR)
	appBundle, err := ioutil.ReadFile(appBundlePath)
	if err != nil {
		log.Panicf("error read js app bundle %s", err.Error())
	}

	_, err = v8Ctx.RunScript(string(appBundle), "bundle.js")
	if err != nil {
		e := err.(*v8.JSError)    // JavaScript errors will be returned as the JSError struct
		log.Println(e.StackTrace) // the full stack trace of the error, if available
		log.Panicln("error in bundling app")
	}

	_, err = v8Ctx.RunScript(`
    function renderReact(arg) {
        return GO_APP.render(arg);
    }
    function runMatchRoutes(url) {
        return GO_APP.getMatchRoutes(url);
    }
    `, "register_main.js")
	if err != nil {
		log.Panicln("error in registering app func", err)
	}

	reactHtmlPath := fmt.Sprintf("%s/react.tmpl", constants.BASE_TEMPLATE_DIR)
	reactHtml, err := ioutil.ReadFile(reactHtmlPath)
	if err != nil {
		log.Panicf("error read react html %s", err.Error())
	}

	tmpl := template.Must(template.New("react.html").Parse(string(reactHtml)))
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

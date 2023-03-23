package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	v8 "rogchap.com/v8go"
)

type reactHandler struct {
	v8Ctx *v8.Context
	tmpl  *template.Template
}

type ReactHandlerOpts struct {
	V8Ctx *v8.Context
	Tmpl  *template.Template
}

func NewReactHandler(opts ReactHandlerOpts) *reactHandler {
	return &reactHandler{
		v8Ctx: opts.V8Ctx,
		tmpl:  opts.Tmpl,
	}
}

type ReactHandlerRouter interface {
	RenderReact(w http.ResponseWriter, r *http.Request)
}

type routeMatch struct {
	Params       interface{} `json:"params"`
	Pathname     string      `json:"pathname"`
	PathnameBase string      `json:"pathnameBase"`
}

type markupValue struct {
	ReactMarkup string `json:"markup"`
}

func (rH *reactHandler) RenderReact(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.String()

	runMatchRoutes := fmt.Sprintf(`
    runMatchRoutes("%s");
    `, reqUrl)
	match, err := rH.v8Ctx.RunScript(runMatchRoutes, "match_routes.js")
	if err != nil {
		log.Println("error run match_routes.js", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !match.IsArray() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mJson, err := match.MarshalJSON()
	if err != nil {
		log.Println("cannot marshall route match")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rM []routeMatch
	err = json.Unmarshal(mJson, &rM)
	if err != nil {
		log.Println("cannot umarshall route match")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var staticData map[string]interface{}
	matchRoute := rM[0].Pathname
	if matchRoute == "/" {
		staticData = make(map[string]interface{})
		staticData["greet"] = r.Header.Get("user-agent")
	}

	reactAppArgs := map[string]interface{}{
		"url":        reqUrl,
		"staticData": staticData,
	}
	appArgs, err := json.Marshal(reactAppArgs)
	if err != nil {
		log.Println("json marshall app args error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	runReactApp := fmt.Sprintf(`
    runReact(%s);
    `, appArgs)

	val, err := rH.v8Ctx.RunScript(runReactApp, "main.js")
	if err != nil {
		log.Println("error run main.js", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !val.IsPromise() {
		log.Println("value is not a promise")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	promiseVal, err := rH.resolvePromise(val)
	if err != nil {
		log.Println("fail to resolve promise", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !promiseVal.IsObject() {
		log.Println("value is not an object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valM, err := promiseVal.MarshalJSON()
	if err != nil {
		log.Println("json marshall value markup error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var markup markupValue
	json.Unmarshal(valM, &markup)

	var templateData = struct {
		ReactApp string
	}{ReactApp: markup.ReactMarkup}
	err = rH.tmpl.ExecuteTemplate(w, "react.html", templateData)
	if err != nil {
		log.Println("error execute template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rH *reactHandler) resolvePromise(val *v8.Value) (*v8.Value, error) {
	for {
		p, _ := val.AsPromise()
		state := p.State()
		switch state {
		case v8.Fulfilled:
			return p.Result(), nil
		case v8.Rejected:
			return nil, errors.New(p.Result().DetailString())
		case v8.Pending:
			rH.v8Ctx.PerformMicrotaskCheckpoint() // run VM to make progress on the promise
			// go round the loop again...
		default:
			return nil, fmt.Errorf("illegal v8.Promise state %d", p) // unreachable
		}
	}
}

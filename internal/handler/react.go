package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	v8 "rogchap.com/v8go"
)

type reactHandler struct {
	v8Ctx         *v8.Context
	tmpl          *template.Template
	withHydration bool
}

type ReactHandlerOpts struct {
	V8Ctx         *v8.Context
	Tmpl          *template.Template
	WithHydration bool
}

func NewReactHandler(opts ReactHandlerOpts) *reactHandler {
	return &reactHandler{
		v8Ctx:         opts.V8Ctx,
		tmpl:          opts.Tmpl,
		withHydration: opts.WithHydration,
	}
}

type ReactHandlerRouter interface {
	RenderReact(w http.ResponseWriter, r *http.Request)
}

type route struct {
	Path string `json:"path"`
}

type routeMatch struct {
	Params       interface{} `json:"params"`
	Pathname     string      `json:"pathname"`
	PathnameBase string      `json:"pathnameBase"`
	Route        route       `json:"route"`
}

type markupValue struct {
	ReactMarkup string   `json:"markup"`
	EmotionCss  string   `json:"emotionCss"`
	EmotionIds  []string `json:"emotionIds"`
	EmotionKey  string   `json:"emotionKey"`
}

type templateData struct {
	ReactApp      string
	EmotionCss    string
	EmotionIds    string
	EmotionKey    string
	WithHydration bool
	AppState      string
}

func (rH *reactHandler) RenderReact(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.String()

	runMatchRoutes := fmt.Sprintf(`GO_APP.getMatchRoutes("%s")`, reqUrl)
	match, err := rH.v8Ctx.RunScript(runMatchRoutes, "match_routes.js")
	if err != nil {
		e := err.(*v8.JSError)
		fmt.Println(e.StackTrace)

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
	matchRoute := rM[0].Route.Path
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

	runReactApp := fmt.Sprintf(`GO_APP.render(%s)`, appArgs)

	val, err := rH.v8Ctx.RunScript(runReactApp, "render.js")
	if err != nil {
		log.Println("error run render.js", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valM, err := val.MarshalJSON()
	if err != nil {
		log.Println("json marshall value markup error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var markup markupValue
	json.Unmarshal(valM, &markup)

	eIds := strings.Join(markup.EmotionIds, " ")
	tData := templateData{
		ReactApp:      markup.ReactMarkup,
		EmotionCss:    markup.EmotionCss,
		EmotionIds:    eIds,
		EmotionKey:    markup.EmotionKey,
		WithHydration: rH.withHydration,
		AppState:      "{}",
	}
	err = rH.tmpl.ExecuteTemplate(w, "react.html", tData)
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

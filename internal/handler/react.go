package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/adityalstkp/gov8react/internal/constants"
	"github.com/adityalstkp/gov8react/internal/usecase"
	v8 "rogchap.com/v8go"
)

type reactHandler struct {
	v8Ctx         *v8.Context
	tmpl          *template.Template
	withHydration bool
	introUsecase  usecase.IntroUsecase
}

type ReactHandlerOpts struct {
	V8Ctx         *v8.Context
	Tmpl          *template.Template
	WithHydration bool
	IntroUsecase  usecase.IntroUsecase
}

func NewReactHandler(opts ReactHandlerOpts) *reactHandler {
	return &reactHandler{
		v8Ctx:         opts.V8Ctx,
		tmpl:          opts.Tmpl,
		withHydration: opts.WithHydration,
		introUsecase:  opts.IntroUsecase,
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

type emotionMarkup struct {
	Css string   `json:"css"`
	Key string   `json:"key"`
	Ids []string `json:"ids"`
}

type markupValue struct {
	Html    string        `json:"html"`
	Emotion emotionMarkup `json:"emotion"`
}

type emotionTemplate struct {
	Css string
	Key string
	Ids string
}

type templateData struct {
	ReactApp      string
	Emotion       emotionTemplate
	WithHydration bool
	AppState      string
}

func (rH *reactHandler) RenderReact(w http.ResponseWriter, r *http.Request) {
	reqUrl := r.URL.String()

	httpStatusCode := http.StatusOK
	runMatchRoutes := fmt.Sprintf(`GO_APP.getMatchRoutes("%s")`, reqUrl)
	match, err := rH.v8Ctx.RunScript(runMatchRoutes, "match_routes.js")

	if err != nil {
		httpStatusCode = http.StatusInternalServerError
		e := err.(*v8.JSError)
		log.Println(e.StackTrace)
		log.Println("error run match_routes.js", err)
		w.WriteHeader(httpStatusCode)
		return
	}

	if !match.IsArray() {
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	mJson, err := match.MarshalJSON()
	if err != nil {
		log.Println("cannot marshall route match")
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	var rM []routeMatch
	err = json.Unmarshal(mJson, &rM)
	if err != nil {
		log.Println("cannot umarshall route match")
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	initialData := map[string]interface{}{}
	matchRoute := rM[0].Route.Path

	if rM[0].Route.Path == "*" {
		httpStatusCode = http.StatusNotFound
	}

	sD, err := rH.getInitialData(matchRoute)
	if err != nil {
		httpStatusCode = http.StatusInternalServerError
		log.Println("cannot get initial data", err)
		w.WriteHeader(httpStatusCode)
		return
	}

	if sD != nil {
		initialData[matchRoute] = sD
	}

	reactAppArgs := map[string]interface{}{
		"url":         reqUrl,
		"initialData": initialData,
	}
	appArgs, err := json.Marshal(reactAppArgs)
	if err != nil {
		log.Println("json marshall app args error", err)
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	runReactApp := fmt.Sprintf(`GO_APP.render(%s)`, appArgs)
	val, err := rH.v8Ctx.RunScript(runReactApp, "render.js")
	if err != nil {
		log.Println("error run render.js", err)
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	valM, err := val.MarshalJSON()
	if err != nil {
		log.Println("json marshall value markup error", err)
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	var markup markupValue
	json.Unmarshal(valM, &markup)

	sS, err := json.Marshal(initialData)
	if err != nil {
		log.Println("cannot marshal static data", err)
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}

	eIds := strings.Join(markup.Emotion.Ids, " ")
	tData := templateData{
		ReactApp: markup.Html,
		Emotion: emotionTemplate{
			Css: markup.Emotion.Css,
			Ids: eIds,
			Key: markup.Emotion.Key,
		},
		WithHydration: rH.withHydration,
		AppState:      string(sS),
	}

	w.WriteHeader(httpStatusCode)
	err = rH.tmpl.ExecuteTemplate(w, "react.html", tData)
	if err != nil {
		log.Println("error execute template", err)
		httpStatusCode = http.StatusInternalServerError
		w.WriteHeader(httpStatusCode)
		return
	}
}

func (rH *reactHandler) getInitialData(matchRoute string) (interface{}, error) {
	switch matchRoute {
	case constants.REACT_ROUTE_INDEX_PATH:
		d, err := rH.introUsecase.Greet()
		if err != nil {
			return nil, err
		}

		return d, nil
	default:
		return nil, nil
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

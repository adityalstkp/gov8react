package main

import (
	"fmt"
	"log"

	"github.com/adityalstkp/gov8react/internal/constants"
	"github.com/adityalstkp/gov8react/internal/utilities"
	v8 "rogchap.com/v8go"
)

type v8Polyfills struct {
	name string
}

func newV8Ctx() (*v8.Context, error) {
	iso := v8.NewIsolate()

	global := createGlobalObject(iso)
	v8Ctx := v8.NewContext(iso, global)

	polyfills := []v8Polyfills{
		{
			name: "text_encoder",
		},
		{
			name: "buffer",
		},
	}
	err := injectPolyfills(v8Ctx, polyfills)
	if err != nil {
		return nil, err
	}

	appBundlePath := fmt.Sprintf("%s/main.js", constants.BASE_ARTIFACTS_DIR)
	appBundle, err := utilities.ReadFile(appBundlePath)
	if err != nil {
		return nil, err
	}

	_, err = v8Ctx.RunScript(string(appBundle), "bundle.js")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return v8Ctx, nil
}

func createGlobalObject(iso *v8.Isolate) *v8.ObjectTemplate {
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
	return global
}

func injectPolyfills(v8Ctx *v8.Context, polyfills []v8Polyfills) error {
	for _, p := range polyfills {
		pName := fmt.Sprintf("%s/polyfills.%s.js", constants.BASE_ARTIFACTS_DIR, p.name)
		pB, err := utilities.ReadFile(pName)
		if err != nil {
			log.Panicf("error read js polyfill %s", err.Error())
		}

		pO := fmt.Sprintf("polyfills_%s.js", p.name)
		_, err = v8Ctx.RunScript(string(pB), pO)
		if err != nil {
			log.Panicln("error in bundling polyfill", err)
		}
	}
	return nil
}

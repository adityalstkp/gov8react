# gov8react

React SSR with Go V8 binding

## Usage
```sh
Usage of .bin/gov8react:
  -http_addr string
        http listen address (default "0.0.0.0:3000")
  -with_hydration
        render with hydration
```

## Command
For those who want to try
```sh
make setup # install deps
make build-all-client # build all client artifacts
go run cmd/http/*.go -with_hydration true
```

## Current State
1. Hydration will load large bundle

## Stack
1. Go (v8 binding)
2. React
3. `rspack`
4. `@emotion/css`

## Limitations
These are the limitation that we found and these may comes because we lack of knowledges:
1. Node capabilities, v8 binding does not comes with node.js capabilities likes `fs`, `buffer`, etc. Working around with some polyfills.
2. `@emotion/react` not working, `cache` seems always empty, still don't know.
3. Code split will be hard things to do, because v8 context cannot refer/lookup webpack split files easily.

## Questions
1. How's the performance?

# gov8react

React SSR with Go V8 binding

## Stack
1. Go (v8 binding)
2. React
3. rspack
4. emotion

## Limitations
1. Node capabilities, v8 binding does not comes with node.js capabilities likes `fs`, `buffer`, etc. Working around with some polyfills.
2. `@emotion/react` not working, `cache` seems always empty, still don't know.
3. It's only for *static* SSR, once sent to the browser it will not be interactive.
4. Data sent from `go` to `react`, `react` will not query the data, `go` is the one that provides.

## Questions
1. How's the performance?

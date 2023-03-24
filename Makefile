
clean-artifacts:
	rm -rf .artifacts

build-react-client: 
	(cd client && NODE_ENV=$(node_env) pnpm rspack -c rspack.client.config.js)

build-react-server: 
	(cd client && NODE_ENV=$(node_env) pnpm rspack -c rspack.server.config.js)

build-polyfills: 
	(cd client && NODE_ENV=$(node_env) pnpm rspack -c rspack.polyfills.config.js)

setup: clean-artifacts
	go mod tidy
	(cd client && pnpm i)

run-dev: clean-artifacts build-react-client build-react-server build-polyfills
	go run cmd/http/*.go

run: clean-artifacts build-react-client build-react-server build-polyfills
	go build -o .bin/gov8react cmd/http/*.go && .bin/gov8react

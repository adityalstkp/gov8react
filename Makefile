NODE_ENV:="development"

.PHONY: clean-artifacts build-react-client build-react-server build-polyfills build-all-client setup 

clean-artifacts:
	rm -rf .artifacts

build-react-client: 
	(cd client && NODE_ENV=$(NODE_ENV) pnpm rspack -c rspack.client.config.js)

build-react-server: 
	(cd client && NODE_ENV=$(NODE_ENV) pnpm rspack -c rspack.server.config.js)

build-polyfills: 
	(cd client && NODE_ENV=$(NODE_ENV) pnpm rspack -c rspack.polyfills.config.js)

setup: clean-artifacts
	go mod tidy
	(cd client && pnpm i)

build-all-client: clean-artifacts build-react-client build-react-server build-polyfills

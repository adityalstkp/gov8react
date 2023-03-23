clean-artifacts:
	rm -rf .artifacts

setup: clean-artifacts
	go mod tidy
	(cd client && pnpm i)

run-dev: clean-artifacts
	(cd client && pnpm build:dev)
	go run cmd/http/main.go

run: clean-artifacts
	(cd client && pnpm build:prod)
	go build -o .bin/gov8react cmd/http/main.go && .bin/gov8react

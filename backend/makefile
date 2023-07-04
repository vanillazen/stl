app = stl

.PHONY: direnv
direnv:
	direnv allow .

.PHONY: build
build:
	go build ./...

.PHONY: run
run:
	go run main.go

# Testing
.PHONY: test
test:
	make -f makefile.test test-selected

# OpenAPI
.PHONY: openapi/gen
openapi/gen:
	swag init -g main.go -o api/openapi
	rm -rf api/openapi/docs.go
	mv api/openapi/swagger.json api/openapi/openapi.json
	mv api/openapi/swagger.yaml api/openapi/openapi.yaml
	redocly build-docs api/openapi/openapi.yaml -o api/openapi/openapi.html

# Tools
.PHONY: redocly/install
redocly/install:
	npm i -g @redocly/cli@latest


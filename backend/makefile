app = stl

.PHONY: direnv
direnv:
	direnv allow .

.PHONY: build
build:
	make openapi/gen
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

# API
.PHONY: api/get-list
api/get-list:
	./scripts/curl/get-list.sh -h localhost -p 8080 -l cdc7a443-3c6a-431b-b45a-b14735953a19 --pretty

# Tools (sudo reuired)
.PHONY: install/redocly
install/redocly:
	npm i -g @redocly/cli@latest

.PHONY: install/jq/arch
install/jq/arch:
	pacman -Sy jq

.PHONY: install/jq/debian
install/jq/debian:
	apt update;
	apt install -y jq

.PHONY: install/jq/centos
install/jq/centos:
	yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
	yum install jq -y




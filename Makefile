all:
	go build -o ./build/a.out .

clear:
	rm -rf ./build/*

install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

generate-oapi:
	oapi-codegen -generate chi-server,types,spec -package api ./api/openapi.yml > ./api/api.gen.go
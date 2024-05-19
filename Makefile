.PHONY: run test

setup:
	go get golang.org/x/text/language
	go get golang.org/x/text/message
	go get golang.org/x/sys/unix

build:
	go build -x .

test:
	go test -v \
		lexer.go parser.go jsonValue.go \
		*_test.go

DEFAULT_PAIRS := 3
generate:
	go run haversine-generate-json.go -pairs=$(if $(pairs),$(pairs),$(DEFAULT_PAIRS))

run:
	go run haversine-compute-real-json-parser.go \
		lexer.go parser.go jsonValue.go \
		haversine.go

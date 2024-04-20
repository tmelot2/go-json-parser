.PHONY: run test

test:
	go test -v \
		lexer.go parser.go \
		*_test.go

DEFAULT_PAIRS := 3
generate:
	go run haversine-generate-json.go -pairs=$(if $(pairs),$(pairs),$(DEFAULT_PAIRS))

run:
	go run haversine-compute-real-json-parser.go \
		lexer.go parser.go \
		haversine.go

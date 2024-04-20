.PHONY: run test

test:
	go test -v \
		lexer.go parser.go \
		*_test.go

run:
	go run haversine-compute-real-json-parser.go \
		lexer.go parser.go \
		haversine.go

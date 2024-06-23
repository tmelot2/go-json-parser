# Go JSON Parser with Haversine Calculator

A Haversine calculator that uses a simple JSON parser written in Go. Still a work-in-progress.

Created with Go 1.22.1, no external dependencies.

## Origins

This project was created as part of the coursework for the ["Performance-Aware Programming" course](https:#www.computerenhance.com/p/table-of-contents) I am taking. Some of the coursework artifacts are still in this project, like the Haversine stuff. It seemed useful to make the parser stuff into a separate module to gain experience with that while learning Go.

I am also learning Go, & this is the 1st project I am using it with, so there's bound to be mistakes!

## Usage

Run JSON parser unit tests:

```sh
cd internal/jsonParser
go test -v .
```

Generate `pairs.json` file:
```sh
cd cmd/generateJson

# Generate with default num pairs
go run .

# Generate with custom num pairs
go run . -pairs=100000
```

Run Haversine compute with my JSON parser:
```sh
cd cmd/myJsonParser

# Run app
go run .

# Run app with profiling
go run -tags=profile .
```


## Progress

- Parser works!
- Supported types: Object, array, string, int, float
- Parsed data is of type JsonValue, which you can use to pick out typed data.
- There are unit tests for the lexer & parser, which will continue to be expanded.

## Todo

- Currently, to use, you call `parser.ParseJson(fileData)`. This requires the entire file is loaded into memory. Bad for big files!
	- [ ] Refactor to use a streaming reader, maybe make a JsonParser class that recieves the exernal methods

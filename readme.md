# Go JSON Parser with Haversine Calculator

A Haversine calculator that uses a from-scratch, eventually highly-performance JSON parser written in Go. It's a project for the ["Performance-Aware Programming" course](https://www.computerenhance.com) I am taking.

Project pieces:

- A from-scratch JSON parser.
- A from-scratch block profiler to analyze CPU cycles, memory bandwidth, & more.
- A from-scratch repetition tester for further performance analysis.
- Future: An extrenely optimized & performant JSON parser. In the course we are going down to the assembly-language level of analysis & optimization, leveraging CPU architecture features like caching, vector instructions, etc.

Still a work-in-progress.

Created with Go 1.22.1, no external dependencies.

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

- Parser
	- It works!
	- Supported types: Object, array, string, int, float. Missing: bool.
	- Parsed data is of type JsonValue, which you can use to pick out typed data.
	- There are unit tests for the lexer & parser, which will continue to be expanded.
	- Currently ~9x slower than Go's builtin parser. Lots of room for improvement!
- Block profiler
	- It also works! And it's so cool to use it!
	- For each block, measures CPU cycles, hit count, & optionally memory bandwidth.
- Repetition tester
	- WIP

## Todo

- Currently, to use, you call `parser.ParseJson(fileData)`. This requires the entire file is loaded into memory. Bad for big files!
	- [ ] Refactor to use a streaming reader, maybe make a JsonParser class that recieves the exernal methods

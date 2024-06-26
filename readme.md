# Go JSON Parser with Haversine Calculator

A Haversine calculator that uses a from-scratch, eventually highly-performance JSON parser written in Go. It's a project for the ["Performance-Aware Programming" course](https://www.computerenhance.com) I am taking.

Project pieces:

- A from-scratch JSON parser.
- A from-scratch block profiler to analyze CPU cycles, memory bandwidth, & more.
- A from-scratch repetition tester for further performance analysis.
- Future: An extremely optimized & performant JSON parser. In the course we are going down to the assembly-language level of analysis & optimization, leveraging CPU architecture features like caching, vector instructions, etc.

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

Run repetition tester:
```sh
cd cmd/repetitionTest

# Run app
go run .
```

## Example Profiler Output

A run with 1 million pairs on an AMD Ryzen 9 5900X @ 3.7 Ghz:

```
Count:         1,000,000
Haversine sum: 8,623,222,205.7044734954833984
Haversine avg: 8,623.2222057044727990

[CPU profiling stats]
Total time: 7,248.6769ms (CPU freq  3,700,051,600 Hz)
          Block                Cycles |  Hit Cnt | Percent
        Startup:               16,835 |        1 | 0.00%
           Read:          114,389,146 |        1 | 0.43%, 105.301mb at 3.33gb/s
      ReadToStr:          108,329,118 |        1 | 0.40%
         Parser:               16,909 |        1 | 0.00%, 95.23% w/children
     Parser.Lex:       22,827,895,456 |        1 | 85.11%
   Parser.Parse:        2,712,289,181 |        1 | 10.11%
   SumHaversine:        1,056,989,343 |        2 | 3.94%, 30.518mb at 0.10gb/s
     MiscOutput:              178,858 |        1 | 0.00%
============================================================
       Profiler:               59,200 |       18 | 0.00%
          Total:       26,820,478,472 |        0 | 100.00%
```

## Example Repetition Tester Output

A run that repeatedly loads 10 million pairs:

```
OS.ReadFile :
Min: 715707021 (193.432730ms) 5.282137gb/s
Max: 2174110084 (587.592460ms) 1.738855gb/s
Avg: 759530706 (205.276871ms) 4.977367gb/s

ioutil.ReadFile :
Min: 721926739 (195.113721ms) 5.236629gb/s
Max: 828885191 (224.021172ms) 4.560900gb/s
Avg: 738089858 (199.482096ms) 5.121955gb/s

bufio.Reader :
Min: 1440673033 (389.367870ms) 2.624095gb/s
Max: 1864848749 (504.009007ms) 2.027222gb/s
Avg: 1495594008 (404.211254ms) 2.527733gb/s

bytes.Buffer :
Min: 1456971513 (393.772828ms) 2.594740gb/s
Max: 10149646361 (2743.124979ms) 0.372472gb/s
Avg: 1823100121 (492.725687ms) 2.073645gb/s
```


## Progress

- Parser
	- It works!
	- Supported types: Object, array, string, int, float. Missing: bool.
	- Parsed data is of type JsonValue, which you can use to pick out typed data.
	- There are unit tests for the lexer & parser, which will continue to be expanded.
	- Currently ~9x slower than Go's builtin parser. GOOD, lots of room for improvement!
- Block profiler
	- It also works! And it's so cool to use it!
	- For each block, measures CPU cycles, hit count, & optionally memory bandwidth.
	- Supports nested profiled blocks.
	- TODO: Support recursive profiled blocks. You can do this now, but the numbers get crazy & meaningless.
- Repetition tester
	- It also also works!
	- Repeatedly measures calls to different Go stl file read functions in order to find the "stars align" best-case speed.

## Todo

- Currently, to use, you call `parser.ParseJson(fileData)`. This requires the entire file is loaded into memory. Bad for big files!
	- TODO: Refactor to use a streaming reader, maybe make a JsonParser class that recieves the exernal methods

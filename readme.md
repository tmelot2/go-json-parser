# Go JSON Parser + "Performance-Aware Programming" Coursework

This project is coursework for the ["Performance-Aware Programming" course](https://www.computerenhance.com) I'm taking. The core app is a Haversine calculator that uses a from-scratch, eventually highly-performance JSON parser written in Go.

Project pieces:

- A from-scratch JSON parser.
- A from-scratch block profiler to analyze CPU cycles, memory bandwidth, & more.
- A from-scratch repetition tester for further performance analysis.
- Future: An extremely optimized & performant JSON parser. In the course we are going down to the assembly-language level of analysis & optimization, leveraging CPU architecture features like caching, vector instructions, etc.

Still a work-in-progress.

Created with Go 1.22.1 with minimal dependencies.

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

Run repetition tester (with file loading function comparisons):
```sh
cd cmd/repetitionTest

# Run app
go run loadFile.go
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
Using alloc type: None

--- OS.ReadFile ---
Min: 888444 (0.240113ms) 4.292762gb/s PF: 272 (3.9736k/fault)
Max: 3351756 (0.905853ms) 1.137875gb/s PF: 278 (3.8878k/fault)
Avg: 1203035 (0.325135ms) 3.170215gb/s PF: 272 (3.9736k/fault)

--- WriteToAllBytes ---
Min: 1661892 (0.449146ms) 2.294902gb/s
Max: 2981037 (0.805661ms) 1.279380gb/s
Avg: 1751498 (0.473363ms) 2.177495gb/s

--- ioutil.ReadFile ---
Min: 887371 (0.239823ms) 4.297953gb/s PF: 271 (3.9883k/fault)
Max: 2685497 (0.725788ms) 1.420176gb/s PF: 272 (3.9736k/fault)
Avg: 1183749 (0.319922ms) 3.221866gb/s PF: 272 (3.9736k/fault)
```

## Running Custom Assembly Routines

We have to examine & tweak custom assembly language routines. Here's an example of how to do that using `cmd/repetitionTest/nopLoop.go`, which compares 4 routines:

1. Write the assembly (see `nopLoop.asm`).
2. Assemble it with nasm:
	```sh
	nasm -f win64 -o nopLoop.obj nopLoop.asm
	```
3. Create a linkable library from the assembled obj:
	```sh
	lib nopLoop.obj
	```
4. Setup cgo in Go source:
	```go
	/*
	#cgo CFLAGS: -I.
	#cgo LDFLAGS: -L. -lnopLoop

	typedef char u8;
	typedef long long unsigned u64;
	void MOVAllBytesASM(u64 count, u8 *data);
	void NOPAllBytesASM(u64 count);
	void CMPAllBytesASM(u64 count);
	void DECAllBytesASM(u64 count);
	*/
	import "C"
	```
5. Call your functions like `C.MOVAllBytesASM(...)`

## Progress

- Parser
	- Works!
	- Supported types: Object, array, string, int, float, bool.
	- Parsed data is type `JsonValue`, which you can use to get typed data.
	- There are unit tests for the lexer & parser, which will continue to be expanded.
	- Currently ~9x slower than Go's builtin parser. GOOD, lots of room for improvement!
	- See `./internal/jsonParser/jsonValue.go` for usage.
- Block profiler
	- Also works! And it's so cool to use it!
	- For each block, measures CPU cycles, hit count, & optionally memory bandwidth.
	- Supports nested profiled blocks.
	- TODO: Support recursive profiled blocks. You can do this now, but the numbers get crazy & meaningless.
- Repetition tester
	- Also also works!
	- Repeatedly measures calls to different functions to find the "stars align" best-case speed. Outputs best, worse, & average.

## Todo

- Currently, to use, you call `parser.ParseJson(fileData)`. This requires the entire file is loaded into memory. Bad for big files!
	- TODO: Refactor to use a streaming reader, maybe make a JsonParser class that recieves the exernal methods

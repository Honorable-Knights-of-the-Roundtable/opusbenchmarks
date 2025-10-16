# OPUS Benchmarks

Benchmarking some OPUS libraries in Golang.

### Running this project

Run the benchmarks using `go test -bench=.`. Be warned the benchmarks encompass a wide number of sample rates, channel counts, and frame durations, meaning they will take a long time to complete!

Change the constant `ENCDEC_TYPE` at the top of `benchmark_test.go` to alter what encoderdecoder is used during benchmarking. Output the result of each encoderdecoder to a file with `go test -bench=. > file.txt`. Then, compare the trials of different encoderdecoders using [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat): `benchstat file1.txt file2.txt`

# Implementations

### [`hraban/opus`](github.com/hraban/opus)

The "original" Go implementation of OPUS, relying on several dev-dependencies that make building a project awkward. Namely, this project relies on CGo, which makes cross-compilation (and Windows compilation) difficult.

### [`jj11hh/opus`](github.com/jj11hh/opus)

A fork of `hraban/opus` that exchanges the CGo dependency for a [`wazero`](https://github.com/wazero/wazero) dependency. `wazero` is a WebAssembly engine built entirely in Go, which allows us to use a WASM build of OPUS that is embedded in the project itself. This make the `jj11hh/opus` library much more portable than `hraban/opus`.
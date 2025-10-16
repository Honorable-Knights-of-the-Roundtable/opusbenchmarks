# OPUS Benchmarks

Benchmarking some OPUS libraries in Golang.

### Running this project

Run the benchmarks using `go test -bench=.`. Be warned the benchmarks encompass a wide number of sample rates, channel counts, and frame durations, meaning they will take a long time to complete!

Change the constant `ENCDEC_TYPE` at the top of `benchmark_test.go` to alter what encoderdecoder is used during benchmarking. Output the result of each encoderdecoder to a file with `go test -bench=. > file.txt`. Then, compare the trials of different encoderdecoders using [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat): `benchstat file1.txt file2.txt`

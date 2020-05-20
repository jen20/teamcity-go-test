# TeamCity Go Test Runner

`teamcity-go-test` is a replacement test runner for TeamCity for Go. It is
heavily inspired by [Pavel Gulbin's work][1] <other repo>, but modified to fit
our use case, where we have lots of tests which take a long time but are safe
to run in parallel. The workflow is:

1. Compile a test binary using `go test -c`
2. Pipe a list of test names, one per line, into `teamcity-go-test` , with the
   `-test` parameter pointing to the executable. To run tests in parallel, set
   `-parallelism` to a number greater than 1.


[1]: https://github.com/2tvenom/go-test-teamcity

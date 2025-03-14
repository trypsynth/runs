# runs
This is a simple HTTP server written in Go that tracks how many times a specific application has been run. The data is persisted in a JSON file.

## Getting Started

Clone this repository and run `go build` to generate a binary. after doing that, run the binary, and make a POST request to http://127.0.0.1:7867/runs?name=<app_name> to see the app working.

# Go Hiring Project

Provide a `CV.pdf` to the web service and scan & group it for relevant topics.

## Usage

If you want to run it, use the commands:

```
go build main.go
./main
```

This runs a service on port 8080.

To test this service, go to `localhost:8080`.

## Example Output

Using my `CV.pdf`, example part of the output:

```
...
Topics with non-zero values in Column 10
- Topic: aachen
- Topic: thesis
- Topic: scientic
- Topic: team
- Topic: english
- Topic: metadata
- Topic: profession
...
```

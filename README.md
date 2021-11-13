# Dr. Cat

Quickly render markdown text to the console. 
Written in Go with [glamour](https://github.com/charmbracelet/glamour).

## usage

```
fmt-md-text [-f md-filename] [-l]
```

```
fmt-md-text -f README.md
```

## build

```bash
git clone
go build .
```

## test

You need to build the binary to run the tests.

```bash
go build .
go test -v
```

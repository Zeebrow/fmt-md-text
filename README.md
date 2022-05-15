# Dr. Cat

Quickly render markdown text to the console. 
Written in Go with [glamour](https://github.com/charmbracelet/glamour).

## example usage

Don't run as root.

```bash
fmt-md-text [-f md-filename] [-l]
```

```bash
fmt-md-text < README.md
fmt-md-text -f README.md
fmt-md-text -f README.md | less -R
```

`NICE TO HAVE:` detect content length exceeding terminal height and drop to `$PAGER`

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

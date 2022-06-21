# Dr. Cat

Quickly render markdown text to the console. 
Written in Go with [glamour](https://github.com/charmbracelet/glamour).

## download release

Debian (apt)
```
wget https://github-artifacts-zeebrow.s3.amazonaws.com/drcat/releases/debian/v1.0.0/amd64/drcat-1.0.0.deb
wget https://github-artifacts-zeebrow.s3.amazonaws.com/drcat/releases/debian/v1.0.0/amd64/drcat-1.0.0.deb.md5
md5sum -c drcat-1.0.0.deb.md5
```

Binary linux amd64
```
wget https://github-artifacts-zeebrow.s3.amazonaws.com/drcat/releases/binaries/v1.0.0/linux/amd64/drcat-1.0.0
wget https://github-artifacts-zeebrow.s3.amazonaws.com/drcat/releases/binaries/v1.0.0/linux/amd64/drcat-1.0.0.md5
md5sum -c drcat-1.0.0.md5
mv drcat-1.0.0 drcat
chmod +x drcat
./drcat -version
```

Binary windows amd64
```

```

## example usage

Be careful running as root!

```bash
fmt-md-text [-f md-filename] [-l]
```

```bash
fmt-md-text < README.md
fmt-md-text -f README.md
fmt-md-text -f README.md | less -R
```

`NICE TO HAVE:` detect content length exceeding terminal height and drop to `$PAGER`

`NICE TO HAVE:` `-h` (or some option) gives markdown formatting guide

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

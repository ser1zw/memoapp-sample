# Memo app

Sample application in Go

## Requirements
- [Go](https://github.com/golang/go)
- [dep](https://github.com/golang/dep)
- [SQLite3](https://www.sqlite.org/index.html)

## Installation & Run

```console
$ git clone https://github.com/ser1zw/memoapp-sample.git
$ cd memoapp-sample
$ dep ensure
$ sqlite3 db/memoapp.db < db/init.sql
$ go run app.go
```

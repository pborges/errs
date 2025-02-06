package main

import (
	"errors"
	"github.com/pborges/errs"
	"io"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := a()
	log.Println(err)
	log.Println(errs.Detailed(err))
	log.Println(errs.Detailed(f()))
	log.Println(errs.Detailed(noerr()), noerr() == nil)
	log.Println(errs.Detailed(errs.Push(io.EOF)))
	log.Println(errs.Detailed(errs.Wrap(io.EOF, errors.New("i am no good"))))

	log.Println(errors.Is(err, io.EOF))
	log.Println(errors.Is(err, io.ErrNoProgress))
}

func a() error {
	return errs.Push(b())
}

func b() error {
	return errs.Push(c())
}

func c() error {
	return errs.Wrap(d(), io.ErrNoProgress)
}

func d() error {
	return errs.Push(e())
}

func e() error {
	return errs.Push(io.EOF)
}

func f() error {
	return errs.Push(io.EOF)
}

func noerr() error {
	return errs.Push(nil)
}

/**
OUTPUT:
2025/02/05 19:09:24 main.go:16: multiple Read calls return no data or error
2025/02/05 19:09:24 main.go:17: multiple Read calls return no data or error
│┬ main.go:28 (main.a) ^
│└┬ main.go:32 (main.b) ^
│ └┬ main.go:36 (main.c) multiple Read calls return no data or error
│  └┬ main.go:40 (main.d) ^
│   └─ main.go:44 (main.e) EOF
2025/02/05 19:09:24 main.go:18: EOF
2025/02/05 19:09:24 main.go:19:  true
2025/02/05 19:09:24 main.go:20: EOF
2025/02/05 19:09:24 main.go:21: i am no good
│┬ main.go:21 (main.main) i am no good
│└─ main.go:21 (main.main) EOF
2025/02/05 19:09:24 main.go:23: true
2025/02/05 19:09:24 main.go:24: true
*/

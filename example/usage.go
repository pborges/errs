package main

import (
	"errors"
	"fmt"
	"github.com/pborges/errs"
	"io"
	"log"
	"os"
)

type JSError struct {
	Message  string
	Location string
}

func (e JSError) Error() string {
	return e.Message
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := a()
	fmt.Println("* Normal print")
	log.Println(err)

	fmt.Println("\n* Detailed print")
	log.Println(errs.Detailed(err))

	fmt.Println("\n* Nil error")
	log.Println(errs.Detailed(noerr()))

	fmt.Println("\n* Wrapped Error")
	log.Println(errs.Detailed(errs.Wrap(io.EOF, errors.New("i am no good"))))

	fmt.Println("\n* Is io.EOF")
	log.Println(errors.Is(err, io.EOF))

	fmt.Println("\n* Is io.ErrNoProgress")
	log.Println(errors.Is(err, io.ErrNoProgress))

	fmt.Println("\n* Is io.ErrClosedPipe")
	log.Println(errors.Is(err, io.ErrClosedPipe))

	fmt.Println("\n* Multiline")
	log.Println(errs.Detailed(multiline()))

	fmt.Println("\n* JSError with no transformer")
	log.Println(errs.Detailed(errs.Push(errs.Push(JSError{Message: "foo", Location: "bar:12"}))))

	// Add a transformer
	errs.Transform(func(err error) (bool, string) {
		var jerr JSError
		if errors.As(err, &jerr) {
			return true, jerr.Message + " => " + jerr.Location
		}
		return false, ""
	})

	fmt.Println("\n* JSError with with transformer")
	log.Println(errs.Detailed(errs.Push(errs.Push(JSError{Message: "foo", Location: "bar:12"}))))
}

func multiline() error {
	err := errs.Push(errors.New("multiline error\nline1\nline2"))
	err = errs.Push(err)
	err = errs.Push(err)
	err = errs.Push(err)
	return err
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
2025/02/06 18:41:19 main.go:16: multiple Read calls return no data or error
2025/02/06 18:41:19 main.go:17: multiple Read calls return no data or error
│┬ main.go:29 (main.a) ^
│└┬ main.go:33 (main.b) ^
│ └┬ main.go:37 (main.c) multiple Read calls return no data or error
│  └┬ main.go:41 (main.d) ^
│   └─ main.go:45 (main.e) EOF
2025/02/06 18:41:19 main.go:18: EOF
2025/02/06 18:41:19 main.go:19:  true
2025/02/06 18:41:19 main.go:20: EOF
2025/02/06 18:41:19 main.go:21: i am no good
│┬ main.go:21 (main.main) i am no good
│└─ main.go:21 (main.main) EOF
2025/02/06 18:41:19 main.go:23: Is io.EOF true
2025/02/06 18:41:19 main.go:24: Is io.ErrNoProgress true
2025/02/06 18:41:19 main.go:25: Is io.ErrClosedPipe false
*/

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
* Normal print
2025/03/07 18:40:06 usage.go:27: multiple Read calls return no data or error

* Detailed print
2025/03/07 18:40:06 usage.go:30: multiple Read calls return no data or error
│┬ usage.go:75 (main.a) ^
│└┬ usage.go:79 (main.b) ^
│ └┬ usage.go:83 (main.c) multiple Read calls return no data or error
│  └┬ usage.go:87 (main.d) ^
│   └─ usage.go:91 (main.e) EOF

* Nil error
2025/03/07 18:40:06 usage.go:33:

* Wrapped Error
2025/03/07 18:40:06 usage.go:36: i am no good
│┬ usage.go:36 (main.main) i am no good
│└─ usage.go:36 (main.main) EOF

* Is io.EOF
2025/03/07 18:40:06 usage.go:39: true

* Is io.ErrNoProgress
2025/03/07 18:40:06 usage.go:42: true

* Is io.ErrClosedPipe
2025/03/07 18:40:06 usage.go:45: false

* Multiline
2025/03/07 18:40:06 usage.go:48: multiline error
│ line1
│ line2
│┬ usage.go:70 (main.multiline) ^
│└┬ usage.go:69 (main.multiline) ^
│ └┬ usage.go:68 (main.multiline) ^
│  └─ usage.go:67 (main.multiline) multiline error
│     line1
│     line2

* JSError with no transformer
2025/03/07 18:40:06 usage.go:51: foo
│┬ usage.go:51 (main.main) ^
│└─ usage.go:51 (main.main) foo

* JSError with with transformer
2025/03/07 18:40:06 usage.go:63: foo => bar:12
│┬ usage.go:63 (main.main) ^
│└─ usage.go:63 (main.main) foo => bar:12
*/

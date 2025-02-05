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
	log.Println(errs.Detailed(errs.Wrap(io.EOF)))
	log.Println(errs.Detailed(errs.Wrapf(io.EOF, "i am no good")))

	log.Println(errors.Is(err, io.EOF))
}

func a() error {
	return errs.Wrap(b())
}

func b() error {
	return errs.Wrapf(c(), "unable to call c")
}

func c() error {
	return errs.Wrap(d())
}

func d() error {
	return errs.Wrapf(io.EOF, "unable to read")
}

func f() error {
	return errs.Wrap(io.EOF)
}

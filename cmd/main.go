package main

import (
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
	log.Println(errs.Dump(err))
	log.Println(errs.Dump(errs.Wrap(io.EOF)))
	log.Println(errs.Dump(errs.Wrapf(io.EOF, "i am no good")))
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

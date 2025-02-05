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
	return errs.Errorf("unable to read %w", io.EOF)
}

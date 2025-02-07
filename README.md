# errs

An opinionated, non-intrusive error stacker.

Inbune errors with Caller info and bubble it up the stack.

## Motivation
I am lazy. I often just catch errors and throw them up the call chain. 

Then I log the errors at the very top, loosing a ton of contextual information!

Or when I am feeling less lazy... I actually wrap errors using `fmt.Errorf`,

However after a few nested wraps the error messages start to become... long.

### Code
```go
func a0() error {
	if err := a1(); err != nil {
		return fmt.Errorf("unable to a1: %w", err)
	}
	return nil
}

func a1() error {
	if err := a2(); err != nil {
		return fmt.Errorf("unable to a2: %w", err)
	}
	return nil
}

func a2() error {
	if err := a3(); err != nil {
		return fmt.Errorf("unable to a3: %w", err)
	}
	return nil
}

func a3() error {
	return io.EOF
}

fmt.Println(a0())
```
### Output
```
unable to a1: unable to a2: unable to a3: EOF
```

This simple library builds something akin to a Java stack trace, but just stacks the errors.

And yes! `errors.Is` and `errors.As` still work!

## Usage

Simple to use, where you would normally bubble up an error,

```go
func MyFunction() error {
    err := myErrorProducingFunction()
    return err
}
```
Push it instead

# Push an error
```go
func MyFunction() error {
    err := myErrorProducingFunction()
    return errs.Push(err)
}
```
Where you would normally `fmt.Errorf` to wrap an error, instead...

# Wrap an error
```go
func MyFunction() error {
    err := myErrorProducingFunction()
    return errs.Wrap(err, errors.New("yea that didnt go as planned...."))
}
```

# Dump the error stack
```go
fmt.Println(errs.Detailed(err))
```
The `errs.Detailed` method will operate on any error, if the error is an `errs.Stack` it will be printed.

For example:
```
2025/02/05 19:09:24 main.go:17: multiple Read calls return no data or error
│┬ main.go:28 (main.a) ^
│└┬ main.go:32 (main.b) ^
│ └┬ main.go:36 (main.c) multiple Read calls return no data or error
│  └┬ main.go:40 (main.d) ^
│   └─ main.go:44 (main.e) EOF
```

package errs

type Transformer func(error) (bool, string)

var transformers []Transformer

func transform(err error) string {
	for _, t := range transformers {
		if ok, msg := t(err); ok {
			return msg
		}
	}
	return err.Error()
}

func Transform(transformer func(error) (bool, string)) {
	transformers = append(transformers, transformer)
}

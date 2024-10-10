package utils

func Must[T any](fn func() (T, error)) T {
	result, err := fn()
	if err != nil {
		panic(err) // Panic if there's an error
	}
	return result // Return the result if no error
}

func Must1[T any, A1 any](fn func(A1) (T, error), a1 A1) T {
	result, err := fn(a1)
	if err != nil {
		panic(err) // Panic if there's an error
	}
	return result // Return the result if no error
}

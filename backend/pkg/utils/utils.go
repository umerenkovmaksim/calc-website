package utils

import (
	"errors"
	"io"
	"log"
	"strconv"
)

var ErrArrayEmpty = errors.New("array is empty")

func CloseResponseBody(body io.Closer) {
	if err := body.Close(); err != nil {
		log.Printf("close body error: %v", err)
	}
}

func Pop[T any](array *[]T) (T, error) {
	if len(*array) == 0 {
		var zeroVar T
		return zeroVar, ErrArrayEmpty
	}

	elem := (*array)[len(*array)-1]
	*array = (*array)[:len(*array)-1]

	return elem, nil
}

func IsNumber(line string) bool {
	_, err := strconv.ParseFloat(line, 10)
	if err != nil {
		return false
	}
	return true
}

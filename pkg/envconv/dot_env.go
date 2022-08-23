package envconv

import (
	"io"
	"os"

	"github.com/joho/godotenv"
)

func MapFromReader(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

func MapFromFile(name string) (map[string]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return MapFromReader(f)
}

func ArrFromReader(r io.Reader) ([]string, error) {
	m, err := MapFromReader(r)
	if err != nil {
		return nil, err
	}

	return MapToArr(m), nil
}

func ArrFromFile(name string) ([]string, error) {
	m, err := MapFromFile(name)
	if err != nil {
		return nil, err
	}

	return MapToArr(m), nil
}

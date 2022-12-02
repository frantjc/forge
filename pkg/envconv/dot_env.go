package envconv

import (
	"io"
	"os"

	"github.com/joho/godotenv"
)

// MapFromReader takes a Reader with the content of a .env file e.g.
//
//		FOO=bar
//	 KEY=val
//
// and returns a corresponding map
//
//	{
//		"FOO": "bar",
//		"KEY": "val"
//	}
func MapFromReader(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

// MapFromFile takes a path to a file with the content of a .env file e.g.
//
//		FOO=bar
//	 KEY=val
//
// and returns a corresponding map
//
//	{
//		"FOO": "bar",
//		"KEY": "val"
//	}
func MapFromFile(name string) (map[string]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return MapFromReader(f)
}

// ArrFromReader takes a Reader with the content of a .env file e.g.
//
//		FOO=bar
//	 KEY=val
//
// and returns a corresponding environment array
//
// ["FOO=bar", "KEY=val"].
func ArrFromReader(r io.Reader) ([]string, error) {
	m, err := MapFromReader(r)
	if err != nil {
		return nil, err
	}

	return MapToArr(m), nil
}

// ArrFromFile takes a path to a file with the content of a .env file e.g.
//
//		FOO=bar
//	 KEY=val
//
// and returns a corresponding environment array
//
// ["FOO=bar", "KEY=val"].
func ArrFromFile(name string) ([]string, error) {
	m, err := MapFromFile(name)
	if err != nil {
		return nil, err
	}

	return MapToArr(m), nil
}

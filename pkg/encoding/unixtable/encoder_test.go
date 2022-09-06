package unixtable_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/frantjc/forge/pkg/encoding/unixtable"
)

func TestEncoderStruct(t *testing.T) {
	var (
		ut = &Unixtable{
			One: "general",
			Two: "kenobi",
		}
		expected = []byte(`One       Two
general   kenobi
`)
		buf = new(bytes.Buffer)
	)

	if err := unixtable.NewEncoder(buf).Encode(ut); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var (
		actualStr   = buf.String()
		expectedStr = string(expected)
	)

	if !strings.EqualFold(actualStr, expectedStr) {
		t.Error("actual\n", actualStr, "does not equal expected\n", expectedStr)
		t.FailNow()
	}
}

func TestEncoderTaggedStruct(t *testing.T) {
	var (
		ut = &UnixtableTagged{
			One: "general",
			Two: "kenobi",
		}
		expected = []byte(`one       two
general   kenobi
`)
		buf = new(bytes.Buffer)
	)

	if err := unixtable.NewEncoder(buf).Encode(ut); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var (
		actualStr   = buf.String()
		expectedStr = string(expected)
	)

	if !strings.EqualFold(actualStr, expectedStr) {
		t.Error("actual\n", actualStr, "does not equal expected\n", expectedStr)
		t.FailNow()
	}
}

func TestEncoderSlice(t *testing.T) {
	var (
		ut = []*Unixtable{
			{
				One: "hello",
				Two: "there",
			},
			{
				One: "general",
				Two: "kenobi",
			},
		}
		expected = []byte(`One       Two
hello     there
general   kenobi
`)
		buf = new(bytes.Buffer)
	)

	if err := unixtable.NewEncoder(buf).Encode(ut); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var (
		actualStr   = buf.String()
		expectedStr = string(expected)
	)

	if !strings.EqualFold(actualStr, expectedStr) {
		t.Error("actual\n", actualStr, "does not equal expected\n", expectedStr)
		t.FailNow()
	}
}

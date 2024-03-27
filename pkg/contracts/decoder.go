package contracts

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/ajg/form"
)

type Decoder interface {
	Decode(interface{}) error
}

func ContentDecoder(contentType string) (func(r io.Reader) Decoder, *error) {
	switch contentType {
	case "application/json", "":
		return func(r io.Reader) Decoder { return json.NewDecoder(r) }, nil
	case "application/xml", "text/xml":
		return func(r io.Reader) Decoder { return xml.NewDecoder(r) }, nil
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return func(r io.Reader) Decoder { return form.NewDecoder(r) }, nil
	default:
		return nil, nil
	}
}

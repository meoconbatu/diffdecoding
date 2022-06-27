package diff

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
)

// decode func performs string decoding, then return.
// Supported encoding types are: gz, gzip, gz+base64, gzip+base64, gz+b64, gzip+b64, b64, base64, text/plain
// if encondingType is 'text/plain', just return s
func decode(s, encodingType string) (string, error) {
	var err error
	v := []byte(s)
	switch encodingType {
	case "b64", "base64":
		v, err = base64Decode(s)
	case "gz", "gzip":
		v, err = gunzipData([]byte(s))
	case "text/plain":
	default:
		v, err = base64DecodeGunzip(s)
	}
	return string(v), err
}

// base64DecodeGunzip decodes a string containing a base64 sequence and then uncompresses the result with gzip
func base64DecodeGunzip(data string) ([]byte, error) {
	resData, err := base64Decode(data)
	if err != nil {
		return nil, errors.New("Invalid character in input stream")
	}
	resData, err = gunzipData(resData)
	if err != nil {
		return nil, errors.New("Unknown compression format")
	}
	return resData, nil
}

// base64Decode decodes a string containing a base64 sequence.
// Terraform uses the "standard" Base64 alphabet as defined in RFC 4648 section 4.
func base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
func gzipData(data []byte) (compressedData []byte, err error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	_, err = w.Write(data)
	if err != nil {
		return
	}
	err = w.Flush()
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}
	compressedData = buf.Bytes()
	return
}

// gunzip decompress a gzip compressed string
func gunzipData(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return data, err
	}
	var resBuf bytes.Buffer
	_, err = resBuf.ReadFrom(r)
	if err != nil {
		return data, err
	}
	return resBuf.Bytes(), nil
}

package diff

import (
	"bytes"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
)

type part struct {
	header textproto.MIMEHeader
	body   []byte
}

func (p part) isYAML() bool {
	return p.header.Get("Content-Type") == "text/cloud-config"
}
func parse(b []byte) (map[string]string, []*part) {
	msg, _ := mail.ReadMessage(bytes.NewBuffer(b))
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	parts := make([]*part, 0)
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return params, parts
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := io.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(p.Header)
			parts = append(parts, &part{p.Header, slurp})
			// fmt.Printf("Part %q: %q\n", p.Header.Get("Content-Type"), slurp)
		}
	}
	return nil, nil
}

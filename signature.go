package gotk

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type SigningMd5 struct {
	secrete, key string
	lowcase      bool
}

func NewSigningMd5(secrete, key string, lowcase bool) (*SigningMd5, error) {
	if secrete == "" || key == "" {
		return nil, fmt.Errorf("secrete or key is empty")
	}
	return &SigningMd5{secrete: secrete, key: key, lowcase: lowcase}, nil
}

func (sign *SigningMd5) SignValues(param map[string]string) (value string) {
	var (
		k, format string
		keys      []string
		pairs     []string
	)

	keys = make([]string, 0, len(param))

	for k = range param {
		if k != sign.key { // !!
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	pairs = make([]string, 0, len(param))
	for _, k = range keys {
		pairs = append(pairs, k+param[k])
	}

	if format = "%X"; sign.lowcase {
		format = "%x"
	}

	value = fmt.Sprintf(format, md5.Sum([]byte(sign.secrete+strings.Join(pairs, "")+sign.secrete)))

	return value
}

func (sign *SigningMd5) SignQuery(param map[string]string) (query string) {
	var (
		value string
		pairs []string
	)

	value = sign.SignValues(param)
	pairs = make([]string, 0, len(param)+1)

	for k, v := range param {
		pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	pairs = append(pairs, url.QueryEscape(sign.key)+"="+url.QueryEscape(value))

	return strings.Join(pairs, "&")
}

func (sign *SigningMd5) VerifyQuery(query string) (err error) {
	var (
		value  string
		param  map[string]string
		values url.Values
	)

	if values, err = url.ParseQuery(query); err != nil {
		return err
	}

	if value = values.Get(sign.key); value == "" || len(value) != 32 {
		return fmt.Errorf("invalid signature")
	}

	param = make(map[string]string, len(values)-1)

	for k := range values {
		if k == sign.key {
			continue
		}
		if len(values[k]) == 0 {
			param[k] = ""
		} else {
			param[k] = values[k][0]
		}
	}

	if value != sign.SignValues(param) {
		return fmt.Errorf("signature not match")
	}

	return nil
}

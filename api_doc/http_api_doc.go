package api_doc

import (
	// "fmt"
	"time"
)

type Field struct {
	Type     string `mapstructure:"type"` // int32, int64, float32, bool, []xx
	Key      string `mapstructure:"key"`
	Value    string `mapstructure:"value"` // can be deserialize to Type
	Required bool   `mapstructure:"required"`
	Note     string `mapstructure:"note"`
}

type HttpAPI struct {
	Title           string        `mapstructure:"title"`
	Note            string        `mapstructure:"note"`
	Interval        time.Duration `mapstructure:"interval"`
	Request         string        `mapstructure:"request"` // method@path
	Headers         []Field       `mapstructure:"headers"`
	Parameters      []Field       `mapstructure:"parameters"`
	Body            []Field       `mapstructure:"body"`
	ResponseHeaders []Field       `mapstructure:"response_headers"`
	ResponseBody    string        `mapstructure:"response_body"`
}

type HttpAPIs struct {
	Name      string   `mapstructure:"name"`
	Links     []string `mapstructure:"links"`
	Note      string   `mapstructure:"note"`
	Basepath  string   `mapstructure:"basepath"`
	Signature string   `mapstructure:"signature"` // TODO:
	Public    struct {
		Headers []Field `mapstructure:"headers"` // without subfield Type
		Body    []Field `mapstructure:"body"`
	} `mapstructure:"public"`
	List []HttpAPI `mapstructure:"list"`
}

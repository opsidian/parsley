package json_test

import (
	"io/ioutil"
	"testing"

	encoding_json "encoding/json"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/examples/json/json"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func benchmarkParsleyJSON(b *testing.B, jsonFilePath string) {
	input, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		b.Fatal(err)
	}

	s := combinator.Sentence(json.NewParser())
	f := text.NewFile("test", input[0:len(input)-1])
	r := text.NewReader(f)
	h := parser.NewHistory()
	if _, err = parsley.Evaluate(h, r, s, nil); err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		h := parser.NewHistory()
		_, _ = parsley.Evaluate(h, r, s, nil)
	}
}

func BenchmarkParsleyJSON1k(b *testing.B)  { benchmarkParsleyJSON(b, "../example_1k.json") }
func BenchmarkParsleyJSON10k(b *testing.B) { benchmarkParsleyJSON(b, "../example_10k.json") }

func benchmarkEncodingJSON(b *testing.B, jsonFilePath string) {
	input, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		b.Fatal(err)
	}

	var val interface{}
	if err := encoding_json.Unmarshal(input, &val); err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		var val interface{}
		_ = encoding_json.Unmarshal(input, &val)
	}
}

func BenchmarkEncodingJSON1k(b *testing.B)  { benchmarkEncodingJSON(b, "../example_1k.json") }
func BenchmarkEncodingJSON10k(b *testing.B) { benchmarkEncodingJSON(b, "../example_10k.json") }

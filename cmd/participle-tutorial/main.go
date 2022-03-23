package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/rs/zerolog/log"
)

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

// @ captures into a field
// @@ captures recursively

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	String  *string  ` @String`
	Number  *float64 `| @Float`
	Integer *int64   `| @Int`
}

func main() {
	parser, err := participle.Build(&INI{})
	if err != nil {
		log.Error().Err(err).Send()
	}
	ini := &INI{}
	err = parser.ParseString("", `
age = 21.34
age = 21
name = "Bob Smith"

[address]
city = "Beverly Hills"
postal_code = 90210
`, ini)
	if err != nil {
		log.Error().Err(err).Send()
	}

	log.Info().Interface("ini", ini).Send()

	err = parser.ParseString("", `
foo foo foo
name = "Bob Smith"

[address]
city = "Beverly Hills"
postal_code = 90210
`, ini)
	if err != nil {
		log.Error().Err(err).Send()
	}

	log.Info().Interface("ini", ini).Send()
}

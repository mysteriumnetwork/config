package config

import (
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// delim is the delimeter used to parse nested structs.
// Its used in koanf.New call.
var delim = "."

// Parser exposes a `Parse` method to load data from config files.
type Parser[T any] struct {
	files []string
}

// NewParser returns a new parser.
func NewParser[T any](configFiles []string) *Parser[T] {
	return &Parser[T]{
		files: configFiles,
	}
}

// Parse will read all possible files loading them to provided `cfg` in order:
// 1. Defaults
// 2. Files
// Each step overrides the last.
func (p *Parser[T]) Parse(cfg *T, defaults T) error {
	k := koanf.New(delim)

	err := k.Load(structs.Provider(defaults, "koanf"), nil)
	if err != nil {
		return fmt.Errorf("error loading defaults when loading config: %w", err)
	}

	for _, path := range p.files {
		parser, err := getParser(path)
		if err != nil {
			return fmt.Errorf("error getting parser when loading config: %w", err)
		}

		if err := k.Load(file.Provider(path), parser); err != nil {
			return fmt.Errorf("error reading file when loading config: %w", err)
		}
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}

// SetDelim will set the global delimiter used for nested structs in config.
func SetDelim(d string) {
	delim = d
}

func getParser(file string) (koanf.Parser, error) {
	ext := filepath.Ext(file)
	switch ext {
	case ".json":
		return json.Parser(), nil
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	default:
		return nil, fmt.Errorf("%s file format is not supported, require either yaml or json", ext)
	}
}

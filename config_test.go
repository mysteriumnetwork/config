package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exampleConfig struct {
	Field1  string  `koanf:"field1"`
	Field2  int     `koanf:"field2"`
	Struct1 struct1 `koanf:"struct1"`
	Struct2 struct2 `koanf:"struct2"`
}

type struct1 struct {
	SubField1 float64  `koanf:"subfield1"`
	SubField2 []string `koanf:"subfield2"`
}

type struct2 struct {
	SubField3 map[string]string `koanf:"subfield3"`
	SubField4 bool              `koanf:"subfield4"`
}

func TestConfig(t *testing.T) {
	defaults := exampleConfig{
		Field1: "default1",
		Field2: 2,
		Struct1: struct1{
			SubField1: 6.28,
			SubField2: []string{"default1", "default2"},
		},
		Struct2: struct2{
			SubField3: map[string]string{
				"defaultkey1": "defaultvalue1",
				"defaultkey2": "defaultvalue2",
			},
			SubField4: true,
		},
	}

	var cfg exampleConfig
	err := Parse([]string{"./example_config.yaml"}, &cfg, defaults)
	assert.NoError(t, err)

	assert.Equal(t, "value1", cfg.Field1)
	assert.Equal(t, 2, cfg.Field2)
	assert.Equal(t, 3.14, cfg.Struct1.SubField1)
	require.Len(t, cfg.Struct1.SubField2, 2)
	assert.Equal(t, "default1", cfg.Struct1.SubField2[0])
	assert.Equal(t, "default2", cfg.Struct1.SubField2[1])
	require.Len(t, cfg.Struct2.SubField3, 2)
	assert.Equal(t, "value2", cfg.Struct2.SubField3["key1"])
	assert.Equal(t, "value3", cfg.Struct2.SubField3["key2"])
	assert.True(t, cfg.Struct2.SubField4)
}

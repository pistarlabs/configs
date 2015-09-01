package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Config represents a configuration with convenient access methods.
type Config struct {
	Root interface{}
}

// String returns a string according to a dotted path.
func (c *Config) String(path string) (string, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return "", err
	}

	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n), nil
	case string:
		return n, nil
	}

	return "", fmt.Errorf("Type mismatch: expected string; got %T", n)
}

// UString returns a string according to a dotted path or default or "".
func (c *Config) UString(path string, defaults ...string) string {
	value, err := c.String(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return ""
}

// Bool returns a bool according to a dotted path.
func (c *Config) Bool(path string) (bool, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return false, err
	}

	switch n := n.(type) {
	case bool:
		return n, nil
	case string:
		return strconv.ParseBool(n)
	}

	return false, fmt.Errorf("Type mismatch: expected bool; got %T", n)
}

// UBool retirns a bool according to a dotted path or default value or false.
func (c *Config) UBool(path string, defaults ...bool) bool {
	value, err := c.Bool(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return false
}

// Float64 returns a float64 according to a dotted path.
func (c *Config) Float64(path string) (float64, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return 0, err
	}

	switch n := n.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case string:
		return strconv.ParseFloat(n, 64)
	}

	return 0, fmt.Errorf("Type mismatch: expected float64; got %T", n)
}

// UFloat64 returns a float64 according to a dotted path or default value or 0.
func (c *Config) UFloat64(path string, defaults ...float64) float64 {
	value, err := c.Float64(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return float64(0)
}

// Int returns an int according to a dotted path.
func (c *Config) Int(path string) (int, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return 0, err
	}
	switch n := n.(type) {
	case float64:
		// encoding/json unmarshals numbers into floats, so we compare
		// the string representation to see if we can return an int.
		if i := int(n); fmt.Sprint(i) == fmt.Sprint(n) {
			return i, nil
		}
		return 0, fmt.Errorf("Value can't be converted to int: %v", n)
	case int:
		return n, nil
	case string:
		if v, err := strconv.ParseInt(n, 10, 0); err == nil {
			return int(v), nil
		}
		return 0, err
	}
	return 0, fmt.Errorf("Type mismatch: expected int; got %T", n)
}

// UInt returns an int according to a dotted path or default value or 0.
func (c *Config) UInt(path string, defaults ...int) int {
	value, err := c.Int(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return 0
}

// List returns a []interface{} according to a dotted path.
func (c *Config) List(path string) ([]interface{}, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return nil, err
	}

	if value, ok := n.([]interface{}); ok {
		return value, nil
	}

	return nil, fmt.Errorf("Type mismatch: expected []interface{}; got %T", n)
}

// UList returns a []interface{} according to a dotted path or defaults or []interface{}.
func (c *Config) UList(path string, defaults ...[]interface{}) []interface{} {
	value, err := c.List(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}
	return make([]interface{}, 0)
}

// Map returns a map[string]interface{} according to a dotted path.
func (c *Config) Map(path string) (map[string]interface{}, error) {
	n, err := Get(c.Root, path)
	if err != nil {
		return nil, err
	}

	if value, ok := n.(map[string]interface{}); ok {
		return value, nil
	}

	return nil, fmt.Errorf("Type mismatch: expected map[string]interface{}; got %T", n)
}

// UMap returns a map[string]interface{} according to a dotted path or default or map[string]interface{}.
func (c *Config) UMap(path string, defaults ...map[string]interface{}) map[string]interface{} {
	value, err := c.Map(path)

	if err == nil {
		return value
	}

	for _, def := range defaults {
		return def
	}

	return map[string]interface{}{}
}

// Get returns a child of the given value according to a dotted path.
func Get(cfg interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	// Normalize path.
	for k, v := range parts {
		if v == "" {
			if k == 0 {
				parts = parts[1:]
			} else {
				return nil, fmt.Errorf("Invalid path %q", path)
			}
		}
	}
	// Get the value.
	for pos, part := range parts {
		switch c := cfg.(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				if int(i) < len(c) {
					cfg = c[i]
				} else {
					return nil, fmt.Errorf(
						"Index out of range at %q: list has only %v items",
						strings.Join(parts[:pos+1], "."), len(c))
				}
			} else {
				return nil, fmt.Errorf("Invalid list index at %q",
					strings.Join(parts[:pos+1], "."))
			}

		case map[string]interface{}:
			if value, ok := c[part]; ok {
				cfg = value
			} else {
				return nil, fmt.Errorf("Nonexistent map key at %q",
					strings.Join(parts[:pos+1], "."))
			}

		default:
			return nil, fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(parts[:pos+1], "."), cfg)

		}
	}

	return cfg, nil
}

// Load reads a JSON configuration from given filename
func Load(filename string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parse(cfg)
}

// parse perform JSON parsing
func parse(cfg []byte) (*Config, error) {
	var out interface{}
	var err error

	if err = json.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}

	if out, err = normalize(out); err != nil {
		return nil, err
	}

	return &Config{Root: out}, nil
}

// normalizes a unmarshalled value. This is needed because
// encoding/json doesn't support marshalling map[interface{}]interface{}.
func normalize(value interface{}) (interface{}, error) {
	switch value := value.(type) {

	case map[interface{}]interface{}:
		node := make(map[string]interface{}, len(value))

		for k, v := range value {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Unsupported map key: %#v", k)
			}

			item, err := normalize(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}

			node[key] = item
		}

		return node, nil

	case map[string]interface{}:
		node := make(map[string]interface{}, len(value))

		for key, v := range value {

			item, err := normalize(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}

			node[key] = item
		}

		return node, nil

	case []interface{}:
		node := make([]interface{}, len(value))
		for key, v := range value {
			item, err := normalize(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported list item: %#v", v)
			}
			node[key] = item
		}
		return node, nil

	case bool, float64, int, string:
		return value, nil

	}

	return nil, fmt.Errorf("Unsupported type: %T", value)
}

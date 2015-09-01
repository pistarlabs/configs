package configs

import (
	"reflect"
	"testing"
)

var configTests = []struct {
	path string
	kind string
	want interface{}
	ok   bool
}{
	// OK
	{"development.database.host", "String", "localhost", true},
	// Failed
	{"development.database.something", "String", "", false},
}

func TestLoadJson(t *testing.T) {
	_, err := Load("config.json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFunction(t *testing.T) {
	cfg, err := Load("config.json")
	if err != nil {
		t.Fatal(err)
	}

	expect(t, cfg.UString("development.database.host"), "localhost")
	expect(t, cfg.UString("development.database.username"), "root")
	expect(t, cfg.UString("development.database.password"), "12345")
	expect(t, cfg.UInt("development.database.port"), 12345)
	expect(t, cfg.UString("development.database.name"), "dev")

	expect(t, cfg.UString("production.database.host"), "localhost")
	expect(t, cfg.UString("production.database.username"), "root")
	expect(t, cfg.UString("production.database.password"), "12345")
	expect(t, cfg.UInt("production.database.port"), 12345)
	expect(t, cfg.UString("production.database.name"), "dev")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

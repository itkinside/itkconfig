package itkconfig

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/string.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse string value: %s", err.Error())
	}

	want := Config{
		Foo: "bar",
	}
	if want != config {
		t.Fatalf(`
Could not parse config containing string.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestBool(t *testing.T) {
	type Config struct {
		Foo bool
	}

	config := Config{
		Foo: false,
	}
	err := LoadConfig("test_configs/bool.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse string value: %s", err.Error())
	}

	want := Config{
		Foo: true,
	}
	if want != config {
		t.Fatalf(`
Could not parse config containing bool.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestInt(t *testing.T) {
	type Config struct {
		Foo int
	}

	config := Config{
		Foo: 0,
	}
	err := LoadConfig("test_configs/int.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse string value: %s", err.Error())
	}

	want := Config{
		Foo: -1,
	}
	if want != config {
		t.Fatalf(`
Could not parse config containing int.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestUint(t *testing.T) {
	type Config struct {
		Foo uint
	}

	config := Config{
		Foo: 0,
	}
	err := LoadConfig("test_configs/uint.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse string value: %s", err.Error())
	}

	want := Config{
		Foo: 1,
	}
	if want != config {
		t.Fatalf(`
Could not parse config containing uint.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestFloat(t *testing.T) {
	type Config struct {
		Foo float32
	}

	config := Config{
		Foo: 0,
	}
	err := LoadConfig("test_configs/float.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse string value: %s", err.Error())
	}

	want := Config{
		Foo: 0.5,
	}
	if want != config {
		t.Fatalf(`
Could not parse config containing float.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestConfigWrongType(t *testing.T) {
	a := 5
	err := LoadConfig("test_configs/empty.cfg", &a)

	if err == nil {
		t.Fatal("Parsed config with invalid type.")
	}
}

func TestConfigWithoutPointer(t *testing.T) {
	type Config struct{}
	err := LoadConfig("test_configs/empty.cfg", Config{})
	if err == nil {
		t.Fatal("Parsed config without config being a pointer.")
	}
}

func TestLoadingEmptyConfig(t *testing.T) {
	type Config struct{}
	err := LoadConfig("test_configs/empty.cfg", &Config{})
	if err != nil {
		t.Fatalf("Could not parse empty config. Err: %s", err.Error())
	}
}

func TestLoadingExampleConfig(t *testing.T) {
	type Config struct {
		Port            int
		TemplatesFolder string
		Debug           bool
		AdminEmail      []string
	}
	config := Config{
		Port:            8080,
		TemplatesFolder: "tmpl",
		Debug:           false,
		AdminEmail:      []string{"test@example.org"},
	}

	want := Config{
		Port:            8000,
		TemplatesFolder: "templates",
		Debug:           true,
		AdminEmail:      []string{"test@example.org", "foo@mailinator.com", "bar@mailinator.com"}, // Slices append to defaults
	}

	err := LoadConfig("test_configs/example.cfg", &config)
	if err != nil {
		t.Fatal("Could not parse example.cfg")
	}

	if !reflect.DeepEqual(want, config) {
		t.Fatalf(`
Example config not equal 
	expected: %#v, 
	got:      %#v
`, want, config)
	}
}

func TestConfigSpecifiesNonExistent(t *testing.T) {
	type Config struct{}
	err := LoadConfig("test_configs/example.cfg", &Config{})
	if err == nil {
		t.Fatal("Parsed config even though the field was not requested.")
	}
}

func TestQuotedValues(t *testing.T) {
	type Config struct {
		Foo string
	}
	config := Config{
		Foo: "baz",
	}

	err := LoadConfig("test_configs/quotedvalues.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config: %s", err.Error())
	}

	if config.Foo != "bar" {
		t.Fatalf("Parsed config incorrectly. Expected: 'bar', got: '%s'.", config.Foo)
	}
}

func TestUnexportedValues(t *testing.T) {
	type Config struct {
		unexported string
	}

	err := LoadConfig("test_configs/unexported.cfg", &Config{
		unexported: "foo",
	})

	if err == nil {
		t.Fatal("Loading config that sets unexported values should error.")
	}
}

func TestNoSpace(t *testing.T) {
	type Config struct {
		Foo int
	}

	config := Config{
		Foo: 2,
	}
	err := LoadConfig("test_configs/nospaces.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with no spaces: %s", err.Error())
	}
}

func TestWierdSpaces(t *testing.T) {
	type Config struct {
		Foo int
	}

	config := Config{
		Foo: 2,
	}
	err := LoadConfig("test_configs/nospaces.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with no spaces: %s", err.Error())
	}
	if config.Foo != 1 {
		t.Fatalf("Config parsed incorrectly. Expected: %d, got: %d", 1, config.Foo)
	}
}

func TestWeirdQuote(t *testing.T) {
	type Config struct {
		Foo string
		Bar string
	}

	config := Config{
		Foo: "",
		Bar: "",
	}
	err := LoadConfig("test_configs/quotesinside.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with quotes inside of value: %s", err.Error())
	}

	want := Config{
		Foo: "quote",
		Bar: "string",
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with quotes inside correctly.
	expected: %#v
	got:      %#v`, want, config)
	}
}
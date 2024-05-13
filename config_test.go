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
		AdminEmail:      []string{"foo@mailinator.com", "bar@mailinator.com"},
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

func TestStringWithComment(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/stringwithcomment.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with escaped quote: %s", err.Error())
	}

	want := Config{
		Foo: "#something",
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with comment inside of string.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestEndOfLineComment(t *testing.T) {
	type Config struct {
		Foo int
	}

	config := Config{
		Foo: 0,
	}
	err := LoadConfig("test_configs/endoflinecomment.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with escaped quote: %s", err.Error())
	}

	want := Config{
		Foo: 1,
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with comment at the end of the line.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestFullLineComment(t *testing.T) {
	type Config struct{}
	err := LoadConfig("test_configs/fulllinecomment.cfg", &Config{})
	if err != nil {
		t.Fatalf("Could not parse config with full line comment: %s", err.Error())
	}
}

func TestKeyWithQuote(t *testing.T) {
	type Config struct {
		Key string
	}
	err := LoadConfig("test_configs/keywithquote.cfg", &Config{Key: "test"})
	if err == nil {
		t.Fatal("Key with quote should not be allowed.")
	}
}

func TestNoEquals(t *testing.T) {
	type Config struct {
		Foo string
	}
	err := LoadConfig("test_configs/noequals.cfg", &Config{Foo: ""})
	if err == nil {
		t.Fatal("Key with no equals should not be allowed.")
	}
}

func TestEmptyKey(t *testing.T) {
	type Config struct{}
	err := LoadConfig("test_configs/emptykey.cfg", &Config{})
	if err == nil {
		t.Fatal("Empty key is not allowed.")
	}
}

func TestWierdEdgeCase(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/wierdedgecase.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with wierd syntax: %s", err.Error())
	}

	want := Config{
		Foo: "hello",
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with wierd syntax correctly.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestEscapedQuote(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/escapedquote.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with escaped quote: %s", err.Error())
	}

	want := Config{
		Foo: "hel\"lo",
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with quotes inside correctly.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestTwoQuoteStart(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/twoquotesstart.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with escaped quote: %s", err.Error())
	}

	want := Config{
		Foo: "test",
	}
	if want != config {
		t.Fatalf(`
Could not parse config correctly with quotes inside correctly.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestTwoQuote(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/twoquotes.cfg", &config)
	if err != nil {
		t.Fatalf("Could not parse config with escaped quote: %s", err.Error())
	}

	want := Config{
		Foo: "testing",
	}
	if want != config {
		t.Fatalf(`
Could not parse config with multiple quotes correctly.
	expected: %#v
	got:      %#v`, want, config)
	}
}

func TestMultipleDefinitions(t *testing.T) {
	type Config struct {
		Foo string
	}

	config := Config{
		Foo: "",
	}
	err := LoadConfig("test_configs/multipledefinitions.cfg", &config)
	if err == nil {
		t.Fatalf("Should not be allowed to redefine a key.")
	}
}

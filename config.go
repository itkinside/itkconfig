// Copyright (c) 2014 Trygve Aaberge and contributors
// Released under the LGPLv2.1, see LICENSE

// Package itkconfig implements parsing of configuration files through the use
// of reflection.
package itkconfig

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// parseField parses a field based on its field type.
func parseField(key, value string, fieldType reflect.Type) (reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.ValueOf(nil), fmt.Errorf("invalid bool \"%s\" in key \"%s\": %s", value, key, err)
		}
		return reflect.ValueOf(v), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, fieldType.Bits())
		if err != nil {
			return reflect.ValueOf(nil), fmt.Errorf("invalid int \"%s\" in key \"%s\": %s", value, key, err)
		}
		return reflect.ValueOf(i).Convert(fieldType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(value, 10, fieldType.Bits())
		if err != nil {
			return reflect.ValueOf(nil), fmt.Errorf("invalid uint \"%s\" in key \"%s\": %s", value, key, err)
		}
		return reflect.ValueOf(i).Convert(fieldType), nil
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, fieldType.Bits())
		if err != nil {
			return reflect.ValueOf(nil), fmt.Errorf("invalid float \"%s\" in key \"%s\": %s", value, key, err)
		}
		return reflect.ValueOf(i).Convert(fieldType), nil
	default:
		return reflect.ValueOf(nil), fmt.Errorf("unsupported type: %s", fieldType.Kind())
	}
}

func parseKey(rawKey string) (*string, error) {
	key := strings.TrimSpace(rawKey)
	if strings.Contains(key, "\"") {
		return nil, errors.New("key cannot contain \"")
	}
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}
	return &key, nil
}

func parseVal(rawVal string) (*string, error) {
	val := strings.TrimSpace(rawVal)

	quoteCommentGroup := regexp.MustCompile(`^(".*?"|[^"]*?)(\s*#.*)$`)
	groups := quoteCommentGroup.FindStringSubmatchIndex(val)
	if groups != nil {
		val = val[:groups[2*2]]
	}

	// Remove non-escaped quotes and replace escaped quotes.
	var sb strings.Builder
	for i, r := range val {
		if r == '"' {
			continue
		}

		if val[i] == '\\' && val[i+1] == '"' {
			sb.WriteRune('"')
		} else {
			sb.WriteRune(r)
		}
	}
	val = sb.String()

	return &val, nil
}

// LoadConfig loads the provided configuration file and parses it through the
// use of reflection according to the type definition of config, which has to be
// a pointer to a struct.
func LoadConfig(filename string, config interface{}) error {
	// Use reflect to place config keys into the right element in the struct
	configPtrReflect := reflect.ValueOf(config)
	if configPtrReflect.Kind() != reflect.Ptr {
		return errors.New("config argument must be a pointer")
	}
	configReflect := configPtrReflect.Elem()
	if configReflect.Kind() != reflect.Struct {
		return errors.New("config argument must be a pointer to a struct")
	}

	lastUpdate := make(map[string]uint)
	for _, field := range reflect.VisibleFields(configReflect.Type()) {
		lastUpdate[field.Name] = 0
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fh := bufio.NewScanner(f)

	lineNr := uint(0)
	syntaxError := func(message string) error {
		return fmt.Errorf("syntax error parsing config (%s:%d): %s", filename, lineNr, message)
	}

	for fh.Scan() {
		line := fh.Text()
		lineNr++

		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		keyVal := strings.SplitN(line, "=", 2)
		if len(keyVal) != 2 {
			return syntaxError("line must contain '='")
		}

		key, err := parseKey(keyVal[0])
		if err != nil {
			return syntaxError(err.Error())
		}

		value, err := parseVal(keyVal[1])
		if err != nil {
			return syntaxError(err.Error())
		}

		field := configReflect.FieldByName(*key)
		if !field.IsValid() {
			return syntaxError(fmt.Sprintf("the config key '%s' is not defined", *key))
		}
		if !field.CanSet() {
			return syntaxError(fmt.Sprintf("cannot set unexported field: '%s'", *key))
		}

		switch field.Kind() {
		case reflect.Slice:
			if lastUpdate[*key] == 0 {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			}

			v, err := parseField(*key, *value, field.Type().Elem())
			if err != nil {
				return syntaxError(err.Error())
			}

			field.Set(reflect.Append(field, v))
		default:
			if lastUpdate[*key] != 0 {
				return syntaxError(fmt.Sprintf("key '%s' was defined multiple times, initially on line %d (did you mean to define a slice?)", *key, lastUpdate[*key]))
			}

			v, err := parseField(*key, *value, field.Type())
			if err != nil {
				return syntaxError(err.Error())
			}
			field.Set(v)
		}
		lastUpdate[*key] = lineNr
	}

	return nil
}

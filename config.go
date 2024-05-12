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

	val = strings.ReplaceAll(val, "\"", "")
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

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fh := bufio.NewScanner(f)

	lineNr := 0
	for fh.Scan() {
		line := fh.Text()
		lineNr++

		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		keyVal := strings.SplitN(line, "=", 2)
		if len(keyVal) != 2 {
			return fmt.Errorf("syntax error parsing config (%s:%d): config line must contain '='", filename, lineNr)
		}

		key, err := parseKey(keyVal[0])
		if err != nil {
			return fmt.Errorf("syntax error parsing config (%s:%d): %s", filename, lineNr, err.Error())
		}

		value, err := parseVal(keyVal[1])
		if err != nil {
			return fmt.Errorf("syntax error parsing config (%s:%d): %s", filename, lineNr, err.Error())
		}

		field := configReflect.FieldByName(*key)
		if !field.IsValid() {
			return fmt.Errorf("syntax error parsing config (%s:%d): config key is not valid: '%s'", filename, lineNr, *key)
		}
		if !field.CanSet() {
			return fmt.Errorf("syntax error parsing config (%s:%d): cannot set unexported field: '%s'", filename, lineNr, *key)
		}

		switch field.Kind() {
		case reflect.Slice:
			// Create a empty slice, if no slice exists for this key already.
			if field.IsNil() {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			}

			// Convert the value (string) to Value struct defined in reflect.
			v, err := parseField(*key, *value, field.Type().Elem())
			if err != nil {
				return fmt.Errorf("syntax error parsing config (%s:%d): %s", filename, lineNr, err.Error())
			}

			field.Set(reflect.Append(field, v))
		default:
			v, err := parseField(*key, *value, field.Type())
			if err != nil {
				return fmt.Errorf("syntax error parsing config (%s:%d): %s", filename, lineNr, err.Error())
			}
			field.Set(v)
		}
	}

	return nil
}

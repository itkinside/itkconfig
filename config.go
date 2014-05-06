// Copyright (c) 2014 Trygve Aaberge
// Released under the MIT License, http://opensource.org/licenses/MIT

package itkconfig

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func LoadConfig(filename string, config interface{}) error {
	// Use reflect to place config keys into the right element in the struct
	configPtrReflect := reflect.ValueOf(config)
	if configPtrReflect.Kind() != reflect.Ptr {
		return errors.New("Config argument must be a pointer")
	}
	configReflect := configPtrReflect.Elem()
	if configReflect.Kind() != reflect.Struct {
		return errors.New("Config argument must be a pointer to a struct")
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	fh := bufio.NewScanner(f)

	for fh.Scan() {
		var key, value string
		line := fh.Text()
		lineParts := strings.Split(line, "\"")
		// Split the line on ", because we want to keep parts
		// inside "" unchanged.
		for i, part := range lineParts {
			if i%2 == 0 {
				commentIndex := strings.Index(part, "#")
				if commentIndex != -1 {
					// Remove comments
					part = part[:commentIndex]
				}
				if i == 0 {
					// If first part, we want to fetch the key
					keyVal := strings.SplitN(part, "=", 2)
					key = strings.TrimSpace(keyVal[0])
					if len(keyVal) < 2 {
						// The part didn't contain a =
						if i != len(lineParts)-1 && commentIndex == -1 {
							// Not the last line, which means there is a " before the = (if any)
							return fmt.Errorf("\" are not allowed in key: %s", line)
						}
						if key != "" {
							// Last line, which means no =
							return fmt.Errorf("Config line must contain \"=\": %s", line)
						}
						// The line is only comments
						break
					} else if key == "" {
						// Line had a =, but only spaces before it
						return fmt.Errorf("Key can't be empty: %s", line)
					}
					// We want to trim space at the start of the value
					part = strings.TrimLeftFunc(keyVal[1], unicode.IsSpace)
				}
				if i == len(lineParts)-1 || commentIndex != -1 {
					// Last part, we want to trim space at the end of the value
					part = strings.TrimRightFunc(part, unicode.IsSpace)
				}
				if commentIndex != -1 {
					// The part had a comment char, ignore the rest of the parts
					value += part
					break
				}
			}
			value += part
		}

		if key == "" {
			// The line is only comments
			continue
		}

		if value == "" {
			return fmt.Errorf("Value of key \"%s\" can't be empty.", key)
		}

		// Fetch the field in the config struct with the same name as the key
		field := configReflect.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("Config key is not valid: %s", key)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid int \"%s\" in key \"%s\": %s", value, key, err)
			}
			field.SetInt(i)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("Invalid bool \"%s\" in key \"%s\": %s", value, key, err)
			}
			field.SetBool(v)
		default:
			return fmt.Errorf("Unsupported type: %s", field.Kind())
		}
	}

	return nil
}

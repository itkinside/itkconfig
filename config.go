// Copyright (c) 2014 Trygve Aaberge
// Released under the MIT License, http://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Config struct {
	Server  string
	Channel string
	Nr      int
	IPv6    bool
}

func loadConfig(filename string) (config Config, err error) {
	return config, loadConfigWithDefaults(filename, &config)
}

func loadConfigWithDefaults(filename string, defaultConfig *Config) error {
	configReflect := reflect.ValueOf(defaultConfig).Elem()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	fh := bufio.NewScanner(f)

	for fh.Scan() {
		var key, value string
		lineParts := strings.Split(fh.Text(), "\"")
		for i, part := range lineParts {
			if i%2 == 0 {
				commentIndex := strings.Index(part, "#")
				if commentIndex != -1 {
					part = part[:commentIndex]
				}
				if i == 0 {
					keyVal := strings.SplitN(part, "=", 2)
					key = strings.TrimSpace(keyVal[0])
					if len(keyVal) < 2 {
						break
					}
					part = strings.TrimLeftFunc(keyVal[1], unicode.IsSpace)
				}
				if i == len(lineParts)-1 || commentIndex != -1 {
					part = strings.TrimRightFunc(part, unicode.IsSpace)
				}
				if commentIndex != -1 {
					value += part
					break
				}
			}
			value += part
		}

		if key == "" {
			continue
		}

		if value == "" {
			return fmt.Errorf("Value of key \"%s\" can't be empty.", key)
		}

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

func main() {
	config, err := loadConfig("itkbot.conf")
	if err != nil {
		fmt.Println("Couldn't load config:", err)
		return
	}

	fmt.Println("Server:", config.Server)
	fmt.Println("Channel:", config.Channel)
	fmt.Println("Nr:", config.Nr)
	fmt.Println("IPv6:", config.IPv6)
}

# ITKconfig - A small, powerful configuration parser for Golang

Originally started as an internal project at the [Student society in
Trondheim](http://samfundet.no) this package has now been open sourced, as we
believe it is the simplest and best way to manage configuration files for
Go-projects. It serves its purpose for our projects, but we would love to hear
your use-cases and feedback, if any.

## Features and core-principles

* Makes writing **Key-Value** configuration files easy.
* Allows, in contrast to JSON, comments in your files - just prepend them with a
  `#`.
* If you want to use `#` in a value, or preserve leading and trailing spaces,
  wrap the value in double quotes: `"#value"`
* Use the same methods as when demarshalling JSON-files, just define your
  configuration struct with your wanted types and let ITKconfig take care of the
  rest.
* Source code is simple and short, which makes it easy to understand the flow
  of the program, but also make changes to the library if you like.

## Example configuration file

An example scenario is given where you want to provide a configuration file to
your Web-application. It could look like:

```bash
# Port that the webservice is listening to
Port = 8000

# Folder where we find our templates
TemplatesFolder = templates

# Enable or disable debug mode, giving more output to the user.
Debug = true

# Various contact points for the admins
AdminEmail = foo@mailinator.com
AdminEmail = bar@mailinator.com
```

Then, provided that this file is called `myapp.config` we can load it into our
application by the following simple code:

```go
package main

import (
  "fmt"
  "github.com/itkinside/itkconfig"
  "log"
)

type Config struct {
  Port            int
  TemplatesFolder string
  Debug           bool
  AdminEmail      []string
}

func main() {
  // Some sane defaults for our project.
  config := &Config{
    Port:            80,
    TemplatesFolder: "temps",
    Debug:           false,
    AdminEmail:      []string{"admin@mailinator.com"},
  }

  // Override (or append on) defaults with config-file.
  err := itkconfig.LoadConfig("myapp.config", config)
  if err != nil {
    log.Fatal(err)
  }

  // Print our variables, just to show off.
  fmt.Printf("Port: %d\n", config.Port)
  fmt.Printf("Templates: %s\n", config.TemplatesFolder)
  fmt.Printf("Debug: %v\n", config.Debug)
  for i, email := range config.AdminEmail {
    fmt.Printf("Admin email %d: %s\n", i, email)
  }
}
```

Could it be more simple, and yet so powerful?

## Some useful tips

#### Comments

The hash symbol is your friend, and you can use it wherever you want. You may
also use it inside a variable by escaping it:

```bash
# This is a comment
Key = some value # Also a comment
Foo = "#something" # This is first comment on this line.
```

#### String parsing

Double quotes are removed when parsing strings.

This is the case when writing a string without surrounding them with quotes:
```bash
Foo = ba"r
# Gets parsed as "bar"
```

Or when writing strings with quotes surrounding them:
```bash
Foo = "ba"r"
# Gets parsed as "bar"
```

If a string containing double quotes is desired, they can be escaped:
```bash
Foo = "ba\"r"
# Gets parsed as "ba"r"
```

#### Lists of key-values

Often a simple Key => Value mapping is not sufficient, and you want a key
mapping to an array of values. This if fully supported and you can define your
struct as:

```go
type Config struct {
  Foo []string
  Bar []float64
  Zoo []int
}
```

And then in your config-file:

```bash
Foo = string number one.
Foo = string number two.
Bar = 1.0
Bar = 2.0
Zoo = 1
Zoo = 2
```

Which, you guessed it, will map to the arrays `Foo{"string number one.",
"string number two"}`, `Bar{1.0,2.0}` and `Zoo{1,2}`.

#### Which types are valid?

At the moment the following types are valid to use when unmarshaling your
config-file:

* String
* Int, Int8, Int16, Int32 and Int64
* Uint, Uint8, Uint16, Uint32 and Uint64
* Float32 and Float64
* Bool

And every one of those as slices, as well. For type definitions and more details
about other types in Golang please refer to [their doc on the
subject](http://golang.org/ref/spec#Types).

#### Using defaults

There are three parts to parsing and defining a config in your application,
given you want to set default values different from those used by Golang.

First you need to define your Config-type. This is done in order to unmarshal
correctly. It is an important step for a type-safe language.  An example
definition looks like:

```go
type Config struct {
  Foo string
}
```


Second you need to create a default-variable of the type you defined in the
previous step.

```go
cfg := &Config{
  Foo: "My default string",
}
```

As you can see our variable `cfg` is a pointer to a Config-type. This pointer is
passed on to ITKconfig which sets the appropriate fields based on your config
file.

Third you use ITKconfig to parse your config-file, validate it and then override
your defaults. This is simply done by:

```go
itkconfig.LoadConfig("filename.conf", cfg)
```

Non-slice keys can only be defined once per config file. Multiple definitions
will produce an error.

## Authors

* Trygve Aaberge ([trygveaa@samfundet.no](mailto:trygveaa@samfundet.no))
* Herman Schistad ([hermansc@samfundet.no](mailto:hermansc@samfundet.no))

Pull-request, your issues and any feedback is greatly appricated.

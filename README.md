# ITKconfig - A small, powerful configuration parser for Golang

Originally started as an internal project at the [Student society in
Trondheim](http://samfundet.no) this package has now been open sourced, as we
believe it is the simplest and best way to manage configuration files for
Go-projects and we'd think others may agree with us as well.

## Features and core-principles

* Makes writing **Key-Value** configuration files easy.
* Allows, in contrast to JSON, comments in your files - just prepend it with a
  `#`.
* Use the same methods as when demarshalling JSON-files, just define your
  configuration struct with your wanted types and let ITKconfig take care of the
  rest.
* Source code is simple and less than 100 lines, this makes it easy to
  understand the flow, but also make changes if you'd like.

## Example configuration file

An example scenario is given where you want to provide a configuration file to
your Web-application. It could look like:

    # Port that the webservice is listening to
    Port = 8000

    # Folder where we find our templates
    TemplatesFolder = "templates"

    # Enable or disable debug mode, giving more output to the user.
    Debug = true

Then, provided that this file is called `myapp.config` we can load it into our
application by the following simple code:

    package main

    import (
      "github.com/itkinside/itkconfig"
      "log"
      "fmt"
    )

    type Config struct {
      Port            int
      TemplatesFolder string
      Debug           bool
    }

    func main() {
      // Some sane defaults for our project.
      config := &Config{
        Port: 80,
        TemplatesFolder: 'temps',
        Debug: false,
      }

      // Override defaults with config-file, if any.
      err := itkconfig.LoadConfig("myapp.config", config)
      if (err != nil) {
        log.Fatal(err)
      }

      // Print our variables, just to show off.
      fmt.Printf("Port: %d\n", config.Port)
      fmt.Printf("Templates: %s\n", config.TemplatesFolder)
      fmt.Printf("Debug: %v\n", config.Debug)
    }

Could it be more simple, and yet so powerful?

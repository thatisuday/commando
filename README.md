# Commando

![logo](/assets/logo.png)

![go-version](https://img.shields.io/github/go-mod/go-version/thatisuday/commando?label=Go%20Version) &nbsp;
![Build](https://github.com/thatisuday/commando/workflows/CI/badge.svg?style=flat-square)

Commando helps you create beautiful CLI applications with ease. It parses [**"getopt(3)"**](http://man7.org/linux/man-pages/man3/getopt.3.html) style command-line arguments, supports sub-command architecture, allows a short-name alias for flags and captures required & optional arguments. The motivation behind creating this library is to provide easy-to-use APIs to create simple command-line tools.

> This library uses [**`clapper`**](https://github.com/thatisuday/clapper) package to parse command-line arguments. Visit the documentation of this package to understand supported **arguments** and **flag** patterns.

![logo](/demo/reactor.gif)

## Documentation
[**pkg.go.dev**](https://pkg.go.dev/github.com/thatisuday/commando)

## Tutorial
[**Medium**](https://medium.com/@thatisuday/building-simple-command-line-cli-applications-in-go-using-commando-8a8e0edbd48a)

## Installation
```
$ go get -u "github.com/thatisuday/commando"
```

## Terminology
Let's imagine we are building a CLI tool `reactor` to generate and manage React front-end projects.

#### Root command
```
$ reactor --version
$ reactor -v
$ reactor --help
$ reactor -h
```

In the example above, the `reactor` alone is the **root-command**. Here, the `--version` and `-v` are flags. The `--version` is a long flag name and `-v` is a short flag name. The `-v` is an alias for `--version`. Similarly `--help` and `-h` are flags.

> Commando adds `--version` and `--help` flags (_along with their aliases_) automatically for the root-command.

#### Sub command
```
$ reactor build --dir ./out/dir
```

In the above example, `build` is a **sub-command**. This sub-command has `--dir` flag, and unlike `--help` flag, this flag takes a user-input value like `./out/dir` we have provided above.

> Commando adds `--help` flag automatically for a sub-command. Also, it registers `help` and `version` sub-commands automatically to display CLI usage and version number respectively.

#### Arguments
```
$ reactor create <name> --dir <dir>
```

In the above example, `create` is a sub-command. The `<name>` placeholder is for the `name` argument value that the user is supposed to provide. For example, `$ reactor create form --dir ./components/form`. Here, `form` is the value for the `name` argument.

> The **root-command** can also be configured to take arguments. For example, `$ reactor <category>` command is valid. If the `category` argument value doesn't match with any registered sub-command, then the root-command is executed, else the sub-command is executed.


## How to configure?
#### Step 1: executable, version and description setup
First, we need to specify the **name of the executable** file using [`commando.SetExecutableName`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#SetExecutableName) function. This value is same as your root-command name. You can use [**`os.Executable`**](https://golang.org/pkg/os/#Executable) function to get the name of the executable. Normally, you just specify the name of your package because when the user install your module using `go get` command, Go creates an executable file with the name same as your package name.

Then we need to set the **version** and the **description** of our CLI application using [`CommandRegistry.SetVersion`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#CommandRegistry.SetVersion) and [`CommandRegistry.SetDescription`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#CommandRegistry.SetDescription) function respectively. This version will be displayed with the `--version` flag (_or `version` sub-command_) and description will be displayed with `--help` flag (_or `help` sub-command_).

```go
import "github.com/thatisuday/commando"

func main() {
    commando.
        SetExecutableName("reactor").
        SetVersion("v1.0.0").
        SetDescription("This CLI tool helps you create and manage React projects.")
}
```

#### Step 2: Register a sub-command
A sub-command is registered using [`commando.Register("<sub-command>")`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Register) function. This function returns the [`*commando.Command`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command) object. If the sub-command is already registered, it returns the registered [`*commando.Command`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command) object.

If `nil` is passed as an argument, then the **root-command** is registered. However, Commando automatically registers the root-command and `commando.Register(nil)` returns the [`*commando.Command`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command) object of the root-command.

> Commando automatically registers the root-command to provide built-in support for `--help` and `--version` flags, hence you do not need to manually register the root-command unless you want to configure it.

The [`*commando.Command`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command) object is a struct and it provides `SetDescription`, `AddArgument`, `AddFlag`, and `SetAction` methods.

- The `SetDescription` method sets the description of the command.
- The `AddArgument` method registers an argument with the command.
- The `AddFlag` method registers a flag with the command.
- The `SetAction` method registers a function that will be executed with **argument values** and **flag values** provided by the user when the root-command or a sub-command is executed by the user.

#### Step 3: Set a description of a sub-command
```go
commando.
  Register("<sub-command>").
  SetDescription("<sub-command-description>")
  SetShortDescription("<sub-command-short-description>")
```

The [`SetDescription`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command.SetDescription) and [`SetShortDescription`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command.SetShortDescription) method takes a `string` argument tp set the long and a short description of the sub-command and returns the [`*commando.Command`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command) object of the same command. These descriptions are printed when user executes the `$ reactor <sub-command> --help` command.

> You can set a description of the **root-command** by passing `nil` as an argument to the [`Register()`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Register) function.


#### Step 4: Add an argument
```go
commando.
  Register("<sub-command>").
  AddArgument("<name>", "<description>", "<default-value>")
```

The [`AddArgument`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command.AddArgument) method registers an argument with the command. The first argument is the name of the argument and the second argument is the description of that argument.

The third argument is a `string` value which is the default-value of the argument. If user doesn't provide the value of this argument, this argument will get the value from the default-value. If the default-value is an empty string (""), then it becomes a **required argument**. If the value of a required argument is not provided by the user, an error message is displayed.

**Ideally, you should register all required arguments before the optional arguments**. Since these are positional values, it is mandatory to do so, else you would get inappropriate results. 

If the argument name ends with `...` suffix, then it is considered as a **variadic argument**. A variadic argument stores all the leftover argument values and concatenate them using command comma (`,`). Hence a command should only contain one variadic argument and it should be registered after all arguments are registered.

> If the argument is already registered, then registration of the argument is skipped without returning an error. You can configure arguments of the **root-command** by passing `nil` as an argument to the [`Register()`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Register) function.

#### Step 5: Add a flag
```go
commando.
  Register("<sub-command>").
  AddFlag("<long-name>,<short-name>", "<description>", <dataType>, <defaultValue>).
```

The [`AddFlag`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command.AddFlag) method registers a flag with the command. The first argument is `string` value containing the long-name and the short-name of the flag separated by a comma. For example, `dir,d` is a valid argument value for `--dir` and `-d` flags. You can skip the short-name registration if you do not need one by providing only long-name value, like `dir`.

The second argument sets the description of the flag. This will be displayed with the usage of the command (_`--help` flag_).

The third argument is the **data-type** of the value that will be provided by the user for this flag. The value of this argument could be either [`commando.Bool`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#pkg-constants), [`commando.Int`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#pkg-constants) or [`commando.String`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#pkg-constants). If the data-type is `commando.Bool`, then the flag does not take any user input (like `--version` flag).

The last argument is the **default-value** of the flag. The value of this argument must be of the data-type provided in the previous argument. If `nil` value is provided, then the flag doesn't have any default-value and it becomes required to be provided by the user, except if the data-type is [`commando.Bool`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#pkg-constants) in which case the default-value is `false` automatically.

If the flag name starts with `no-` prefix, for example `no-clean`, then it is considered as an **inverted flag**. Default value of an inverted flag is `true`. When `--no-clean` flag is provided, value of this flag becomes `false`. This flag is stored without `no-` prefix, like `clean` here, however, `--clean` is not a valid flag.

> If the flag is already registered, then registration of the flag is skipped without returning an error. You should avoid using the same short-name for multiple flags. You can configure flags of the **root-command** by passing `nil` as an argument to the [`Register()`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Register) function.

#### Step 6: Register an action
```go
commando.
  Register("<sub-command>").
  SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
      
  })
```

The [`SetAction`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Command.SetAction) function registers a callback function that will be executed when the root-command or the sub-command is executed by the user. This function is called with the argument values and the flag values provided by the user. If the required argument or a required flag is not provided by the user, this function won't be executed and an error message is shown to the user.

The first argument is a `map` that contains the argument values. The keys of this map are names of the arguments and values are the values of the arguments provided by the user. The values of this map are structs of type [`commando.ArgValue`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#ArgValue) that contains the argument value and other meta-data provided by you during the registration of the argument.

The second argument is a `map` that contains the flag values. The keys of this map are long-names of the flags and values are the values of the flags provided by the user (_or the default-values_). The values of this map are structs of type [`commando.FlagValue`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#FlagValue) that contains the flag value and other meta-data provided by you during the registration of the flag.

The data-type of the `Value` field of the [`commando.ArgValue`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#ArgValue) type is `string`. However, the data-type of the [`commando.FlagValue`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#FlagValue) type is an empty interface `interface{}`. The concrete value of this field can be a `bool`, an `int` or a `string` based on the data-type specified in the flag registration. You should manually extract the concrete value using [**type-assertion**](https://medium.com/rungo/interfaces-in-go-ab1601159b3a#4231). The [`commando.FlagValue`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#FlagValue) also provides `GetBool`, `GetInt` and `GetString` methods to return the flag-value in the correct format. 

#### Step 7: Parse the command-line arguments
```go
commando.Parse(nil)
```

The [`Parse`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#Parse) function parses the command-line arguments provided by the user. Commando executes a sub-command or the root-command based on these values. If the sub-command is not registered or arguments/flags values provided by the user are invalid, then an error message is shown to the user. Else, the action function registered with the command is executed.

```go
commando.Parse([]string{})
```

You can also pass a custom list of command-line arguments. In this case, Commando won't read command-line arguments from the standard-input (_terminal_).

## Example

```go
package main

import (
	"fmt"

	"github.com/thatisuday/commando"
)

func main() {

	// set CLI executable, version and description
	commando.
		SetExecutableName("reactor").
		SetVersion("v1.0.0").
		SetDescription("Reactor is a command-line tool to generate React projects.\nIt helps you create components, write test cases, start a development server and much more.").
		SetEventListener(func(eventName string) {
			//fmt.Println("event-name: ", eventName)
		})

	// configure the root-command
	// $ reactor <category>  --verbose|-V  --version|-v  --help|-h
	commando.
		Register(nil).
		AddArgument("category", "category of the information to look for", ""). // required
		AddFlag("verbose,V", "display log information ", commando.Bool, nil).   // optional
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// print arguments
			for k, v := range args {
				fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			// print flags
			for k, v := range flags {
				fmt.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
			}
		})

	// register `create` sub-command
	// $ reactor create <name> [type]  --dir|-d <dir>  --type|-t [type]  --timeout [timeout]  --verbose|-v  --help|-h
	commando.
		Register("create").
		SetDescription("This command creates a component of a given type and outputs component files in the project directory.").
		SetShortDescription("creates a component").
		AddArgument("name", "name of the component to create", "").                                  // required
		AddArgument("version", "version of the component", "1.0.0").                                 // optional
		AddArgument("files...", "files to remove once component is created", "").                    // variadic, optional
		AddFlag("dir, d", "output directory for the component files", commando.String, nil).         // required
		AddFlag("type, t", "type of the component to create", commando.String, "simple_type").       // optional
		AddFlag("timeout", "operation timeout in seconds", commando.Int, 60).                        // optional
		AddFlag("verbose,v", "display logs while creating the component files", commando.Bool, nil). // optional
		AddFlag("no-clean", "avoid cleanup of the component directory", commando.Bool, nil).         // optional, inverted flag
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// print arguments
			for k, v := range args {
				fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			// print flags
			for k, v := range flags {
				fmt.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
			}
		})

	// register `serve` sub-command
	// $ reactor serve --verbose|-v  --help|-h
	commando.
		Register("serve").
		SetDescription("This command starts the Webpack dev-server on an available port.").
		SetShortDescription("starts a development server").
		AddFlag("verbose,v", "display logs while serving the project", commando.Bool, nil). // optional
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// print arguments
			for k, v := range args {
				fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			// print flags
			for k, v := range flags {
				fmt.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
			}
		})

	// register `build` sub-command
	// $ reactor build  --dir|-d <dir>  --verbose|-v  --help|-h
	commando.
		Register("build").
		SetDescription("This command builds the project with Webpack and outputs the build files in the given directory.").
		SetShortDescription("creates build artifacts").
		AddFlag("dir,d", "output directory of the build files", commando.String, nil).      // required
		AddFlag("verbose,v", "display logs while serving the project", commando.Bool, nil). // optional
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// print arguments
			for k, v := range args {
				fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			// print flags
			for k, v := range flags {
				fmt.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
			}
		})

	// parse command-line arguments from the STDIN
	commando.Parse(nil)
}
```

#### Usage of the root-command

```
$ reactor --help
$ reactor -h
$ reactor help

Reactor is a command-line tool to generate React projects.
It helps you create components, write test cases, start a development server and much more.

Usage:
   reactor <category> {flags}
   reactor <command> {flags}

Commands: 
   build                         creates build artifacts
   create                        creates a component
   help                          displays usage informationn
   serve                         starts a development server
   version                       displays version number

Arguments: 
   category                      category of the information to look for 

Flags: 
   -h, --help                    displays usage information of the application or a command (default: false)
   -V, --verbose                 display log information  (default: false)
   -v, --version                 displays version number (default: false)
```
> You can use [`CommandRegistry.SetEventListener`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#CommandRegistry.SetEventListener) function to add a callback function when **usage information** is displayed.

#### Version of the CLI application
```
$ reactor version
$ reactor --version
$ reactor -v

Version: v1.0.0
```
> You can use [`CommandRegistry.SetEventListener`](https://pkg.go.dev/github.com/thatisuday/commando?tab=doc#CommandRegistry.SetEventListener) function to add a callback function when **version information** is displayed.

#### Usage of the sub-command
```
$ reactor create --help
$ reactor create -h

This command creates a component of a given type and outputs component files in the project directory.

Usage:
   reactor <name> [version] [files] {flags}

Arguments: 
   name                          name of the component to create 
   version                       version of the component (default: 1.0.0)
   files                         files to remove once component is created {variadic}

Flags: 
   --no-clean                    avoid cleanup of the component directory (default: true)
   -d, --dir                     output directory for the component files 
   -h, --help                    displays usage information of the application or a command (default: false)
   --timeout                     operation timeout in seconds (default: 60)
   -t, --type                    type of the component to create (default: simple_type)
   -v, --verbose                 display logs while creating the component files (default: false)
```

#### Executing the root-command
```
$ reactor --verbose
Error: value of the category argument can not be empty.
```

Oops! A user needs to provide the value of the **category** argument since it is a **required** argument.

```
$ reactor service --verbose
$ reactor service -V
arg -> category: service(string)
flag -> version: false(bool)
flag -> help: false(bool)
flag -> verbose: true(bool)
```

#### Executing a sub-command

```
$ reactor create    
Error: value of the name argument can not be empty.
```

Oops! A user needs to provide the value of the **name** argument since it is a **required** argument.
```
$ reactor create my-service 
Error: value of the --dir flag can not be empty.
```

Oops! A user needs to provide the value of the **dir** flag since it is a **required** argument (_no default-value_).

```
$ reactor create my-service -t service --dir ./service/my-service --timeout 10sec 
Error: value of the --timeout flag must be an integer.
```

Oops! Since we have specified that the value of the `--timeout` flag must be an integer, a user needs to provide an integer value.

```
$ reactor create my-service -t service --dir ./service/my-service --timeout 10    
arg -> name: my-service(string)
arg -> version: 1.0.0(string)
arg -> files: (string)
flag -> dir: ./service/my-service(string)
flag -> type: service(string)
flag -> timeout: 10(int)
flag -> verbose: false(bool)
flag -> clean: true(bool)
flag -> help: false(bool)
```

Here, the value of the `version` argument the default value we provided earlier since the user did not provide any value for it. Also, since it is an optional argument, Commando did not print any errors.

```
$ reactor create my-service 2.0.5 file1.txt file2.txt file3.txt -t service --dir=./service/my-service --timeout 10 -v --no-clean
arg -> name: my-service(string)
arg -> version: 2.0.5(string)
arg -> files: file1.txt,file2.txt,file3.txt(string)
flag -> help: false(bool)
flag -> dir: ./service/my-service(string)
flag -> type: service(string)
flag -> timeout: 10(int)
flag -> verbose: true(bool)
flag -> clean: false(bool)
```

## How to create a CLI application?
The example above is a clear demonstration of how a CLI application can be created, however, you can follow this tutorial on [**Medium**](https://medium.com/@thatisuday/building-simple-command-line-cli-applications-in-go-using-commando-8a8e0edbd48a). Here are a few things you should be concerned about.

- Keep your argument names and flag names as simple as possible.
- Try not to override the `--help` or `--version` flags and their short-names.
- Do not configure the **root-command** unless necessary. You do not need to set an **action** function for the root-command. If the action function is missing for the root-command, it won't generate any output or an error.
- Register all optional arguments of a command before the required arguments.
- Do not modify commands after `commando.Parse()` is called.

Your code must be part of the `main` package like we have seen in the previous example. It's better if your work with the [**Go modules**] so that a user can install your application from anywhere on the system.

A user can install the CLI application using `GO111MODULE=on go get "github.com/<username>/<module-name>"` command. Since your code is part of the `main` package, Go creates `<module-name>` binary executable file inside `GOBIN` directory that is supposed to be in the `PATH` of the system.

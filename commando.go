// MIT License

// Copyright (c) 2020 Uday Hiwarale

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package commando helps you create CLI applications with ease.
// It parses "getopt(3)" style command-line arguments, supports sub-command architecture,
// allows a short-name alias for flags and captures required and optional arguments.
package commando

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/thatisuday/clapper"
)

/********************************************/

// data type declarations for the flag values
const (
	// bool data type
	Bool = iota

	// int data type
	Int

	// string data type
	String
)

// root-command name
var rootCommandName = ""

// automatic command and flag descriptions
var (
	helpCommandName         = "help"
	helpCommandDesc         = "This command displays the usage information of this CLI application."
	helpCommandShortDesc    = "displays usage informationn"
	versionCommandName      = "version"
	versionCommandDesc      = "This command displays the version number of this CLI application"
	versionCommandShortDesc = "displays version number"
	helpFlagName            = "help"
	helpFlagShortName       = "h"
	helpFlagDesc            = "displays usage information of the application or a command"
	versionFlagName         = "version"
	versionFlagShortName    = "v"
	versionFlagDesc         = "displays version number"
)

/********************************************/

// CommandRegistry holds the registered command configurations.
// It also stores the version of the CLI application and its description.
type CommandRegistry struct {

	// executable name of the command registry
	Executable string

	// version of the command-line interface
	Version string

	// description of the command-line interface
	Desc string

	// registered command configurations
	Commands map[string]*Command

	// registry to hold `clapper` registry object
	registry clapper.Registry
}

// AddCommand adds a command in the registry.
func (cr *CommandRegistry) addCommand(name string) *Command {

	_name := strings.ReplaceAll(name, " ", "") // remove whitespaces

	// if command is already registered, return the command-config from the registry
	if _c, ok := cr.Commands[_name]; ok {
		return _c
	}

	/*---------------------------*/

	// create a command config
	c := &Command{
		Carg:   cr.registry.Register(_name), // register the command with clapper and store the config-config
		Name:   _name,
		IsRoot: _name == rootCommandName,
		Args:   make(map[string]*Arg),
		Flags:  make(map[string]*Flag),
	}

	/*---------------------------*/

	// register a command config inside the registry
	cr.Commands[_name] = c

	/*---------------------------*/

	// return command config
	return c

}

// SetExecutableName sets the executable name of the registry.
func (cr *CommandRegistry) SetExecutableName(name string) *CommandRegistry {

	if _name := strings.ReplaceAll(name, " ", ""); _name == "" {
		fmt.Printf("Error: executable name must be a non-empty string.\n")
		os.Exit(0)
	} else {
		cr.Executable = _name
	}

	return cr
}

// SetVersion sets the version for the CLI application.
// This version will be printed with the `--version` flag on the root-command.
func (cr *CommandRegistry) SetVersion(version string) *CommandRegistry {

	cr.Version = strings.Trim(version, " ") // trim whitespaces

	return cr
}

// SetDescription sets the description for the CLI application.
// This description will be printed with `--help` flag on the root-command.
func (cr *CommandRegistry) SetDescription(desc string) *CommandRegistry {

	cr.Desc = strings.Trim(desc, " ") // trim whitespaces

	return cr
}

// Register registers a command in the registry and adds `--help` flag automatically.
// If the root-command is registered, it adds the `--version` flags to display the version.
// The "name" argument must be a string. If `nil` is passed, the root-command is registered.
func (cr *CommandRegistry) Register(name interface{}) *Command {

	// name string value
	var nameString string

	// print error if `name` is not a string
	if _, ok := name.(string); name != nil && !ok {
		fmt.Printf("Error: value of the command must be a string.\n")
		os.Exit(0)
	}

	// if `name` is `nil`, it is a root-command
	if name == nil {
		nameString = rootCommandName
	} else {
		nameString = name.(string)
	}

	/*---------------------------*/

	// (replace all whitespaces)
	_name := strings.ReplaceAll(nameString, " ", "") // remove all whitespaces

	// add command to the registry
	c := cr.addCommand(_name)

	/*---------------------------*/

	// add version flag (to print version of the CLI with --version flag)
	if _name == rootCommandName {
		c.AddFlag(fmt.Sprintf("%s,%s", versionFlagName, versionFlagShortName), versionFlagDesc, Bool, nil)
	}

	/*---------------------------*/

	// add help flag (to print usage of the command with --help flag)
	c.AddFlag(fmt.Sprintf("%s,%s", helpFlagName, helpFlagShortName), helpFlagDesc, Bool, nil)

	/*---------------------------*/

	return c
}

// Parse parses the command-line arguments and executes the action function registered with the command.
// If there is an usage-error while parsing the command-line arguments,
// it will display a message in the console without returning an error.
// If osArgs is `nil`, Parse uses arguments received from `os.Args[1:]`.
func (cr *CommandRegistry) Parse(osArgs []string) {

	// get command-line arguments
	_osArgs := osArgs
	if nil == osArgs {
		_osArgs = os.Args[1:]
	}

	// a store for argument values in `ArgValue` format
	argValues := make(map[string]ArgValue)

	// a store for flag values in `ArgFlag` format
	flagValues := make(map[string]FlagValue)

	/*---------------------------*/

	// parse arguments with `clapper` and get the result.
	// `result` is a struct of type `*clapper.Carg`
	result, err := cr.registry.Parse(_osArgs)

	/*---------------------------*/

	// check for errors
	if err != nil {
		switch err.(type) {

		// unknown command
		case clapper.ErrorUnknownCommand:
			errorUnknownCommand := err.(clapper.ErrorUnknownCommand)
			fmt.Printf("Error: %s is not a valid command.\n", errorUnknownCommand.Name)

		// unknown flag
		case clapper.ErrorUnknownFlag:
			errorUnknownFlag := err.(clapper.ErrorUnknownFlag)

			if errorUnknownFlag.IsShort {
				fmt.Printf("Error: -%s is not a valid flag.\n", errorUnknownFlag.Name)
			} else {
				fmt.Printf("Error: --%s is not a valid flag.\n", errorUnknownFlag.Name)
			}
		// unsupported flag
		case clapper.ErrorUnsupportedFlag:
			errorUnsupportedFlag := err.(clapper.ErrorUnsupportedFlag)
			fmt.Printf("Error: %s is not a supported flag.\n", errorUnsupportedFlag.Name)

		// other error
		default:
			fmt.Printf("Error: %s.\n", err)
		}

		// exit process
		os.Exit(0)
	}

	/*---------------------------*/

	// get command configuration from the registry
	c := cr.Commands[result.Cmd]

	/*---------------------------*/

	// if `help` command is provided, display usage of the root-command
	if result.Cmd == helpCommandName {
		cr.PrintHelp(cr.Commands[rootCommandName]) // usage of the root-command
		os.Exit(0)
	}

	// if `--help` or `-h` flag is provided, display usage of the command
	if result.Flags[helpFlagName].Value == "true" {
		cr.PrintHelp(c)
		os.Exit(0)
	}

	// if `version` command or `--version` flag is provided for the root-command, display version number
	if result.Cmd == versionCommandName || (c.IsRoot && result.Flags[versionFlagName].Value == "true") {
		cr.PrintVersion()
		os.Exit(0)
	}

	/*---------------------------*/

	// check if action function is missing
	if c.Action == nil {

		// show error message only for non-root-command
		if !c.IsRoot {
			fmt.Printf("Error: action function for the %s command is not registered.\n", c.Name)
		}

		os.Exit(0)
	}

	/*---------------------------*/

	// for each argument, validate the argument value
	for name, arg := range c.Args {

		// get default-value and user-value of the argument from the `result`
		defaultValue := result.Args[name].DefaultValue
		userValue := result.Args[name].Value

		// get final value
		value := userValue
		if len(userValue) == 0 {
			value = defaultValue
		}

		/*------------*/

		// if argument is required but value is missing, display an error message and exit the process
		if arg.IsRequired && len(value) == 0 {
			fmt.Printf("Error: value of the %s argument can not be empty.\n", name)
			os.Exit(0)
		}

		// save flag display-value inside `argValues`
		argValues[name] = ArgValue{
			Arg:   *arg,
			Value: value,
		}
	}

	/*---------------------------*/

	// for each flag, validate the flag value
	for name, flag := range c.Flags {

		// get default-value and user-value of the flag from the `result`
		defaultValue := result.Flags[name].DefaultValue
		userValue := result.Flags[name].Value

		// get final value
		value := userValue
		if len(userValue) == 0 {
			value = defaultValue
		}

		/*------------*/

		// if flag is required but value is missing, display an error message and exit the process
		if flag.IsRequired && len(value) == 0 {
			fmt.Printf("Error: value of the --%s flag can not be empty.\n", name)
			os.Exit(0)
		}

		/*------------*/

		// convert `value` to an appropriate data type
		var safeValue interface{}
		switch flag.DataType {
		case Bool:
			if value == "true" {
				safeValue = true
			} else {
				safeValue = false
			}
		case Int:
			if _value, err := strconv.ParseInt(value, 10, 64); err == nil {
				safeValue = int(_value)
			} else {
				fmt.Printf("Error: value of the --%s flag must be an integer.\n", name)
				os.Exit(0)
			}
		case String:
			safeValue = value
		}

		/*------------*/

		// save flag display-value inside `argValues`
		flagValues[name] = FlagValue{
			Flag:  *flag,
			Value: safeValue,
		}
	}

	/*---------------------------*/

	// execute action function
	c.Action(argValues, flagValues)

}

// NewCommandRegistry returns a new value of registry and registers the root-command
// with a bare-minimum configuration to enable `--help` and `--version` command.
func NewCommandRegistry() *CommandRegistry {
	registery := &CommandRegistry{
		registry: clapper.NewRegistry(),
		Commands: make(map[string]*Command),
	}

	// add root-command automatically
	registery.Register(nil)

	// add version command automatically
	registery.Register(versionCommandName).SetDescription(versionCommandDesc).SetShortDescription(versionCommandShortDesc)

	// add help command automatically
	registery.Register(helpCommandName).SetDescription(helpCommandDesc).SetShortDescription(helpCommandShortDesc)

	return registery
}

// PrintVersion prints version of the CLI application.
func (cr *CommandRegistry) PrintVersion() {
	// template data
	templateData := struct {
		Version string
	}{
		Version: cr.Version,
	}

	// parse version template
	if tmpl, err := template.New("version").Parse(versionTemplate); err != nil {
		panic(err)
	} else {
		// compile and output template result
		tmpl.Execute(os.Stdout, templateData)
	}
}

// PrintHelp prints the usage of the command or CLI application.
func (cr *CommandRegistry) PrintHelp(c *Command) {

	// get executable name
	exeName := cr.Executable

	// commands (without root-command)
	commands := make(map[string]*Command)
	for k, v := range cr.Commands {
		if !v.IsRoot {
			commands[k] = v
		}
	}

	// template data
	templateData := struct {
		CliDesc       string
		Executable    string
		IsRootCommand bool
		Desc          string
		Args          map[string]*Arg
		Flags         map[string]*Flag
		Commands      map[string]*Command
	}{
		CliDesc:       cr.Desc,
		Executable:    exeName,
		IsRootCommand: c.IsRoot,
		Desc:          c.Desc,
		Args:          c.Args,
		Flags:         c.Flags,
		Commands:      commands,
	}

	// parse help template
	if tmpl, err := template.New("help").Parse(usageTemplate); err != nil {
		panic(err)
	} else {
		// compile and output template result
		tmpl.Execute(os.Stdout, templateData)
	}
}

// DefaultCommandRegistry holds the default registry.
var DefaultCommandRegistry = NewCommandRegistry()

/*---------------------*/

// Command holds the configuration of a command.
type Command struct {

	// command configuration of the `clapper`
	Carg *clapper.Carg

	// name of the command (AKA sub-command)
	Name string

	// description of the command
	Desc string

	// short-description of the command
	ShortDesc string

	// is root-command
	IsRoot bool

	// specific arguments to parse from the command-line arguments
	Args map[string]*Arg

	// flags to parse from the command-line arguments
	Flags map[string]*Flag

	// Action function
	Action func(map[string]ArgValue, map[string]FlagValue)
}

// SetDescription sets the description for a command.
func (c *Command) SetDescription(desc string) *Command {
	c.Desc = strings.Trim(desc, " ") // trim whitespaces

	return c
}

// SetShortDescription sets the short-description for a command.
func (c *Command) SetShortDescription(shortDesc string) *Command {
	c.ShortDesc = strings.Trim(shortDesc, " ") // trim whitespaces

	return c
}

// AddArgument registers an argument for a command.
// When the defaultValue is an empty string, a user needs to provide a value for this argument.
func (c *Command) AddArgument(name string, desc string, defaultValue string) *Command {

	// (replace all whitespaces)
	_name := strings.ReplaceAll(name, " ", "")

	// if the argument is already registered, return
	if _, ok := c.Args[_name]; ok {
		return c
	}

	/*---------------------------*/

	// register the argument with `clapper`
	c.Carg.AddArg(_name, defaultValue)

	/*---------------------------*/

	// create an argument config
	a := &Arg{
		Name:         _name,
		Desc:         strings.Trim(desc, " "), // trim whitespaces
		DefaultValue: defaultValue,
		IsRequired:   defaultValue == "",
	}

	/*---------------------------*/

	// register the argument with the command
	c.Args[_name] = a

	/*---------------------------*/

	return c
}

// AddFlag registers a flag for the command.
// The flagNames argument should contain "long,short" flag names (e.g. "version,v").
// If dataType argument is `commando.Bool` (boolean), then the defaultValue argument is ignored and should be set to `nil`.
// For non-boolean flags, if the defaultValue argument is `nil`, then the flag is required.
func (c *Command) AddFlag(flagNames string, desc string, dataType int, defaultValue interface{}) *Command {

	// (replace all whitespaces)
	_flagNames := strings.ReplaceAll(flagNames, " ", "")

	// split flagNames
	flagNamesList := strings.Split(_flagNames, ",")

	// long flag name
	name := flagNamesList[0]

	// short flag name
	var shortName string
	if len(flagNamesList) > 1 {
		shortName = flagNamesList[1]
	}

	/*---------------------------*/

	// if flag is already registered, return
	if _, ok := c.Flags[name]; ok {
		return c
	}

	/*---------------------------*/

	// format default-value as a string for `clapper`
	var _defaultValue string

	// check if a flag is required
	var _isRequired bool

	// check for correct data type of `defaultValue`
	switch dataType {
	case Bool:
		_defaultValue = "false"
		_isRequired = false
	case Int:
		if defaultValue == nil {
			_isRequired = true
		} else {
			// check if `defaultValue` is a `int`
			if _, ok := defaultValue.(int); !ok {
				fmt.Printf("Error: value of the --%s flag must be an int or nil.\n", name)
				os.Exit(0)
			}

			_defaultValue = fmt.Sprintf("%v", defaultValue.(int))
			_isRequired = false
		}
	case String:
		if defaultValue == nil {
			_isRequired = true
		} else {
			// check if `defaultValue` is a `string`
			if val, ok := defaultValue.(string); !ok {
				fmt.Printf("Error: value of the --%s flag must be a string or nil.\n", name)
				os.Exit(0)
			} else {
				// check for empty string value
				if strings.ReplaceAll(val, " ", "") == "" {
					_isRequired = true
				} else {
					_defaultValue = fmt.Sprintf("%v", defaultValue.(string))
					_isRequired = false
				}
			}
		}
	default:
		fmt.Printf("Error: invalid data type provided for the --%s flag.\n", name)
		os.Exit(0)
	}

	/*---------------------------*/

	// register the flag with `clapper`
	c.Carg.AddFlag(name, shortName, dataType == Bool, _defaultValue)

	/*---------------------------*/

	// create a flag config
	f := &Flag{
		Name:            name,
		ShortName:       shortName,
		Desc:            strings.Trim(desc, " "), // trim whitespaces
		DataType:        dataType,
		DefaultValue:    _defaultValue,
		DefaultValueRaw: defaultValue,
		IsRequired:      _isRequired,
	}

	/*---------------------------*/

	// register the flag with the command
	c.Flags[name] = f

	/*---------------------------*/

	return c
}

// SetAction registers a function to the command configuration that
// will execute after command-line arguments are parsed.
// If an action function is already registered with a command, it won't get registered again.
func (c *Command) SetAction(action func(map[string]ArgValue, map[string]FlagValue)) *Command {

	// set action if not set before
	if c.Action == nil {
		c.Action = action
	}

	return c
}

/*---------------------*/

// Arg defines the configuration of an argument.
type Arg struct {
	Name         string
	Desc         string
	DefaultValue string
	IsRequired   bool
}

// ArgValue represents an argument value to pass as an argument in action function.
type ArgValue struct {
	Arg
	Value string
}

/*---------------------*/

// Flag defines the configuration of a flag.
type Flag struct {
	Name            string
	ShortName       string
	Desc            string
	DataType        int
	DefaultValue    string
	DefaultValueRaw interface{}
	IsRequired      bool
}

// FlagValue represents a flag value to pass as an argument in action function.
// It also provides an easy interface to get value in an appropriate format.
type FlagValue struct {
	Flag
	Value interface{}
}

// GetBool returns `bool` value of a flag.
func (fv FlagValue) GetBool() (bool, error) {
	if fv.DataType == Bool {
		return fv.Value.(bool), nil
	}

	return false, fmt.Errorf("%s flag can not be converted to bool", fv.Name)
}

// GetInt returns `int` value of a flag.
func (fv FlagValue) GetInt() (int, error) {
	if fv.DataType == Int {
		return fv.Value.(int), nil
	}

	return 0, fmt.Errorf("%s flag can not be converted to int", fv.Name)
}

// GetString returns `string` value of a flag.
func (fv FlagValue) GetString() (string, error) {
	if fv.DataType == String {
		return fv.Value.(string), nil
	}

	return "", fmt.Errorf("%s flag can not be converted to string", fv.Name)
}

/*---------------------*/

// SetExecutableName sets the executable name of the `DefaultCommandRegistry` registry.
func SetExecutableName(name string) *CommandRegistry {
	return DefaultCommandRegistry.SetExecutableName(name)
}

// Register registers a command in the `DefaultCommandRegistry` registry.
func Register(name interface{}) *Command {
	return DefaultCommandRegistry.Register(name)
}

// Parse parses the command-line arguments for the `DefaultCommandRegistry` registry.
func Parse(osArgs []string) {
	DefaultCommandRegistry.Parse(osArgs)
}

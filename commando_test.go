package commando

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

/*----------------*/

// executable name should not be empty
func TestEmptyExecutableName(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "tests/empty-exec-name.go")

	// get output
	if output, err := cmd.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: executable name must be a non-empty string.") {
			t.Fail()
		}
	}
}

// a sub-command must have an action function
func TestMissingActionFunction(t *testing.T) {
	// command
	cmdRoot := exec.Command("go", "run", "tests/missing-action-function.go")

	// get output
	if output, err := cmdRoot.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if fmt.Sprintf("%s", output) != "" {
			t.Fail()
		}
	}

	/*---------------*/

	// command
	cmdCreate := exec.Command("go", "run", "tests/missing-action-function.go", "create")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: action function for the create command is not registered.") {
			t.Fail()
		}
	}
}

// default value of a flag must match the data type
func TestInvalidDefaultValue(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/invalid-default-value.go", "create")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the --dir flag must be a string or nil.") {
			t.Fail()
		}
	}
}

/*----------------*/

// unknown command show display an error
func TestUnknownCommand(t *testing.T) {
	// command
	cmdPrint := exec.Command("go", "run", "tests/valid-registry.go", "print")
	cmdPrint.Env = append(os.Environ(), "NO_ROOT=TRUE")

	// get output
	if output, err := cmdPrint.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: print is not a valid command.") {
			t.Fail()
		}
	}
}

// unsupported flag must display an error
func TestUnsupportedFlag(t *testing.T) {

	unsupportedFlags := []string{"---version", "-version"}

	for _, flag := range unsupportedFlags {
		// command
		cmdRoot := exec.Command("go", "run", "tests/valid-registry.go", flag)

		// get output
		if output, err := cmdRoot.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf("Error: %s is not a supported flag.", flag)) {
				t.Fail()
			}
		}
	}
}

// missing argument value of a required argument must display an error
func TestMissingArgument(t *testing.T) {
	// command
	cmdRoot := exec.Command("go", "run", "tests/valid-registry.go")

	// get output
	if output, err := cmdRoot.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the category argument can not be empty.") {
			t.Fail()
		}
	}

	/*----------------*/

	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the name argument can not be empty.") {
			t.Fail()
		}
	}
}

// missing flag value of a required flag must display an error
func TestMissingFlag(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create", "my-service")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the --dir flag can not be empty.") {
			t.Fail()
		}
	}
}

// wrong value of a flag must display an error
func TestInvalidFlagValue(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create", "my-service", "-d", "./services/my-service", "--timeout", "10sec")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the --timeout flag must be an integer.") {
			t.Fail()
		}
	}
}

// test if all values of a root-command are valid
func TestValidRootCommand(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "service")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		values := []string{
			"arg -> category: service(string)",
			"flag -> verbose: false(bool)",
			"flag -> version: false(bool)",
			"flag -> help: false(bool)",
		}

		for _, value := range values {
			if !strings.Contains(fmt.Sprintf("%s", output), value) {
				t.Fail()
			}
		}
	}
}

// test if default value of an argument is correct
func TestDefaultArgValue(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create", "my-service", "-t", "service", "--dir=./service/my-service", "--timeout", "10", "-v")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "arg -> version: 1.0.0(string)") {
			t.Fail()
		}
	}
}

// test if all values of a sub-command are valid
func TestValidSubCommand(t *testing.T) {
	// command
	cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create", "my-service", "1.0.0", "file1.txt", "file2.txt", "-t", "service", "--dir=./service/my-service", "--timeout", "10", "-v", "--no-clean")

	// get output
	if output, err := cmdCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		values := []string{
			"arg -> version: 1.0.0(string)",
			"arg -> name: my-service(string)",
			"arg -> files: file1.txt,file2.txt(string)",
			"flag -> dir: ./service/my-service(string)",
			"flag -> type: service(string)",
			"flag -> timeout: 10(int)",
			"flag -> verbose: true(bool)",
			"flag -> help: false(bool)",
			"flag -> clean: false(bool)",
		}

		for _, value := range values {
			if !strings.Contains(fmt.Sprintf("%s", output), value) {
				t.Fail()
			}
		}
	}
}

// test if version is displayed properly
func TestValidVersion(t *testing.T) {

	versionTriggers := []string{"-v", "--version", "version"}

	for _, versionTrigger := range versionTriggers {
		// command
		cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", versionTrigger)

		// get output
		if output, err := cmdCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), "Version: v1.0.0") {
				t.Fail()
			}
		}
	}
}

// test if usage of the root-command is displayed properly
func TestValidRootCommandUsage(t *testing.T) {

	helpTriggers := []string{"-h", "--help", "help"}

	for _, helpTrigger := range helpTriggers {
		// command
		cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", helpTrigger)

		// get output
		if output, err := cmdCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			values := []string{
				"Reactor is a command-line tool to generate React projects.",
				"It helps you create components, write test cases, start a development server and much more.",

				"Usage:",
				"reactor <category> {flags}",
				"reactor <command> {flags}",

				"Commands: ",
				"build                         creates build artifacts",
				"create                        creates a component",
				"help                          displays usage informationn",
				"serve                         starts a development server",
				"version                       displays version number",

				"Arguments: ",
				"category                      category of the information to look for",

				"Flags: ",
				"-h, --help                    displays usage information of the application or a command (default: false)",
				"-V, --verbose                 display log information (default: false)",
				"-v, --version                 displays version number (default: false)",
			}

			for _, value := range values {
				if !strings.Contains(fmt.Sprintf("%s", output), value) {
					t.Fail()
				}
			}
		}
	}
}

// test if usage of the sub-command is displayed properly
func TestValidSubCommandUsage(t *testing.T) {

	helpTriggers := []string{"-h", "--help"}

	for _, helpTrigger := range helpTriggers {
		// command
		cmdCreate := exec.Command("go", "run", "tests/valid-registry.go", "create", helpTrigger)

		// get output
		if output, err := cmdCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			values := []string{
				"This command creates a component of a given type and outputs component files in the project directory.",

				"Usage:",
				"reactor <name> [version] [files] {flags}",

				"Arguments: ",
				"name                          name of the component to create",
				"version                       version of the component (default: 1.0.0)",
				"files                         files to remove once component is created {variadic}",

				"Flags: ",
				"-d, --dir                     output directory for the component files",
				"-h, --help                    displays usage information of the application or a command",
				"--timeout                     operation timeout in seconds (default: 60)",
				"-t, --type                    type of the component to create (default: simple_type)",
				"-v, --verbose                 display logs while creating the component files (default: false)",
				"--no-clean                    avoid cleanup of the component directory (default: false)",
			}

			for _, value := range values {
				if !strings.Contains(fmt.Sprintf("%s", output), value) {
					t.Fail()
				}
			}
		}
	}
}

// test if version and help events are working properly
func TestEvents(t *testing.T) {

	// test `version` event
	cmdRootVersion := exec.Command("go", "run", "tests/valid-registry.go", "--version")
	cmdRootVersion.Env = append(os.Environ(), "LISTEN_EVENTS=TRUE")

	// get output
	if output, err := cmdRootVersion.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf("event-name: %s", EventVersion)) {
			t.Fail()
		}
	}

	/*--------------*/

	// test `help` event
	cmdRootHelp := exec.Command("go", "run", "tests/valid-registry.go", "--help")
	cmdRootHelp.Env = append(os.Environ(), "LISTEN_EVENTS=TRUE")

	// get output
	if output, err := cmdRootHelp.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf("event-name: %s", EventHelp)) {
			t.Fail()
		}
	}
}

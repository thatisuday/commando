package commando

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

/*----------------*/

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
		//fmt.Printf("Output: %s\n", output)

		if !strings.Contains(fmt.Sprintf("%s", output), "Error: action function for the create command is not registered.") {
			t.Fail()
		}
	}
}

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

func TestUnknownCommand(t *testing.T) {
	// command
	cmdPrint := exec.Command("go", "run", "tests/valid-registry.go", "print")

	// get output
	if output, err := cmdPrint.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		//fmt.Printf("Output: %s\n", output)

		if !strings.Contains(fmt.Sprintf("%s", output), "Error: print is not a valid command.") {
			t.Fail()
		}
	}
}

func TestUnsupportedFlag(t *testing.T) {

	unsupportedFlags := []string{"---version", "-version"}

	for _, flag := range unsupportedFlags {
		// command
		rootCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", flag)

		// get output
		if output, err := rootCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), fmt.Sprintf("Error: %s is not a supported flag.", flag)) {
				t.Fail()
			}
		}
	}
}

func TestMissingArgument(t *testing.T) {
	// command
	rootCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go")

	// get output
	if output, err := rootCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the category argument can not be empty.") {
			t.Fail()
		}
	}

	/*----------------*/

	// command
	createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", "create")

	// get output
	if output, err := createCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the name argument can not be empty.") {
			t.Fail()
		}
	}
}

func TestMissingFlag(t *testing.T) {
	// command
	createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", "create", "my-service")

	// get output
	if output, err := createCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the --dir flag can not be empty.") {
			t.Fail()
		}
	}
}

func TestInvalidFlagValue(t *testing.T) {
	// command
	createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", "create", "my-service", "-d", "./services/my-service", "--timeout", "10sec")

	// get output
	if output, err := createCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		if !strings.Contains(fmt.Sprintf("%s", output), "Error: value of the --timeout flag must be an integer.") {
			t.Fail()
		}
	}
}

func TestValidRootCommand(t *testing.T) {
	// command
	createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", "service")

	// get output
	if output, err := createCreate.Output(); err != nil {
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

func TestValidSubCommand(t *testing.T) {
	// command
	createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", "create", "my-service", "services/my-service", "-t", "service", "--dir=./service/my-service", "--timeout", "10", "-v")

	// get output
	if output, err := createCreate.Output(); err != nil {
		fmt.Println("Error:", err)
	} else {
		values := []string{
			"arg -> alias: services/my-service(string)",
			"arg -> name: my-service(string)",
			"flag -> dir: ./service/my-service(string)",
			"flag -> type: service(string)",
			"flag -> timeout: 10(int)",
			"flag -> verbose: true(bool)",
			"flag -> help: false(bool)",
		}

		for _, value := range values {
			if !strings.Contains(fmt.Sprintf("%s", output), value) {
				t.Fail()
			}
		}
	}
}

func TestValidVersion(t *testing.T) {

	versionTriggers := []string{"-v", "--version", "version"}

	for _, versionTrigger := range versionTriggers {
		// command
		createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", versionTrigger)

		// get output
		if output, err := createCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			if !strings.Contains(fmt.Sprintf("%s", output), "Version: v1.0.0") {
				t.Fail()
			}
		}
	}
}

func TestValidHelp(t *testing.T) {

	helpTriggers := []string{"-h", "--help", "help"}

	for _, helpTrigger := range helpTriggers {
		// command
		createCreate := exec.Command("go", "run", "tests/valid-registry-with-root.go", helpTrigger)

		// get output
		if output, err := createCreate.Output(); err != nil {
			fmt.Println("Error:", err)
		} else {
			values := []string{
				"Reactor is a command-line tool to generate ReactJS projects.",
				"It helps you create components, write test cases, start a development server and much more.",

				"Usage:",
				"reactor <category> [flags]",
				"reactor <command> [flags]",

				"Commands: ",
				"build                         creates build artifacts",
				"create                        creates a React component",
				"help                          displays usage informationn",
				"serve                         starts a development server",
				"version                       displays version number",

				"Arguments: ",
				"category                      category of the information to look for",

				"Flags: ",
				"-h, --help                    displays usage information of the application or a command",
				"-V, --verbose                 display log information",
				"-v, --version                 displays version number",
			}

			for _, value := range values {
				if !strings.Contains(fmt.Sprintf("%s", output), value) {
					t.Fail()
				}
			}
		}
	}
}

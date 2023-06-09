/*
Copyright 2023 Arkemis S.r.l.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// InputOptions stores the user input options
type InputOptions struct {
	ProjectName string
	Local       bool
	Frontend    bool
	Backend     bool
	Console     bool
}

// questions to ask the user
var qs = []*survey.Question{
	{
		Name:     "projectName",
		Prompt:   &survey.Input{Message: "How do you want to call your project?"},
		Validate: survey.Required,
	},
	{
		Name: "local",
		Prompt: &survey.Confirm{
			Message: "Do you want to create the app locally?",
			Default: false,
		},
	},
}

// getInputs retrieves user inputs based on the command and arguments
func getInputs(cmd cobra.Command, args []string) InputOptions {
	var options InputOptions

	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		answers := struct {
			ProjectName string
			Local       bool
			Repos       []string
		}{}

		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		options.ProjectName = answers.ProjectName
		options.Local = answers.Local

		if answers.Local {
			repoQs := &survey.MultiSelect{Message: "Which repo do you want to initialize?", Options: []string{"Frontend", "Backend", "Console"}, Default: []string{"Frontend", "Backend", "Console"}}
			err := survey.AskOne(repoQs, &answers.Repos)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		for _, repo := range answers.Repos {
			switch repo {
			case "Frontend":
				options.Frontend = true
			case "Backend":
				options.Backend = true
			case "Console":
				options.Console = true
			}
		}

	} else {
		if len(args) == 0 {
			fmt.Println("Please provide a project name or use the --interactive flag")
			os.Exit(1)
		}

		options.ProjectName = args[0]

		// get flags and args
		options.Local, _ = cmd.Flags().GetBool("local")
		options.Frontend, _ = cmd.Flags().GetBool("frontend")
		options.Backend, _ = cmd.Flags().GetBool("backend")
		options.Console, _ = cmd.Flags().GetBool("console")

		// if none of console, frontend or backend flags are set, set all to true
		if !cmd.Flags().Changed("console") && !cmd.Flags().Changed("frontend") && !cmd.Flags().Changed("backend") {
			options.Frontend = true
			options.Backend = true
			options.Console = true
		}
	}

	return options
}

// createConsoleApp clones the console repository
func createConsoleApp(projectName string) {
	_, err := exec.Command("git", "clone", "https://github.com/arkemishub/arke-console", fmt.Sprintf("%v/arke-console", projectName)).Output()
	if err != nil {
		fmt.Println("Error cloning console repository", err)
		os.Exit(1)
	}

	fmt.Println("arke-console ✅")
}

// createFrontendApp clones the frontend repository
func createFrontendApp(projectName string) {
	_, err := exec.Command("git", "clone", "https://github.com/arkemishub/frontend-starters", "-n", fmt.Sprintf("%v/temp", projectName)).Output()
	if err != nil {
		fmt.Println("Error cloning frontend repository", err)
		os.Exit(1)
	}

	_, err = exec.Command("git", "-C", fmt.Sprintf("%v/temp", projectName), "sparse-checkout", "init", "--cone").Output()
	if err != nil {
		fmt.Println("Error sparse checking out frontend repository", err)
		os.Exit(1)
	}

	_, err = exec.Command("git", "-C", fmt.Sprintf("%v/temp", projectName), "sparse-checkout", "set", "examples/nextjs-base").Output()
	if err != nil {
		fmt.Println("Error sparse checking out frontend repository", err)
		os.Exit(1)
	}

	_, err = exec.Command("git", "-C", fmt.Sprintf("%v/temp", projectName), "checkout", "main").Output()
	if err != nil {
		fmt.Println("Error pulling from frontend repository", err)
		os.Exit(1)
	}

	err = os.Rename(fmt.Sprintf("%v/temp/examples/nextjs-base", projectName), fmt.Sprintf("%v/frontend", projectName))
	if err != nil {
		fmt.Println("Error moving files from frontend repository", err)
		os.Exit(1)
	}

	err = os.RemoveAll(fmt.Sprintf("%v/temp", projectName))
	if err != nil {
		fmt.Println("Error removing temp folder", err)
		os.Exit(1)
	}

	fmt.Println("frontend ✅")
}

// createAppCmd is the Cobra command for creating a new Arke application
var createAppCmd = &cobra.Command{
	Use:   "create-app [project name]",
	Short: "Creates a new Arke application",
	Long:  `Clones from the Arke template repositories and sets up the application.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		options := getInputs(*cmd, args)

		// if local flag is set, create app locally
		if options.Local {
			// creates a folder in current directory with provided argument as name
			if err := os.Mkdir(options.ProjectName, os.ModePerm); err != nil {
				fmt.Println("Error creating directory", err)
				os.Exit(1)
			}

			shouldCreateAll := options.Frontend && options.Backend && options.Console

			if shouldCreateAll || options.Console {
				createConsoleApp(options.ProjectName)
			}

			if shouldCreateAll || options.Frontend {
				createFrontendApp(options.ProjectName)
			}

		} else {
			// check if ARKEPATH exists within env vars
			path := os.Getenv("ARKEPATH")
			if path == "" {
				fmt.Println("ARKEPATH is not set. Please add it to your environment variables.")
				os.Exit(1)
			}

			composePath := path + "/docker-compose.yml"

			// check if docker-compose.yml exists
			if _, err := os.Stat(composePath); os.IsNotExist(err) {
				fmt.Println("Please run arkectl install before init.")
				os.Exit(1)
			}

			// set PROJECT_ID env var
			e := os.Setenv("PROJECT_ID", options.ProjectName)
			if e != nil {
				fmt.Println("Error setting environment variable", e)
				os.Exit(1)
			}

			// runs docker compose up at ARKEPATH file location
			command := exec.Command("docker", "compose", "-f", composePath, "up")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Run()
			if err != nil {
				fmt.Println("Error running docker compose up", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createAppCmd)
	createAppCmd.Flags().BoolP("interactive", "i", false, "Use interactive version of create-app command")
	createAppCmd.Flags().BoolP("local", "l", false, "Clones arke repositories locally instead of using Docker")
	createAppCmd.Flags().BoolP("frontend", "f", false, "Creates only the frontend app. Works only with --local flag")
	createAppCmd.Flags().BoolP("backend", "b", false, "Clones only the console repository. Works only with --local flag")
	createAppCmd.Flags().BoolP("console", "c", false, "Creates only the backend app. Works only with --local flag")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var repos = map[string]string{"arke-console": "arke-console"}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project name]",
	Short: "Initialise a new Arke application",
	Long:  `Clones from the Arke template repositories and sets up the application.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called", args)

		// if --dev flag is provided then it will clone the dev repos
		if dev, _ := cmd.Flags().GetBool("dev"); dev {
			// creates a folder in current directory with provided argument as name
			if err := os.Mkdir(args[0], os.ModePerm); err != nil {
				fmt.Println("Error creating directory", err)
				os.Exit(1)
			}

			// runs arke.new command to create a new backend project
			command := exec.Command("mix arke.new", fmt.Sprintf("%v-backend", args[0]))
			command.Dir = args[0]
			_, err := command.Output()
			if err != nil {
				_, err := exec.Command("mix archive.install hex arke_new").Output()
				if err != nil {
					fmt.Println("Error installing arke_new package make sure that Elixir is installed on your machine", err)
					os.Exit(1)
				}
				command.Output()
			}

			// clones the template repo into the folder
			for key, repo := range repos {
				if _, err := git.PlainClone(fmt.Sprintf("%v/%v-%v", args[0], args[0], key), false, &git.CloneOptions{
					URL:      fmt.Sprintf("https://github.com/arkemishub/%v", repo),
					Progress: os.Stdout,
				}); err != nil {
					fmt.Println("Error cloning repository", err)
					os.Exit(1)
				}
			}

		} else {
			// if --dev flag is not provided then run docker compose up with PROJECT_ID as environment variable
			e := os.Setenv("PROJECT_ID", args[0])
			if e != nil {
				fmt.Println("Error setting environment variable", e)
				os.Exit(1)
			}

			command := exec.Command("docker", "compose", "up")
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
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("dev", "d", false, "Initialise a new Arke application with dev repos")
}

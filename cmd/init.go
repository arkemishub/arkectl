/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var repos = map[string]string{"console": "console", "frontend": "arke-console-starter", "backend": "arke-backend-starter"}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project name]",
	Short: "Initialise a new Arke application",
	Long:  `Clones from the Arke template repositories and sets up the application.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called", args)

		//creates a folder in current directory with provided argument as name
		if err := os.Mkdir(args[0], os.ModePerm); err != nil {
			fmt.Println("Error creating directory", err)
			os.Exit(1)
		}

		//clones the template repo into the folder
		for key, repo := range repos {
			if _, err := git.PlainClone(fmt.Sprintf("%v/%v-%v", args[0], args[0], key), false, &git.CloneOptions{
				URL:      fmt.Sprintf("https://github.com/arkemishub/%v", repo),
				Progress: os.Stdout,
			}); err != nil {
				fmt.Println("Error cloning repository", err)
				os.Exit(1)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

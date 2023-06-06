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

	"github.com/spf13/cobra"
)

var repos = map[string]string{"console": "console", "frontend": "arke-console-starter", "backend": "arke-backend-starter"}

var initCmd = &cobra.Command{
	Use:   "init [project name]",
	Short: "Initialise a new Arke application",
	Long:  `Clones from the Arke template repositories and sets up the application.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
		e := os.Setenv("PROJECT_ID", args[0])
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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

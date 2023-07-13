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

const COMPOSE_FILE = "https://raw.githubusercontent.com/arkemishub/arkectl/main/docker-compose.yml"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Setup local enviroment for working with arkectl",
	Long:  `Downloads docker-compose.yml and arke.new mix archive and other files needed to run arkectl locally to ARKEPATH directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check if ARKEPATH exists within env vars
		path := os.Getenv("ARKEPATH")
		if path == "" {
			fmt.Println("ARKEPATH is not set. Please add it to your environment variables.")
			os.Exit(1)
		}

		// get docker-compose.yml with curl
		_, err := exec.Command("curl", "-o", path+"/docker-compose.yml", COMPOSE_FILE, "--create-dirs").Output()

		if err != nil {
			fmt.Println("Error downloading docker-compose.yml", err)
			os.Exit(1)
		}

		// install arke_new mix archive
		_, err = exec.Command("mix", "archive.install", "hex", "arke_new", "--force").Output()

		if err != nil {
			fmt.Println("Error installing arke_new mix archive", err)
			os.Exit(1)
		}

		fmt.Println("arkectl has been installed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a lokust test",
	Long: `Delete a lokust test by name

	lokust delete [test-name]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the locust test name")
		}

		return nil
	},
	Run: deleteRun,
}

func deleteRun(cmd *cobra.Command, args []string) {
	testName := args[0]
	err := kclientset.CoreV1().ConfigMaps(config.Namespace).Delete(buildResourceName(testName, "configmap"), nil)
	if err != nil {
		fmt.Printf("Failed deleting %s locustfile configmap. err: %s\n", testName, err)
		os.Exit(1)
	}

	err = ltclientset.LoadtestsV1beta1().LocustTests(config.Namespace).Delete(testName, nil)
	if err != nil {
		fmt.Printf("Failed deleting %s test. err: %s\n", testName, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	setGlobalFlags(deleteCmd.Flags())
}

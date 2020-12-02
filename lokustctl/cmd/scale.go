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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Set a new number of worker nodes",
	Long:  `Set a new number of worker nodes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the locust test name")
		}

		return nil
	},
	Run: scaleRun,
}

func scaleRun(cmd *cobra.Command, args []string) {
	config.Name = args[0]

	lt, err := ltclientset.LoadtestsV1beta1().LocustTests(config.Namespace).Get(config.Name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Failed scaling %s test. err: %s", lt.Name, err)
		os.Exit(1)
	}

	lt.Spec.Replicas = &config.Replicas

	lt, err = ltclientset.LoadtestsV1beta1().LocustTests(config.Namespace).Update(lt)
	if err != nil {
		fmt.Printf("Failed scaling %s test. err: %s", lt.Name, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scaleCmd)

	setGlobalFlags(scaleCmd.Flags())

	scaleCmd.Flags().Int32Var(&config.Replicas, "replicas", 1, "worker nodes")
}

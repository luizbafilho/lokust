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
	"fmt"

	"github.com/k0kubun/pp"
	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates new locust test",
	Long: `Creates a new distributed locust test.

	lokust create --name [test-name] -f locustfile.py
	`,
	Run: func(cmd *cobra.Command, args []string) {
		pp.Println(config)
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().StringVar(config.Name, "name", "", "Test name")
	createCmd.Flags().StringVar(&config.Name, "name", "", "Test name")
}

func buildLocustTest(name string, namespace string) (loadtestsv1beta1.LocustTest, error) {
	replicas := int32(1)
	test := loadtestsv1beta1.LocustTest{
		TypeMeta: metav1.TypeMeta{APIVersion: loadtestsv1beta1.GroupVersion.String(), Kind: "LocustTest"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: loadtestsv1beta1.LocustTestSpec{
			Replicas: &replicas, // won't be nil because defaulting
			// Resources: nil,
		},
	}

	return test, nil
}

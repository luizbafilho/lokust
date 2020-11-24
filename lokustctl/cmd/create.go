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
	"os"

	"github.com/k0kubun/pp"
	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates new locust test",
	Long: `Creates a new distributed locust test.

	lokust create --name [test-name] -f locustfile.py
	`,
	Run: createRun,
}

func createRun(cmd *cobra.Command, args []string) {
	pp.Println(config)
	lt, err := ltclientset.LoadtestsV1beta1().LocustTests(config.Namespace).Create(buildLocustTest(config))
	if err != nil {
		fmt.Printf("Failed creating %s test. err: %s", lt.Name, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	viper.BindPFlags(createCmd.Flags())

	createCmd.Flags().StringVar(&config.Name, "name", "", "Test name")
	createCmd.MarkFlagRequired("name")

	createCmd.Flags().Int32Var(&config.Replicas, "replicas", 1, "Worker nodes")
}

func buildLocustTest(config Config) *loadtestsv1beta1.LocustTest {

	test := loadtestsv1beta1.LocustTest{
		TypeMeta: metav1.TypeMeta{APIVersion: loadtestsv1beta1.SchemeGroupVersion.String(), Kind: "LocustTest"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
		Spec: loadtestsv1beta1.LocustTestSpec{
			Replicas: &config.Replicas, // won't be nil because defaulting
			Resources: loadtestsv1beta1.LocustTestResources{
				Master: corev1.ResourceRequirements{
					Limits:   ConvertToCoreV1ResourceList(config.Resources.Master.Limits),
					Requests: ConvertToCoreV1ResourceList(config.Resources.Master.Requests),
				},
				Workers: corev1.ResourceRequirements{
					Limits:   ConvertToCoreV1ResourceList(config.Resources.Workers.Limits),
					Requests: ConvertToCoreV1ResourceList(config.Resources.Workers.Requests),
				},
			},
		},
	}

	return &test
}

func ConvertToCoreV1ResourceList(resourceList map[string]string) corev1.ResourceList {
	capacity := make(corev1.ResourceList)

	if len(resourceList) > 0 {
		for k, v := range resourceList {
			quantity, err := resource.ParseQuantity(v)
			if err != nil {
				fmt.Printf("Failed converting %s into resource.Quantity. err: %s", k, err)
				os.Exit(1)
			}
			capacity[corev1.ResourceName(k)] = quantity
		}
	}

	return capacity
}

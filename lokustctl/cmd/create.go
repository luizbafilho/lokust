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
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
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

	cm, err := kclientset.CoreV1().ConfigMaps(config.Namespace).Create(buildLocustfile(config))
	if err != nil {
		fmt.Printf("Failed creating %s locustfile configmap. err: %s", cm.Name, err)
		os.Exit(1)
	}

	lt, err := ltclientset.LoadtestsV1beta1().LocustTests(config.Namespace).Create(buildLocustTest(config, cm))
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
	createCmd.Flags().StringVarP(&config.Locustfile, "locustfile", "f", "", "Python module file to import, e.g. 'locustfile.py'")
	createCmd.MarkFlagRequired("locustfile")

	createCmd.Flags().Int32Var(&config.Replicas, "replicas", 1, "Worker nodes")
}

func buildLocustTest(config Config, cm *corev1.ConfigMap) *loadtestsv1beta1.LocustTest {
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
			ConfigmapName: cm.Name,
		},
	}

	return &test
}

func buildLocustfile(config Config) *corev1.ConfigMap {
	info, err := os.Stat(config.Locustfile)
	if err != nil {
		fmt.Printf("Error reading locustfile: %s. err: %s\n", config.Locustfile, err)
		os.Exit(1)
	}

	var cmData map[string]string

	if info.IsDir() {
		fmt.Printf("lokust does not support using directories in the moment. \nPlease open an issue on GitHub if you want to see that feature.")
		os.Exit(1)
	} else {
		cmData, err = buildConfigmapFromFile(config.Locustfile)
		if err != nil {
			fmt.Printf("Error reading locustfile. err: %s\n", err)
			os.Exit(1)
		}
	}

	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildResourceName(config.Name, "configmap"),
			Namespace: config.Namespace,
		},
		Data: cmData,
	}

	return &configMap
}

func buildConfigmapFromFile(filename string) (map[string]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	configMapData := make(map[string]string, 0)
	configMapData["locustfile.py"] = string(file)

	return configMapData, nil
}

func buildResourceName(testName string, resourceType ...string) string {
	name := fmt.Sprintf("lokust-%s", testName)

	if len(resourceType) > 1 {
		name += "-" + resourceType[0]
	}

	return name
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

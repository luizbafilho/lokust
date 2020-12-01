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
	"path/filepath"

	loadtestsclientset "github.com/luizbafilho/lokust/clientset/versioned"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/spf13/viper"
)

var (
	config      Config = Config{}
	cfgFile     string
	ltclientset *loadtestsclientset.Clientset
	kclientset  *kubernetes.Clientset
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lokustctl",
	Short: "lokustctl controls lokust operator",
	Long: `lokustctl helps you manager your distributed tests on kubernetes by
	creating all necessary resources so you don't have to manage multiple yaml files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	ltclientset, err = loadtestsclientset.NewForConfig(config)
	if err != nil {
		fmt.Printf("Unable to create loadtest kubernetes client: %s", err)
		os.Exit(1)
	}
	// create the clientset
	kclientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Unable to create kubernetes client: %s", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "lokustctl" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("lokust")
	}

	// viper.BindPFlags(rootCmd.Flags())
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		os.Exit(1)
	}
}

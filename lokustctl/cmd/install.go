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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/commands/build"
	"sigs.k8s.io/kustomize/pkg/fs"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Output Kubernetes resources to install Lokust operator",
	Long:  `Output Kubernetes resources to install Lokust operator`,
	Run:   installRun,
}

func installRun(cmd *cobra.Command, args []string) {
	dir, err := ioutil.TempDir(".", "lokust-install-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = copyDir("/config", dir)
	if err != nil {
		panic(err)
	}

	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	f := k8sdeps.NewFactory()
	o := build.NewOptions(dir+"/default", "")
	err = o.RunBuild(stream.Out, fs.MakeRealFS(), f.ResmapF, f.TransformerF)
	if err != nil {
		panic(err)
	}
}

func copyDir(source, destination string) error {
	var err error = pkger.Walk(source, func(path string, info os.FileInfo, err error) error {
		path = strings.Split(path, ":")[1]
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			return copyFile(filepath.Join(source, relPath), filepath.Join(destination, relPath))
		}
	})
	return err
}

// File copies a single file from src to dst
func copyFile(src, dst string) error {
	var err error
	var srcfd pkging.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = pkger.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = pkger.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func init() {
	rootCmd.AddCommand(installCmd)
}

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

	"github.com/luizbafilho/lokust/lokustctl/kustomize/kustfile"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/konfig"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
)

var (
	kustomizeNamespace  string
	kustomizeNamePrefix string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Output Kubernetes resources to install Lokust operator",
	Long:  `Output Kubernetes resources to install Lokust operator`,
	Run:   installRun,
}

func installRun(cmd *cobra.Command, args []string) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir, err := ioutil.TempDir(".", "lokust-install-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	defaultDir := dir + "/default"

	err = copyDir("/config", dir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(defaultDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := runSetNamepaceAndNameprefix(); err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(currDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := runKustomizeBuild(defaultDir); err != nil {
		log.Fatal(err)
	}
}

func runKustomizeBuild(dir string) error {
	opts := &krusty.Options{
		// When true, this means sort the resources before emitting them,
		// per a particular sort order.  When false, don't do the sort,
		// and instead respect the depth-first resource input order as
		// specified by the kustomization files in the input tree.
		// TODO: get this from shell (arg, flag or env).
		DoLegacyResourceSort: true,
		// In the kubectl context, avoid security issues.
		LoadRestrictions: types.LoadRestrictionsRootOnly,
		// In the kubectl context, avoid security issues.
		PluginConfig: konfig.DisabledPluginConfig(),
	}
	m, err := krusty.MakeKustomizer(filesys.MakeFsOnDisk(), opts).Run(dir)
	if err != nil {
		return err
	}
	res, err := m.AsYaml()
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(res)
	return err
}

func runSetNamespace(namespace string) error {
	mf, err := kustfile.NewKustomizationFile(filesys.MakeFsOnDisk())
	if err != nil {
		return err
	}
	m, err := mf.Read()
	if err != nil {
		return err
	}
	m.Namespace = namespace
	return mf.Write(m)
}

func runSetNamepaceAndNameprefix() error {
	mf, err := kustfile.NewKustomizationFile(filesys.MakeFsOnDisk())
	if err != nil {
		return err
	}
	m, err := mf.Read()
	if err != nil {
		return err
	}
	m.NamePrefix = kustomizeNamePrefix
	m.Namespace = config.Namespace
	return mf.Write(m)
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
	installCmd.Flags().StringVar(&kustomizeNamePrefix, "name-prefix", "lokust-", "kubernetes resources name prefix override (default lokust-)")
}

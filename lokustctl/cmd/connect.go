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
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
	"github.com/luizbafilho/lokust/common"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwardAPodRequest struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Pod is the selected pod for this port forwarding
	Pod v1.Pod
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int
	// PodPort is the target port for the pod
	PodPort int
	// Steams configures where to write or read input from
	Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connects to a running instace of a locust test",
	Long: `creates a local proxy binding a local port to the locust master instance running on kubernetes.

    lokustctl connect [test-name]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the locust test name")
		}

		return nil
	},
	Run: connectRun,
}

func connectRun(cmd *cobra.Command, args []string) {
	config.Name = args[0]

	podName, err := getMasterPodName(config.Name)
	if err != nil {
		fmt.Printf("Failed getting master pod name. err: %s", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})
	stream := genericclioptions.IOStreams{
		In: os.Stdin,
		// enable for debugging
		// Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Bye...")
		close(stopCh)
		wg.Done()
	}()

	kubeconfig, err := getKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := PortForwardAPod(PortForwardAPodRequest{
			RestConfig: kubeconfig,
			Pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: config.Namespace,
				},
			},
			LocalPort: config.ConnectPort,
			PodPort:   8089,
			Streams:   stream,
			StopCh:    stopCh,
			ReadyCh:   readyCh,
		})
		if err != nil {
			panic(err)
		}
	}()

	<-readyCh
	fmt.Printf("Proxy created!\nYou can now connect to the Locust instance %s test at http://localhost:%d\n", config.Name, config.ConnectPort)

	wg.Wait()
}

func getMasterPodName(testName string) (string, error) {
	test := loadtestsv1beta1.LocustTest{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
		},
	}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(common.MakeLabels(test, "master")).String(),
	}
	pods, err := kclientset.CoreV1().Pods(config.Namespace).List(listOptions)
	if err != nil {
		return "", fmt.Errorf("error fetching %s test pods list. err: %s", testName, err)
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("No pods found for %s test", testName)
	}

	return pods.Items[0].Name, nil
}

func PortForwardAPod(req PortForwardAPodRequest) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(req.RestConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().IntVarP(&config.ConnectPort, "port", "p", 8089, "Connect local port")
}

package k8s

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sync"
)

// https://refactoring.guru/design-patterns/singleton/go/example

var lock = &sync.Mutex{}

var clientset *kubernetes.Clientset

func GetClientset() *kubernetes.Clientset {
	if clientset == nil {
		lock.Lock()
		defer lock.Unlock()
		if clientset == nil {
			var kubeconfig *string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()

			// use the current context in kubeconfig
			config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				panic(err.Error())
			}

			// create the clientset
			initClientset, err := kubernetes.NewForConfig(config)
			clientset = initClientset

			if err != nil {
				panic(err.Error())
			}
		}
	}

	return clientset
}

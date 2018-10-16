package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/brandon-height/ploy/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	k, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)

	}

	clientset, err := kubernetes.NewForConfig(k)
	if err != nil {
		log.Fatal(err)

	}

	// Setup our config
	c := config.NewConfig(clientset)

	// Establish routes for our router
	c.Routes()

	// Listen and serve on 8000
	log.Fatal(http.ListenAndServe(":8000", c.Router))
}

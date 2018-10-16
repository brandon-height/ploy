package config

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Config defines a client that can be used to connect up to kubernetes client API.
type Config struct {
	Client kubernetes.Interface
	Router *mux.Router
}

// NewConfig returns a new instance of `Config`.
func NewConfig(client kubernetes.Interface) *Config {
	return &Config{
		Client: client,
		Router: mux.NewRouter(),
	}
}

// Routes establishes all the routes for the given `Config` object.
func (c *Config) Routes() {
	// /api/v1 endpoints
	c.Router.HandleFunc("/api/v1/namespace", c.GetNamespaces()).Methods("GET")
	c.Router.HandleFunc("/api/v1/namespace/{namespace}", c.GetNamespace()).Methods("GET")
	//c.Router.HandleFunc("/api/v1/namespace/{namespace}", c.CreateNamespace()).Methods("POST")
	//c.Router.HandleFunc("/api/v1/namespace/{namespace}", c.DeleteNamespace()).Methods("DELETE")
}

// GetNamespaces retrieves all v1.Namespace objects within the current kubernetes context.
func (c *Config) GetNamespaces() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ns, err := c.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ns)
	}
}

// GetNamespace retrieves a specific v1.Namespace object within the current kubernetes context.
func (c *Config) GetNamespace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ns, err := c.Client.CoreV1().Namespaces().Get(vars["namespace"], metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ns)
	}
}

// CreateNamespace creates a new v1.Namespace object into the current kubernetes context.
/*func (c *Config) CreateNamespace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ns, err := c.Client.CoreV1().Namespaces().Get(vars["namespace"], metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ns)
	}
}*/

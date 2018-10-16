package config

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_NewConfig(t *testing.T) {
	c := NewConfig(fake.NewSimpleClientset())
	assert.NotNil(t, c.Client)
	assert.NotNil(t, c.Router)
	assert.Implements(t, (*kubernetes.Interface)(nil), c.Client)

}

// TODO: Need to determine how best to test this.
func Test_Routes(t *testing.T) {
	c := NewConfig(fake.NewSimpleClientset())
	c.Routes()
}

func Test_GetNamespaces(t *testing.T) {
	c := NewConfig(fake.NewSimpleClientset())
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v1/namespace", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.GetNamespaces())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_GetNamespace(t *testing.T) {
	tt := []struct {
		routeVariable string
		expected      int
	}{
		// Kubernetes always responds back with 200 OK even
		// if the namespace doesn't exist.
		{"default", http.StatusOK},
		{"nope", http.StatusOK},
	}
	for _, tc := range tt {
		c := NewConfig(fake.NewSimpleClientset())
		path := fmt.Sprintf("/api/v1/namespace/%s", tc.routeVariable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/namespace/{namespace}", c.GetNamespace())
		router.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, tc.expected)
	}
}

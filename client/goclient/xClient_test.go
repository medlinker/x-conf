package goclient

import (
	"testing"

	"github.com/coreos/etcd/client"
)

func TestOpt(t *testing.T) {

	resp, err := Set("/test2/app", "app2_test", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)

	resp, err = Get("/test2", &client.GetOptions{Recursive: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range resp.Node.Nodes {
		t.Log(node.Key + ": " + node.Value)
	}
}

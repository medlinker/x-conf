package goclient

import (
	"fmt"
	"testing"
)

func TestOpt(t *testing.T) {
	/*
		resp, err := Set("/test2/app", "app2_test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp)

		resp, err = Get("/prjs", &client.GetOptions{Recursive: true})
		if err != nil {
			t.Fatal(err)
		}
		for _, node := range resp.Node.Nodes {
			t.Log(node.Key + ": " + node.Value)
		}
	*/
	SetInit("instanceName", "go-test")
	SetInit("prjName", "app1")
	SetInit("env", "prod")
	SetInit("etcd_clinet_urls", "http://127.0.0.1:2379")
	Config()

	if GetLM("redis.host") == "" {
		t.Fatal("Get key from local memory error!")
	}

	Watching(func() {
		fmt.Println(entry)
	})

	WatchingShare(func() {
		fmt.Println(entry)
	})

	ch := make(chan int, 1)
	<-ch
}

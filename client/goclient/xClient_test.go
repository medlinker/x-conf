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

	if GetLM("redis.host") != "" {
		t.Fatal("Get key from local memory error!")
	}

	Watching(func() {
		fmt.Println(GetLM("redis.host"))
	})

	// ch := make(chan int, 1)
	// <-ch
}

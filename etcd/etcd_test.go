package etcd

import (
	"fmt"
	"github.com/kevinvandervlist/consul-etcd-bootstrapper/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cb(t *testing.T) func(err error) {
	return func(err error) {
		t.Fatalf("An error occured: %s\n", err)
	}
}

func TestSizeParameterUnmarshaller(t *testing.T) {
	expectedValue := "3"
	expectedKey := "/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306/_config/size"
	client := CreateEtcdClient("2655b94e64485fc6f7c3d8dae4820306")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306/_config/size","value":"3","modifiedIndex":540153005,"createdIndex":540153005}}`)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	util.AssertNoErrorCallback(err, cb(t))

	response, err := client.unmarshalResponse(res)
	util.AssertNoErrorCallback(err, cb(t))

	util.AssertEquals(response.Node.Key, expectedKey, t)
	util.AssertEquals(response.Node.Value, expectedValue, t)
}

func TestSingleAnnouncedNode(t *testing.T) {
	expectedValue := "10.20.30.40"
	var expectedTTL uint64 = 8
	client := CreateEtcdClient("2655b94e64485fc6f7c3d8dae4820306")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306","dir":true,"nodes":[{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306/abc","value":"10.20.30.40","expiration":"2015-05-17T22:51:11.251063029Z","ttl":8,"modifiedIndex":542751716,"createdIndex":542751716}],"modifiedIndex":540153002,"createdIndex":540153002}}`)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	util.AssertNoErrorCallback(err, cb(t))

	response, err := client.unmarshalResponse(res)
	util.AssertNoErrorCallback(err, cb(t))

	util.AssertEquals(response.Node.Dir, true, t)
	util.AssertEquals(response.Node.Nodes[0].Value, expectedValue, t)
	util.AssertEquals(response.Node.Nodes[0].TTL, expectedTTL, t)
}

func TestListMultipleAnnouncedNodes(t *testing.T) {
	expectedNodeCount := 2
	client := CreateEtcdClient("2655b94e64485fc6f7c3d8dae4820306")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306","dir":true,"nodes":[{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306/192.168.2.15","value":"192.168.2.15","expiration":"2015-05-18T00:33:14.740272675Z","ttl":44,"modifiedIndex":542916902,"createdIndex":542916902},{"key":"/_etcd/registry/2655b94e64485fc6f7c3d8dae4820306/abc","value":"1","expiration":"2015-05-18T00:32:39.301296169Z","ttl":9,"modifiedIndex":542917283,"createdIndex":542917283}],"modifiedIndex":540153002,"createdIndex":540153002}}`)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	util.AssertNoErrorCallback(err, cb(t))

	response, err := client.unmarshalResponse(res)
	util.AssertNoErrorCallback(err, cb(t))

	util.AssertEquals(response.Node.Dir, true, t)
	util.AssertEquals(len(response.Node.Nodes), expectedNodeCount, t)
}

func TestRequestSizeParameterUrl(t *testing.T) {
	expected := "http://discovery.etcd.io/2655b94e64485fc6f7c3d8dae4820306/_config/size"
	client := CreateEtcdClient("2655b94e64485fc6f7c3d8dae4820306")
	actual := client.buildClusterSizeURL()
	util.AssertEquals(expected, actual, t)
}

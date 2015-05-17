package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type EtcdValueNode struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	Expiration    string `json:"expiration"`
	TTL           uint64 `json:"ttl"`
	ModifiedIndex int    `json:"modifiedIndex"`
	CreatedIndex  int    `json:"createdIndex"`
}

type EtcdNode struct {
	Key           string          `json:"key"`
	Value         string          `json:"value"`
	Dir           bool            `json:"dir"`
	Nodes         []EtcdValueNode `json:"nodes"`
	ModifiedIndex int             `json:"modifiedIndex"`
	CreatedIndex  int             `json:"createdIndex"`
}

type EtcdResponse struct {
	Action string   `json:"action"`
	Node   EtcdNode `json:"node"`
}

type EtcdClient struct {
	BaseUrl string
	token   string
}

func CreateEtcdClient(token string) *EtcdClient {
	return &EtcdClient{"http://discovery.etcd.io/", token}
}

func (e *EtcdClient) unmarshalResponse(res *http.Response) (*EtcdResponse, error) {
	_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	parsed := &EtcdResponse{}
	err = json.Unmarshal(_body, parsed)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (e *EtcdClient) buildClusterSizeURL() string {
	return e.BaseUrl + e.token + "/_config/size"
}

func (e *EtcdClient) buildBaseClusterURL() string {
	return e.BaseUrl + e.token
}

func (e *EtcdClient) buildAnnounceURL(name string) string {
	return e.buildBaseClusterURL() + "/" + name
}

func (e *EtcdClient) buildPayload(value string, ttl uint64) url.Values {
	v := url.Values{}

	if value != "" {
		v.Set("value", value)
	}

	if ttl > 0 {
		v.Set("ttl", fmt.Sprintf("%v", ttl))
	}

	return v
}

func (e *EtcdClient) isValidStatusCode(actual int, valid []int) error {
	for _, v := range valid {
		if actual == v {
			return nil
		}
	}
	message := fmt.Sprintf("Invalid statuscode received. Got %d, expected %v\n", actual, valid)
	return errors.New(message)
}

func (e *EtcdClient) GetClusterSize() (int, error) {
	resp, err := http.Get(e.buildClusterSizeURL())
	valid := []int{http.StatusOK}

	if err != nil {
		return 0, err
	}

	err = e.isValidStatusCode(resp.StatusCode, valid)
	if err != nil {
		return 0, err
	}

	response, err := e.unmarshalResponse(resp)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(response.Node.Value)
}

func (e *EtcdClient) GetCurrentCluster() ([]string, error) {
	resp, err := http.Get(e.buildBaseClusterURL())
	valid := []int{http.StatusOK}
	empty := []string{}

	if err != nil {
		return empty, err
	}

	err = e.isValidStatusCode(resp.StatusCode, valid)
	if err != nil {
		return empty, err
	}

	response, err := e.unmarshalResponse(resp)
	if err != nil {
		return empty, err
	}

	return e.MapCurrentCluster(response)
}

func (e *EtcdClient) MapCurrentCluster(response *EtcdResponse) ([]string, error) {
	ret := []string{}
	if response.Node.Dir == false {
		return ret, errors.New("Expected a directory that can be listed.")
	}
	ret = make([]string, len(response.Node.Nodes), len(response.Node.Nodes))
	for i, n := range response.Node.Nodes {
		log.Printf("Discovered node %v with value %v and TTL %v.\n", n.Key, n.Value, n.TTL)
		ret[i] = n.Value
	}
	return ret, nil
}

func (e *EtcdClient) AnnounceNode(ip string, ttl uint64) (string, error) {
	payload := e.buildPayload(ip, ttl).Encode()
	valid := []int{http.StatusOK, http.StatusCreated}
	url := e.buildAnnounceURL(ip)

	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader(payload))
	request.ContentLength = int64(len(payload))
	request.Header.Set("Content-Type",
		"application/x-www-form-urlencoded; param=value")
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	err = e.isValidStatusCode(resp.StatusCode, valid)
	if err != nil {
		return "", err
	}

	response, err := e.unmarshalResponse(resp)
	if err != nil {
		return "", err
	}

	return response.Node.Value, nil
}

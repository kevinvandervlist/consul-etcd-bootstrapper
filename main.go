package main

import (
	"flag"
	"fmt"
	"github.com/kevinvandervlist/consul-etcd-bootstrapper/consul"
	"github.com/kevinvandervlist/consul-etcd-bootstrapper/etcd"
	"log"
	"os"
	"time"
)

func main() {
	ret := start()
	log.Printf("Exitting...")
	os.Exit(ret)
}

func start() int {
	token_default := "https://discovery.etcd.io/<token>"

	nodeChannel := make(chan []string)
	clusterComplete := make(chan bool)
	announceChannel := make(chan bool)

	token := flag.String("token", token_default, "The ETCD token to use.")
	ip := flag.String("ip", "$public_ipv4", "The public ipv4 address of this node.")
	consulBin := flag.String("consulBin", "/bin/consul", "Location of the consul binary")
	consulArgs := flag.String("consulArgs", "\"-data-dir /data -ui-dir /ui -config-dir /config -client 0.0.0.0\"", "Provide other arguments then shown defaults")

	flag.Parse()

	log.Printf("token: %s\n", *token)
	log.Printf("ip: %s\n", *ip)
	log.Printf("consulBin: %s\n", *consulBin)
	log.Printf("consulArgs: %s\n", *consulArgs)
	log.Printf("Note: Arguments agent -server -bind -join are provided by consul-etcd-bootstrapper.")

	if token_default == *token {
		log.Printf("An ETCD discovery token is mandatory. Please request one at https://discovery.etcd.io/new?size=n")
		log.Printf("where n is your desired cluster size.")
		return 1
	}

	if "$public_ipv4" == *ip {
		log.Printf("An ip address to bind on is necessary. Please provide the concrete $public_ipv4 or $private_ipv4 value")
		log.Printf("of this container instance.")
		return 1
	}


	etcd := etcd.CreateEtcdClient(*token)

	clusterCount, err := etcd.GetClusterSize()
	if err != nil {
		log.Printf("An error occurred while retrieving the expected cluster size: %v\n", err)
		return 3
	}
	log.Printf("Expecting cluster size %v\n", clusterCount)

	// Schedule a periodic announcer
	go announcer(announceChannel, 300, 300, *ip, etcd)
	// Schedule a poller for the cluster members
	go pollCluster(nodeChannel, clusterComplete, 30, clusterCount, etcd)

	// If we get something from the channel we have discovered a cluster.
	nodes := <-nodeChannel
	announceChannel <- true
	log.Printf("Initiating join cluster with the following nodes:\n")
	for _, n := range nodes {
		log.Printf("\t- %s\n", n)
	}

	log.Printf("Bootstrapping consul...")

	wrapper := consul.CreateConsulWrapper()
	wrapper.AddAdditionalArguments(*consulArgs)
	return wrapper.Run(*consulBin, *ip, clusterCount, nodes)
}

func announcer(complete chan bool, interval time.Duration, ttl uint64, ip string, etcd *etcd.EtcdClient) {
	log.Printf("Announcing this node every %d seconds", interval)
	log.Printf("Announcing this node as %v\n", ip)

	announce := func() {
		_, err := etcd.AnnounceNode(ip, ttl)
		if err != nil {
			m := fmt.Sprintf("An error occurred while announcing this node: %v\n", err)
			log.Printf(m)
			// If unable to announce it makes no sense to keep running.
			panic(m)
		}
	}

	t := time.NewTicker(time.Second * interval)
	ticker := t.C
	announce()
	for {
		select {
		case <-complete:
			t.Stop()
			return
		case <-ticker:
			announce()
		}
	}
}

func pollCluster(channel chan []string, complete chan bool, interval time.Duration, clusterSize int, etcd *etcd.EtcdClient) {
	log.Printf("Polling the cluster every %d seconds", interval)

	poll := func() {
		nodes, err := etcd.GetCurrentCluster()
		if err != nil {
			log.Printf("Unable to retrieve a list of nodes in the cluster.")
		}
		log.Printf("Got nodes: %v\n", nodes)
		if len(nodes) >= clusterSize {
			log.Printf("All cluster members detected\n")
			channel <- nodes
		}
	}

	t := time.NewTicker(time.Second * interval)
	ticker := t.C
	poll()

	for {
		select {
		case <-complete:
			t.Stop()
			return
		case <-ticker:
			poll()
		}
	}
}

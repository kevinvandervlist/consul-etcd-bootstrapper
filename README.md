# consul-etcd-bootstrapper
An application to simplify and fully automate the setup of a containerized consul cluster

It allows service discovery via the public etcd cluster of CoreOS

# Warning
It is still a proof of concept.

# Usage
First, retrieve a token, providing the requested cluster size:
https://discovery.etcd.io/new?size=3

Then, instantiate a bunch of nodes as can be seen below.
```
docker run --name consul-master-01 --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token <mytoken>
docker run --name consul-master-01 --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token <mytoken>
docker run --name consul-master-01 --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token <mytoken> 
```
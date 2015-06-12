# consul-etcd-bootstrapper
An application and accompanying container to simplify and fully automate the setup of a containerized consul cluster. It leverages etcd's service discovery in order to locate other cluster members. 
This works almost the same as the bootstrap scenario with [Atlas](https://atlas.hashicorp.com/) integration, as [described on the bootstrapping page](https://www.consul.io/docs/guides/bootstrapping.html) of consul.

# Warning
It is a proof of concept.

# Container in the Docker Registry
A Docker container image can be found in the public registry as [sogyo/consul-etcd-bootstrapper](https://registry.hub.docker.com/u/sogyo/consul-etcd-bootstrapper/).
```
docker pull sogyo/consul-etcd-bootstrapper
```

# Usage
First, setup a token with your desired cluster size:

```
curl https://discovery.etcd.io/new?size=3
https://discovery.etcd.io/1155b94e64485fc6f7c3d8dae4820306
```

Then, instantiate a bunch of nodes with the given token.
```
docker run --name consul-master --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token 1155b94e64485fc6f7c3d8dae4820306 -ip 10.10.10.100
docker run --name consul-master --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token 1155b94e64485fc6f7c3d8dae4820306 -ip 10.10.20.100
docker run --name consul-master --net=host -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp -p 8400:8400 -p 8500:8500 -p 8600:8600/udp sogyo/consul-auto-bootstrap -token 1155b94e64485fc6f7c3d8dae4820306 -ip 10.10.30.100
```

# Vagrant
A Vagrant demo can be found in the 'vagrant' folder. Just do a 

```vagrant up``` 

and you will get a consul cluster with three server nodes which will exchange information via the ETCD discovery cluster.

In this case, ip configuration is done via cloud-init, so the whole process can also be repeated on an actual cloud provider.

# User Guide

## Creating a Cluster

Creating a Kubernetes cluster is as simple as:
```
$ kind create cluster
```

This will bootstrap a Kubernetes cluster using a pre-built 
[node image][node image] - you can find it on docker hub
[`kindest/node`][kindest/node]. 
If you desire to build the node image yourself see the 
[building image](#building-images) section.
To specify another image use the `--image` flag.

By default, the cluster will be given the name `kind-1`. "1," here, is the
default context name. 
Use the `--name` flag to assign the cluster a different context name.

Once a cluster is created you can use [kubectl][kubectl] to interact with it by
using the configuration file generated by `kind`:
```
export KUBECONFIG="$(kind get kubeconfig-path)"
kubectl cluster-info
```

`kind get kubeconfig-path` returns the location of the generated confguration
file.

`kind` also allows you to customize your Kubernetes cluster by using a kind
configuration file.
For a sample kind configuration file see 
[kind-example-config][kind-example-config].
To use a kind configuration file use the `--config` flag and pass the path to
the file.


## Building Images

`kind` runs a local Kubernetes cluster by using Docker containers as "nodes."
`kind` uses the [`node-image`][node image] to run Kubernetes artifacts, such
as `kubeadm` or `kubelet`.
The `node-image` in turn is built off the [`base-image`][base image], which
installs all the dependencies needed for Docker and Kubernetes to run in a
container.

See [base image](#base-image), for more advanced information information.

Currently, `kind` supports three different ways to build a `node-image`: via
`apt`, or if you have the [Kubernetes][kubernetes] source in your host machine
(`$GOPATH/src/k8s.io/kubernetes`), by using `docker` or `bazel`.
To specify the build type use the flag `--type`.
`kind` will default to using the build type `docker` if none is specified.

```
$ kind build node-image --type apt
```

Similarly as for the base-image command, you can specify the name and tag of
the resulting node image using the flag `--image`.

If you previously changed the name and tag of the base image, you can use here
the flag `--base-image` to specify the name and tag you used.

## Advanced

### Building The Base Image
To build the `base-image` we use the `build` command:
```
$ kind build base-image
```

If you want to specify the path to the base image source files you can use the
`--source` flag.
If `--source` is not specified, `kind` has enbedded the contents of the in
default base image in [`pkg/build/base/sources`][pkg/build/base/sources] and
will use this to build it.

By default, the base image will be tagged as `kindest/base:latest`.
If you want to change this, you can use the `--image` flag.

```
$ kind build base-image --image base:v0.1.0
```


[node image]: ../design/base-image.md
[base image]: ../design/node-image.md
[kind-example-config]: ./kind-example-config.yaml
[pkg/build/base/sources]: ./../../pkg/build/base/sources
[kubernetes]: https://github.com/kubernetes/kubernetes
[kindest/node]: https://hub.docker.com/r/kindest/node/
[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/

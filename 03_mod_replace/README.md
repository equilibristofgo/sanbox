# go.mod replace vs go mod vendor for Dockerfile

In this sandbox, we play with the options we have when in our go.mod have a replace (outside root folder of Dockerfile)

We take as initial example [this](https://www.codervlogger.com/dockerfile-for-a-go-project-with-mod-replace-directive/) great explanation

## Pkg
This module contains some commons function, that in this case is shared by other two modules (here there are other ways to share, but let it in aside)
consider this package as private module, that aren't upload to internet (git) or public repository (yes, yes... this move us to goproxy and all this stuff but let's in this)
what happend inside docker build when it can't to obtain some dependencies that we use as replace.

## Echo
In this module we can se the problem see [logs](./echo/error.log) check how Dockerfile crash.

## Iris
But here we do the trick... vendoring... we use a [Makefile](./iris/Makefile) to organice previus activities to docker build, and there we vendorize
* We could keep it hidden using some command in Makefile
* We use different approach in [Dockerfile](./iris/Dockerfile) to build
* ...


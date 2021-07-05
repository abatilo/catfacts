# catfacts

[![Continuous Deployment](https://github.com/abatilo/catfacts/actions/workflows/cd.yml/badge.svg)](https://github.com/abatilo/catfacts/actions/workflows/cd.yml)

Visit [https://catfacts.aaronbatilo.dev](https://catfacts.aaronbatilo.dev) to
use this application.

I wanted to play around with the Twilio API and the idea that came to mind was
to make a website and backend that would let me subscribe to a list of random
cat facts like how they used to do many moons ago.

## Contributing

Development of this application is done entirely within a local Kubernetes
cluster. Tools are all managed using [asdf-vm](https://asdf-vm.com/#/)

For convenience, there's a `Makefile` that will call to `asdf` and build the
local Kubernetes for you.

```
⇒  make
help                           View help information
asdf-bootstrap                 Install all tools through asdf-vm
kind-bootstrap                 Create a Kubernetes cluster for local development
helm-bootstrap                 Update used helm repositories
bootstrap                      Perform all bootstrapping to start your project
clean                          Delete local dev environment
up                             Run a local development environment
down                           Shutdown local development and free those resources
psql                           Opens a psql shell to the local postgres instance
```

If you type `make up`, and have `asdf` installed, then all of the tools that
are in `.tool-versions` should be installed, then a local Kubernetes cluster
will be created for you and [tilt](https://tilt.dev/) will load all the
applications with hot reloading into your local Kubernetes cluster.

The first time you do this, expect it to take several minutes as building the
cluster and building each application development container will take a while.
Since development is done within containers themselves, versions of runtimes
that are available on your host operating system are ignored. The development
environment is always reproducible and consistent.

Navigate to `http://localhost:8000` to view and use the application.

After it's all up and running, as you make changes to your code locally, the
applications will reload.

## Pulumi

Pulumi is used for managing secrets locally as well as managing remote
resources for deployment. If you'd like to run this application yourself,
you'll need the appropriate `PULUMI_ACCESS_TOKEN` or you'll need to replace the
existing stacks that are managed under `deployment/pulumi` with your own stack
files.

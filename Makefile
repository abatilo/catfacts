SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

PROJECT_NAME = cf
REGISTRY_PORT = 37893

.PHONY: help
help: ## View help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: asdf-bootstrap
asdf-bootstrap: ## Install all tools through asdf-vm
	asdf plugin-add ctlptl  || asdf install ctlptl
	asdf plugin-add golang  || asdf install golang
	asdf plugin-add helm    || asdf install helm
	asdf plugin-add kind    || asdf install kind
	asdf plugin-add kubectl || asdf install kubectl
	asdf plugin-add nodejs  || asdf install nodejs
	asdf plugin-add pulumi  || asdf install pulumi
	asdf plugin-add tilt    || asdf install tilt
	asdf plugin-add yarn    || asdf install yarn

.PHONY: kind-bootstrap
kind-bootstrap: ## Create a Kubernetes cluster for local development
	ctlptl get registry $(PROJECT_NAME)-registry || ctlptl create registry $(PROJECT_NAME)-registry --port=$(REGISTRY_PORT)
	ctlptl get cluster kind-$(PROJECT_NAME) || ctlptl create cluster kind --name kind-$(PROJECT_NAME) --registry=$(PROJECT_NAME)-registry

.PHONY: helm-bootstrap
helm-bootstrap: asdf-bootstrap kind-bootstrap ## Update used helm repositories
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add traefik https://helm.traefik.io/traefik
	helm repo update # Make sure that tilt can pull the latest helm chart versions

.PHONY: bootstrap
bootstrap: asdf-bootstrap kind-bootstrap helm-bootstrap ## Perform all bootstrapping to start your project

.PHONY: clean
clean: ## Delete local dev environment
	ctlptl get cluster kind-$(PROJECT_NAME) && ctlptl delete cluster kind-$(PROJECT_NAME)
	ctlptl get registry $(PROJECT_NAME)-registry && ctlptl delete registry $(PROJECT_NAME)-registry

.PHONY: up
up: bootstrap ## Run a local development environment
	tilt up --context kind-$(PROJECT_NAME) --file ./build/Tiltfile --hud

.PHONY: down
down: ## Shutdown local development and free those resources
	tilt down --context kind-$(PROJECT_NAME) --file ./build/Tiltfile

.PHONY: psql
psql: ## Opens a psql shell to the local postgres instance
	kubectl --context kind-$(PROJECT_NAME) exec -it postgresql-postgresql-0 -- bash -c "PGPASSWORD=local_password psql -U postgres"

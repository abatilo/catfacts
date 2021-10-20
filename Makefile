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
	asdf plugin-add golang  || asdf install golang
	asdf plugin-add helm    || asdf install helm
	asdf plugin-add k3d     || asdf install k3d
	asdf plugin-add kubectl || asdf install kubectl
	asdf plugin-add nodejs  || asdf install nodejs
	asdf plugin-add pulumi  || asdf install pulumi
	asdf plugin-add tilt    || asdf install tilt
	asdf plugin-add yarn    || asdf install yarn

.PHONY: k8s-bootstrap
k8s-bootstrap: ## Create a Kubernetes cluster for local development
	k3d registry create $(PROJECT_NAME).localhost --port=$(REGISTRY_PORT) || echo "Registry already exists"
	k3d cluster create $(PROJECT_NAME) --registry-use k3d-$(PROJECT_NAME).localhost:$(REGISTRY_PORT) || echo "Cluster already exists"

.PHONY: helm-bootstrap
helm-bootstrap: asdf-bootstrap k8s-bootstrap ## Update used helm repositories
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add traefik https://helm.traefik.io/traefik
	helm repo update # Make sure that tilt can pull the latest helm chart versions

.PHONY: bootstrap
bootstrap: asdf-bootstrap k8s-bootstrap helm-bootstrap ## Perform all bootstrapping to start your project

.PHONY: clean
clean: ## Delete local dev environment
	k3d cluster delete $(PROJECT_NAME) || echo "No cluster found"
	k3d registry delete $(PROJECT_NAME).localhost || echo "No registry found"

.PHONY: up
up: bootstrap ## Run a local development environment
	tilt up --context k3d-$(PROJECT_NAME) --file ./build/Tiltfile --hud

.PHONY: down
down: ## Shutdown local development and free those resources
	tilt down --context k3d-$(PROJECT_NAME) --file ./build/Tiltfile

.PHONY: psql
psql: ## Opens a psql shell to the local postgres instance
	kubectl --context k3d-$(PROJECT_NAME) exec -it postgresql-postgresql-0 -- bash -c "PGPASSWORD=local_password psql -U postgres"

# Default variables

REGISTRY_AUTH_FILE := ocp_pull_secret.json
IMAGE_TDX ?= quay.io/bpradipt/rhcos-layer/ocp-4.20:latest-tdx
IMAGE_SNP ?= quay.io/bpradipt/rhcos-layer/ocp-4.20:latest-snp

TEE ?= tdx

# Default target
.PHONY: build
build: check-podman get-pull-secret get-ocp-release-image ## Default target: Build the RHCOS layer
	@echo "Building RHCOS layer..."
	if [ "$(TEE)" = "tdx" ]; then \
		podman build --authfile $(REGISTRY_AUTH_FILE) --build-arg OCP_RELEASE_IMAGE=$(shell cat ./.ocp-release-image) \
		--secret id=org,src=org_secret.txt --secret id=key,src=key_secret.txt \
		-t $(IMAGE_TDX) -f tdx/Containerfile . ;\
	elif [ "$(TEE)" = "snp" ]; then \
		podman build --authfile $(REGISTRY_AUTH_FILE) --build-arg OCP_RELEASE_IMAGE=$(shell cat ./.ocp-release-image) \
		--secret id=org,src=org_secret.txt --secret id=key,src=key_secret.txt \
		-t $(IMAGE_SNP) -f snp/Containerfile . ;\
	fi

# Target to check if 'jq' is available
.PHONY: check-jq
check-jq: ## Check if 'jq' command is available
	@which jq > /dev/null || (echo "'jq' command not found, please install it" && exit 1)


# Target to check if 'oc' is available
.PHONY: check-oc
check-oc: ## Check if 'oc' command is available
	@which oc > /dev/null || (echo "'oc' command not found, please install it" && exit 1)

# Target to check if 'podman' is available
.PHONY: check-podman
check-podman: ## Check if 'podman' command is available
	@which podman > /dev/null || (echo "'podman' command not found, please install it" && exit 1)

# Target to get the pull secret from 'oc'
.PHONY: get-pull-secret
get-pull-secret: ## Get the pull secret from oc 
	@echo "Fetching pull secret from the cluster..."
	@if [ ! -f "$(REGISTRY_AUTH_FILE)" ]; then \
		echo "\$$REGISTRY_AUTH_FILE is not set. Retrieving the pull secret..."; \
		$(MAKE) download-pull-secret; \
		echo "Downloaded the pull secret to: $(REGISTRY_AUTH_FILE)"; \
	else \
		echo "$(REGISTRY_AUTH_FILE) is already present. Skipping download."; \
	fi

# Target to download the pull secret
.PHONY: download-pull-secret
download-pull-secret: ## Download the pull secret
	@echo "Downloading the pull secret..."
	$(MAKE) check-oc
	$(MAKE) check-jq
	oc get -n openshift-config secret/pull-secret -ojson | jq -r '.data.".dockerconfigjson"' | base64 -d | jq '.' > $(REGISTRY_AUTH_FILE)
	@echo "Pull secret downloaded to auth file:$(REGISTRY_AUTH_FILE)"

# Target to get OCP release image details
.PHONY: get-ocp-release-image
get-ocp-release-image: ## Get the OCP release image to use as base
	@echo "Getting OCP release image to use as base..."
	$(MAKE) check-oc; \
	oc adm release info --image-for rhel-coreos > .ocp-release-image


# Help target
.PHONY: help
help: ## Display this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make [target]\n"} /^[a-zA-Z_-]+:.*##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)


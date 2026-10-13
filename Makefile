
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/infinimesh/platform-operator:latest

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/infinimesh/enterprise-operator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
gen-manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	@unformatted=$$(find . -name '*.go' -not -path './vendor/*' -not -path '*/node_modules/*' -not -path './.tmp/*' -print | xargs gofmt -l); \
	if [ -n "$$unformatted" ]; then \
		echo "Formatting Go files:"; \
		printf '%s\n' "$$unformatted"; \
		printf '%s\n' "$$unformatted" | while IFS= read -r f; do [ -n "$$f" ] && gofmt -w "$$f"; done; \
		echo "fmt passed."; \
	else \
		echo "All Go files are formatted correctly."; \
		echo "fmt passed."; \
	fi

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

generate-manifests:
	kustomize build config -o manifests/operator.yaml
	cp config/crds/infinimesh_v1beta1_platform.yaml ./manifests/crd.yaml

.PHONY: fmt-check workflow-check race test-fuzz code-ql code-ql-summary code-ql-gate commit-check

workflow-check: ## Parse every GitHub Actions workflow
	sh scripts/image-tags_test.sh

fmt-check: ## Fail if any Go file is not gofmt-clean
	@unformatted=$$(find . -name '*.go' -not -path './vendor/*' -not -path '*/node_modules/*' -not -path './.tmp/*' -print | xargs gofmt -l); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not gofmt-formatted:"; \
		printf '%s\n' "$$unformatted"; \
		echo "Run 'make fmt' to fix them."; \
		exit 1; \
	fi

race: ## Run tests with the race detector
	@echo "==> go test -race ./pkg/... ./cmd/..."
	@GOFLAGS=-mod=vendor go test -race ./pkg/... ./cmd/...
	@echo "race passed."

test-fuzz: ## Run Go fuzz tests
	bash scripts/test_fuzz.sh

code-ql: ## Run local CodeQL and write SARIF under .tmp/codeql/
	bash scripts/codeql_local.sh

code-ql-summary: code-ql ## Run CodeQL and print a readable summary from SARIF
	@jq -r '.runs[]?.results[]? | select((.locations[0].physicalLocation.artifactLocation.uri // "") | test("(^|/)vendor/") | not) | "\(.level // "warning")\t\(.ruleId // "no-rule")\t\(.locations[0].physicalLocation.artifactLocation.uri // "unknown"):\(.locations[0].physicalLocation.region.startLine // 0)\t\(.message.text // "no-message")"' .tmp/codeql/*.sarif | sort || true

code-ql-gate: code-ql ## Fail if CodeQL reports any error findings
	@errors=$$(jq '[.runs[]?.results[]? | select((.level // "warning") == "error") | select((.locations[0].physicalLocation.artifactLocation.uri // "") | test("(^|/)vendor/") | not)] | length' .tmp/codeql/*.sarif | awk '{s+=$$1} END{print s+0}'); \
	if [ "$$errors" -gt 0 ]; then \
		echo "CodeQL gate failed: $$errors error finding(s) found."; \
		exit 1; \
	fi; \
	echo "CodeQL gate passed: no error findings found."

commit-check: fmt workflow-check race test-fuzz code-ql-gate ## Format, workflows, race, fuzz, and CodeQL
	@echo "commit-check passed."

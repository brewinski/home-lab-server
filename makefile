# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality

# ==============================================================================
# Go Installation
#
#   You need to have Go version 1.22 to run this code.
#
#   https://go.dev/dl/
#
#   If you are not allowed to update your Go frontend, you can install
#   and use a 1.22 frontend.
#
#   $ go install golang.org/dl/go1.22@latest
#   $ go1.22 download
#
#   This means you need to use `go1.22` instead of `go` for any command
#   using the Go frontend tooling from the makefile.

# ==============================================================================
# Brew Installation
#
#	Having brew installed will simplify the process of installing all the tooling.
#
#	Run this command to install brew on your machine. This works for Linux, Mac and Windows.
#	The script explains what it will do and then pauses before it does it.
#	$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
#
#	WINDOWS MACHINES
#	These are extra things you will most likely need to do after installing brew
#
# 	Run these three commands in your terminal to add Homebrew to your PATH:
# 	Replace <name> with your username.
#	$ echo '# Set PATH, MANPATH, etc., for Homebrew.' >> /home/<name>/.profile
#	$ echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/<name>/.profile
#	$ eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
#
# 	Install Homebrew's dependencies:
#	$ sudo apt-get install build-essential
#
# 	Install GCC:
#	$ brew install gcc

# ==============================================================================
# Install Tooling and Dependencies
#
#   This project uses Docker and it is expected to be installed. Please provide
#   Docker at least 4 CPUs. To use Podman instead please alias Docker CLI to
#   Podman CLI or symlink the Docker socket to the Podman socket. More
#   information on migrating from Docker to Podman can be found at
#   https://podman-desktop.io/docs/migrating-from-docker.
#
#	Run these commands to install everything needed.
#	$ make dev-brew
#	$ make dev-docker
#	$ make dev-gotooling

# ==============================================================================
# Running Test
#
#	Running the tests is a good way to verify you have installed most of the
#	dependencies properly.
#
#	$ make test
#

# ==============================================================================
# Running The Project
#
#	$ make dev-up
#	$ make dev-update-apply
#   $ make token
#   $ export TOKEN=<token>
#   $ make users
#
#   You can use `make dev-status` to look at the status of your KIND cluster.

# ==============================================================================
# Project Tooling
#
#   There is tooling that can generate documentation and add a new domain to
#   the code base. The code that is generated for a new domain provides the
#   common code needed for all domains.
#
#   Generating Documentation
#   $ go run app/tooling/docs/main.go --browser
#   $ go run app/tooling/docs/main.go -out json
#
#   Adding New Domain To System
#   $ go run app/tooling/sales-admin/main.go domain sale

# ==============================================================================
# CLASS NOTES
#
# Kind
# 	For full Kind v0.22 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.22.0
#
# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
# 	$ ./sales-admin genkey
#
# Testing Coverage
# 	$ go test -coverprofile p.out
# 	$ go tool cover -html p.out
#
# Module Call Examples
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/@v/list
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/list
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.info
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.mod
# 	$ curl https://proxy.golang.org/github.com/ardanlabs/conf/v3/@v/v3.1.1.zip
# 	$ curl https://sum.golang.org/lookup/github.com/ardanlabs/conf/v3@v3.1.1
#
# OPA Playground
# 	https://play.openpolicyagent.org/
# 	https://academy.styra.com/
# 	https://www.openpolicyagent.org/docs/latest/policy-reference/

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.22

KIND_CLUSTER    := brewinski-cluster
NAMESPACE       := brewinski 
APP             := volleyball
BASE_IMAGE_NAME := localhost/brewinski/service
SERVICE_NAME    := volley-bot 
VERSION         := 0.0.1

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cespare/reflex@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew instal watch

dev-docker:
	docker pull $(GOLANG)

# ==============================================================================
# Running the app within the local computer

vb:
	go run ./src/metro-volleyball-bot/...

vb-watch:
	reflex -r '\.go' -s -- sh -c 'make vb' 

vb-run-test:
	reflex -r '\.go' -s -- sh -c 'make test && make vb'
	reflex

# ==============================================================================
# Running tests within the local computer
test-race:
	CGO_ENABLED=1 go test -race -count=1 ./src/metro-volleyball-bot/...

test-only:
	CGO_ENABLED=0 go test -count=1 ./src/metro-volleyball-bot/... 

lint:
	CGO_ENABLED=0 go vet ./src/metro-volleyball-bot/... 
	staticcheck -checks=all ./src/metro-volleyball-bot/...

vuln-check:
	govulncheck ./src/metro-volleyball-bot/..

test: test-only lint vuln-check

test-race: test-race vuln-check

test-watch:
	reflex -r '\.go' -s -- sh -c 'make test'

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all


set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

IMAGE_PREFIX := env_var_or_default("IMAGE_PREFIX", "localhost/borg")
IMAGE_VERSION := env_var_or_default("IMAGE_VERSION", "latest")

TOOLS := "droid siegfried tika magika mediainfo jhove verapdf odf-validator ooxml-validator"

default:
	@just --list

build-gui:
	podman build -t "{{IMAGE_PREFIX}}/gui:{{IMAGE_VERSION}}" ./gui

build-server:
	podman build -t "{{IMAGE_PREFIX}}/server:{{IMAGE_VERSION}}" ./server

build-tools:
	for tool in {{TOOLS}}; do \
		podman build -t "{{IMAGE_PREFIX}}/$tool:{{IMAGE_VERSION}}" "./tools/$tool"; \
	done

build-all: build-server build-gui build-tools

push-gui:
	podman push "{{IMAGE_PREFIX}}/gui:{{IMAGE_VERSION}}"

push-server:
	podman push "{{IMAGE_PREFIX}}/server:{{IMAGE_VERSION}}"

push-tools:
	for tool in {{TOOLS}}; do \
		podman push "{{IMAGE_PREFIX}}/$tool:{{IMAGE_VERSION}}"; \
	done

# Start podman compose
podman-up:
    podman-compose up -d --build

# Stop podman compose
podman-down:
    podman-compose down

# Show podman logs
podman-logs:
    podman-compose logs -f

# Run linter of server code
[working-directory: 'server']
lint-server:
    golangci-lint run \
        --output.text.path=stdout \
        --output.text.colors=false \
        --output.text.print-issued-lines=false \
        --output.code-climate.path=gl-code-quality-report.json

# Run tests for server component
[working-directory: 'server']
test-server:
    # CGO_ENABLED=1 -> SQLite wird für Tests benötigt.
    CGO_ENABLED=1 gotestsum \
        --junitfile report.xml \
        --format testname \
        -- -coverprofile=cover.out ./...


push-all: push-server push-gui push-tools

build-and-push-all: build-all push-all

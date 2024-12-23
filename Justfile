# Ships Flight Summarizer!
export PROJECT_NAME := "flight-summarizer"
export GOLANG_VERSION := `grep -E '^go ([0-9]{1,}\.[0-9]{1,})$' $PWD/go.mod | awk '{print $NF}'`

deploy: build test
  #!/usr/bin/env bash
  >&2 echo "Work in progress!"

# builds Flight Summarizer ✈️
build:
  just --one _docker_compose run --rm build

# tests Flight Summarizer ✈️
test:
  just --one _docker_compose run --rm test

# performs Flight Summarizer end-to-end tests
e2e:
  just --one _docker_compose run --rm e2e

_docker_compose *ARGS:
  #!/usr/bin/env bash
  if ! test -e "$HOME/.docker/cli-plugins/docker-compose"
  then
    >&2 echo "ERROR: Docker Compose CLI plugin not found; please install it"
    exit 1
  fi
  docker compose {{ARGS}}

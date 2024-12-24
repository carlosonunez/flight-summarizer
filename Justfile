# Ships Flight Summarizer!
export GOLANG_VERSION := `grep -E '^go ([0-9]{1,}\.[0-9]{1,})$' $PWD/go.mod | awk '{print $NF}'`

deploy: build test
  #!/usr/bin/env bash
  >&2 echo "Work in progress!"

# Starts a local instance of the Flight Summarizer server!
start-server: build
  SERVER_HOST=0.0.0.0 just --one _docker_compose up --build summarizer-server
  just --one _docker_compose down

# builds Flight Summarizer ✈️
build:
  #!/usr/bin/env bash
  for project in cmd/*
  do
    PROJECT_NAME="$(basename $project)" just --one _docker_compose run --rm build
  done

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
  docker-compose --log-level ERROR {{ARGS}}

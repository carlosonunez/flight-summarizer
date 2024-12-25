# Ships Flight Summarizer!
set export := true

export GOLANG_VERSION := `grep -E '^go ([0-9]{1,}\.[0-9]{1,})$' $PWD/go.mod | awk '{print $NF}'`
export BUILD_ARCHS := "linux/arm64 linux/amd64"

# Cleans up our workspace
clean: _remove_buildx_builders

# Deploys Flight Summarizer to Github Container Registry
deploy_to_ghcr: _ensure_gh_creds_available build
  #!/usr/bin/env bash
  set -euo pipefail
  last_commit_sha=$(git log -1 --format=%h)
  tag_at_this_commit=$(git tag --points-at HEAD)
  tags="latest,${last_commit_sha},${tag_at_this_commit}"
  for project in cmd/*
  do
    project_name="$(basename $project)"
    archs_available=$(find $PWD/out -type f | grep -E ".*/summarizer-[a-z]{1,}-[a-z0-9]{1,}" |
      grep -v 'darwin' |
      sed -E "s/.*${project_name}-//" |
      tr '-' '/' |
      tr '\n' ',' |
      sed -E 's/,$//')
    PROJECT_NAME="$project_name" \
      ARCHITECTURES_CSV="$archs_available" \
      BUILD_TAGS="$tags" \
      just --one _docker_compose run --rm deploy-to-ghcr
  done

# Starts a local instance of the Flight Summarizer server!
start-server: build
  SERVER_HOST=0.0.0.0 just --one _docker_compose up --build summarizer-server
  just --one _docker_compose down

# builds Flight Summarizer ✈️
build:
  #!/usr/bin/env bash
  export GOARCH="$$(just --one _arch)" || exit 1
  for project in cmd/*
  do
    for arch in $BUILD_ARCHS
    do
      export GOOS=$(cut -f1 -d '/' <<< "$arch")
      export GOARCH=$(cut -f2 -d '/' <<< "$arch")
      export PROJECT_NAME="$(basename $project)"
      >&2 echo "INFO: Building [$PROJECT_NAME] for $GOOS/$GOARCH"
      just --one _docker_compose run --rm build
    done
  done

# tests Flight Summarizer ✈️
test:
  #!/usr/bin/env bash
  export GOARCH="$(just --one _arch)" || exit 1
  just --one _docker_compose run --rm test

# performs Flight Summarizer end-to-end tests
e2e:
  export GOARCH="$(just --one _arch)" || exit 1
  just --one _docker_compose run --rm e2e

_docker_compose *ARGS:
  #!/usr/bin/env bash
  if ! &>/dev/null docker compose
  then
    >&2 echo "ERROR: Docker Compose CLI plugin not found; please install it"
    exit 1
  fi
  docker-compose --log-level ERROR {{ARGS}}

_arch:
  #!/usr/bin/env bash
  if test -n "$RUNNER_ARCH"
  then echo "${RUNNER_ARCH,,}" && exit 0
  fi
  if grep -iq 'arm' <<< "$(uname -p)"
  then echo "arm64" && exit 0
  fi
  echo "amd64"

_ensure_gh_creds_available:
  #!/usr/bin/env bash
  for key in GITHUB_TOKEN GITHUB_REPOSITORY
  do
    test -n "${!key}" && continue
    >&2 echo "ERROR: Please define $key"
    exit 1
  done

_remove_buildx_builders:
  #!/usr/bin/env bash
  for project in cmd/*
  do
    project_name="$(basename $project)"
    PROJECT_NAME="$project_name" \
      just --one _docker_compose run --rm _remove_buildx_builder || true
  done

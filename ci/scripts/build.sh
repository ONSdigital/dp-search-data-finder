#!/bin/bash -eux

pushd dp-search-data-finder
  make build
  cp build/dp-search-data-finder Dockerfile.concourse ../build
popd

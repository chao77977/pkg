#!/bin/sh

set -e

if [ ! -f "build/env.sh"  ]; then
  echo "$0 must be run from the root of the repository."
  exit 2
fi

workspace="$PWD/build/_workspace"
kpdir="$workspace/src"
if [ ! -L "$kpdir/pkg" ]; then
  mkdir -p $kpdir
  cd $kpdir && ln -sf ../../../. pkg
fi

# Set up the environment to use the workspace.
export GOPATH="$workspace"
export GO111MODULE=on
PWD="$kpdir/pkg"

# Launch the arguments with the configured environment.
cd $PWD && exec "$@"

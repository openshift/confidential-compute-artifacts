#!/bin/bash
set -e

export KATA_VERSION=$(cat /VERSION)
export NVIDIA_DRIVERS_VERSION=$(cat /NVIDIA-VERSION)

if [ -d /host ]; then \
    echo "Copying kata ${KATA_VERSION} artifacts to /host"
    \cp /kata-initrds.tar.gz /host/kata-initrds-${KATA_VERSION}-nvidia-${NVIDIA_DRIVERS_VERSION}-$(arch).tar.gz
    \cp /kata-osbuilder.tar.gz /host/kata-osbuilder-${KATA_VERSION}.tar.gz
    \cp /kata-logs.tar.gz /host/kata-logs-${KATA_VERSION}.tar.gz
    echo "All done!"
else
    echo "Error: /host directory not found."
    echo "Please run with: podman run --rm -v \$(pwd):/host"
    exit 1
fi

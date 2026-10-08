#!/bin/bash
set -e

export KATA_VERSION=$(cat /VERSION)
export KERNEL_VERSION=$(cat /KERNEL-VERSION)

if [ -d /host ]; then \
    echo "Copying kata ${KATA_VERSION} artifacts to /host"
    # On s390x /NVIDIA-VERSION does not exist — use a kernel-only filename.
    # On all other arches (x86_64) include the NVIDIA driver version in the name.
    if [ "$(arch)" = "s390x" ]; then
        \cp /kata-initrds.tar.gz /host/kata-initrds-${KATA_VERSION}-kernel-${KERNEL_VERSION}-$(arch).tar.gz
    else
        export NVIDIA_DRIVERS_VERSION=$(cat /NVIDIA-VERSION)
        \cp /kata-initrds.tar.gz /host/kata-initrds-${KATA_VERSION}-nvidia-${NVIDIA_DRIVERS_VERSION}-kernel-${KERNEL_VERSION}-$(arch).tar.gz
    fi
    \cp /kata-osbuilder.tar.gz /host/kata-osbuilder-${KATA_VERSION}.tar.gz
    \cp /kata-logs.tar.gz /host/kata-logs-${KATA_VERSION}.tar.gz
    \cp /kata-source.tar.gz /host/kata-containers-${KATA_VERSION}.tar.gz
    \cp /kata-vendor.tar.gz /host/kata-containers-${KATA_VERSION}-rh-vendor.tar.gz
    echo "All done!"
else
    echo "Error: /host directory not found."
    echo "Please run with: podman run --rm -v \$(pwd):/host"
    exit 1
fi

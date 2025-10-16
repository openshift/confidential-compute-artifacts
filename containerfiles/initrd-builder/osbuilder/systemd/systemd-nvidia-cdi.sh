#!/bin/bash

# Load drivers
nvidia-ctk -d system create-device-nodes --control-devices --load-kernel-modules

# Start persistenced
nvidia-persistenced --verbose --uvm-persistence-mode

# Generate NVIDIA CDI spec
nvidia-ctk cdi generate --output=/var/run/cdi/nvidia.yaml

# Set confidential compute to ready state
nvidia-smi conf-compute -srs 1

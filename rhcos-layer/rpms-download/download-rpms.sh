#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 1 ]; then
    echo "Usage: $0 <rpms.yaml>"
    exit 1
fi

YAML_FILE="$1"
if [ ! -f "$YAML_FILE" ]; then
    echo "Error: File '$YAML_FILE' not found."
    exit 1
fi

num_sources=$(yq e '. | length' "$YAML_FILE")

for i in $(seq 0 $((num_sources - 1))); do
    base_url=$(yq e ".[$i].base_url" "$YAML_FILE")
    use_subdir=$(yq e ".[$i].use_subdir" "$YAML_FILE")
    num_rpms=$(yq e ".[$i].rpms | length" "$YAML_FILE")

    for j in $(seq 0 $((num_rpms - 1))); do
        rpm=$(yq e ".[$i].rpms[$j]" "$YAML_FILE")

        if [ "$use_subdir" == "true" ]; then
            first_letter="${rpm:0:1}"
            url="${base_url}/${first_letter}/${rpm}"
        else
            url="${base_url}/${rpm}"
        fi

        echo "Downloading $rpm from $url"
        curl -O "$url" || echo "Failed to download $rpm"
    done
done

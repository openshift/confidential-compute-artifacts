## Pre-build initrds

1. Update the arguments in the `argfile.conf` accordingly.
   `KERNEL_VERSION` must be set to the **base kernel version without the arch suffix**
   (e.g. `6.12.0-124.21.1.el10_1`). The correct suffix (`.x86_64` / `.s390x`) is
   appended automatically at build time from the target platform.

2. Create the `org_secret.txt` and `key_secret.txt` with your ORG_ID and ACTIVATION_KEY for subscription
```
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```

### Single-arch build — local testing (x86_64 or s390x)

Use `-t` for local builds. This gives a directly usable named tag with the full image content.

**x86_64:**
```
podman build . --no-cache \
    --platform linux/x86_64 \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    -t kata-initrds:multi
```

**s390x:**
```
podman build . --no-cache \
    --platform linux/s390x \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    -t kata-initrds:s390x
```

On **x86_64** this produces:
- `kata-initrds-<version>-nvidia-<nvidia-version>-kernel-<kernel>-x86_64.tar.gz`

On **s390x** this produces (no NVIDIA):
- `kata-initrds-<version>-kernel-<kernel>-s390x.tar.gz`

### Multi-arch manifest build — for registry push only

Use `--manifest` only when building for a registry push. With `--manifest` the per-arch
images are stored untagged locally; the tag is just a tiny OCI index pointer (~600 B),
not the real image. For local testing always use `-t` above.

Requires QEMU user-static emulation to be registered on the host (for cross-arch builds):
```
podman build . --no-cache \
    --platform linux/x86_64,linux/s390x \
    --manifest kata-initrds:latest \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey

# Then push the manifest list to a registry
podman manifest push kata-initrds:latest docker://registry.example.com/kata-initrds:latest
```

### Extract tarballs from an already-built image

```
# From a -t build
podman run --rm -v $PWD:/host:z kata-initrds:x86_64
podman run --rm -v $PWD:/host:z kata-initrds:s390x

# From a --manifest build (use --platform to resolve the right arch)
podman run --rm --platform linux/x86_64 -v $PWD:/host:z localhost/kata-initrds:latest
podman run --rm --platform linux/s390x  -v $PWD:/host:z localhost/kata-initrds:latest
```

## Debugging the kata-osbuilder.sh script

There are multiple ways, one of them is to generate an image building only the first two stages (i.e. `initrd-builder-setup`)
and then use the resulting image to run/test the `kata-osbuilder.sh` manually.

```
podman build . --no-cache \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    --target initrd-builder-setup \
    -t initrd-builder-setup:1.0
podman run -ti -v $PWD:/host:z localhost/initrd-builder-setup:1.0 /bin/bash
# then run osbuilder/kata-osbuilder.sh as desired
```

## Pre-build initrds

1. Update the arguments in the `argfile.conf` accordingly
2. Create the `org_secret.txt` and `key_secret.txt` with your ORG_ID and ACTIVATION_KEY for subscription
```
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```
3. Pre-build initrds for the target architecture (`linux/amd64` or `linux/s390x`).
   The target architecture dynamically selects the corresponding `PRE_BUILD_INITRDS_<arch>` from `argfile.conf`.
   Tarballs are stored in both the container image and `$PWD`.

**x86_64:**
```
podman build . --no-cache \
    --platform linux/amd64 \
    -v $PWD:/host:z \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    -t kata-initrds:x86_64
```

**s390x:**
```
podman build . --no-cache \
    --platform linux/s390x \
    -v $PWD:/host:z \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    -t kata-initrds:s390x
```

4. Tarballs can also be extracted from the container image
```
# For x86_64
podman run --rm --platform linux/amd64 -v $PWD:/host:z kata-initrds:x86_64

# For s390x
podman run --rm --platform linux/s390x -v $PWD:/host:z kata-initrds:s390x
```

## Debugging the kata-osbuilder.sh script

There are multiple ways, one of them is to generate an image building only the first two stages (i.e. `initrd-builder-setup`)
and then use the resulting image to run/test the `kata-osbuilder.sh` manually.

```
podman build . --no-cache \
    --platform linux/amd64 \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    --target initrd-builder-setup \
    -t initrd-builder-setup:1.0
podman run -ti --rm --platform linux/amd64 -v $PWD:/host:z localhost/initrd-builder-setup:1.0 /bin/bash
# then run osbuilder/kata-osbuilder.sh as desired
```

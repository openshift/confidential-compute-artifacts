## Pre-build initrds

1. Update the arguments in the `argfile.conf` accordingly
2. Create the `org_secret.txt` and `key_secret.txt` with your ORG_ID and ACTIVATION_KEY for subscription
```
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```
4. Pre-build initrds for the runtime classes defined in the step 1. 
   The tarballs will be stored in both locations: within the container image and in $PWD
```
podman build . --no-cache \
    -v $PWD:/host:z \
    --build-arg-file=./argfile.conf \
    -v $PWD/org_secret.txt:/activation-key/org \
    -v $PWD/key_secret.txt:/activation-key/activationkey \
    -t kata-initrds:1.0
```

## Debugging the kata-osbuilder.sh script

There are multiple ways, one of them is to generate an image building only the first two stages (i.e. `initrd-builder-setup`)
and then use the resulting image to run/test the `kata-osbuilder.sh` manually.

```
podman build . --no-cache \
    --build-arg-file=./argfile.conf \
    --secret id=org,src=org_secret.txt \
    --secret id=key,src=key_secret.txt \
    --target initrd-builder-setup \
    -t initrd-builder-setup:1.0
podman run -ti -v $PWD:/host:z localhost/initrd-builder-setup:1.0 /bin/bash
# then run osbuilder/kata-osbuilder.sh as desired
```

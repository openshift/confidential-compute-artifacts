## Create kata-initrds image

1. Download/generate the kata-containers rpm and copy it to the same directory as the `Containerfile`
2. Update the arguments in the `argfile.conf` accordingly
3. Create the `org_secret.txt` and `key_secret.txt` with your ORG_ID and ACTIVATION_KEY for subscription
```
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```
4. Run the command below to build the kata-initrds image
```
podman build . --no-cache \
    -v $PWD:/host \
    --build-arg-file=./argfile.conf \
    --secret id=org,src=org_secret.txt \
    --secret id=key,src=key_secret.txt \
    -t kata-initrds:1.0
```
If `-v $PWD:/host` is provided, the file `kata-initrds.tar.gz` will also be copied to $PWD.

## Debugging the kata-osbuilder.sh script

There are multiple ways, one of them is to generate an image building only the first two stages (i.e. `initrds-builder-setup`)
and then use the resulting image to run/test the `kata-osbuilder.sh` manually.

```
podman build . --no-cache \
    --build-arg-file=./argfile.conf \
    --secret id=org,src=org_secret.txt \
    --secret id=key,src=key_secret.txt \
    --target initrds-builder-setup \
    -t initrds-builder-setup:1.0
podman run -ti -v $PWD:/host localhost/initrds-builder-setup:1.0 /bin/bash
# then run /usr/libexec/kata-containers/osbuilder/kata-osbuilder.sh as desired
```

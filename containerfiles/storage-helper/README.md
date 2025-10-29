## Build image

- Set subscription-manager details

This is an optional step if you are running on an unregistered
host. You'll also need to uncomment the subscription-manager instructions
in the Containerfile. 

```sh
echo <org_id> > org_secret.txt
echo <activation_key> > key_secret.txt
```

- Set registry namespace details

Example: 

```sh
export registry=quay.io/openshift_sandboxed_containers
```

- Build the image

```sh

podman build  --secret id=org,src=org_secret.txt --secret id=key,src=key_secret.txt \
     -t $registry/storage-helper:latest -f Containerfile .
```

# Initrd Customisation

You can use the Containerfile in this folder to create a custom
RHCOS image layer with Kata initrd containing CA certs and Registry auth file.

- Copy the required certs to a single file (eg. tls.crt)

- Copy all the required registry auths to a single json file (eg. auth.json) . The file must be in
  [auth-json format](https://github.com/containers/image/blob/main/docs/containers-auth.json.5.md)

- Build the image

```sh
sudo podman build --build-arg RHCOS_COCO_IMAGE=<rhcos-coco-image> \
             --build-arg CA_CERTS_FILE=tls.crt \
             --build-arg REGISTRY_AUTH_FILE=auth.json \
             --env CA_CERTS_FILE=tls.crt \
             --env REGISTRY_AUTH_FILE=auth.json \
             -t <new-rhcos-coco-image> -f Containerfile .
```

- Push it to registry

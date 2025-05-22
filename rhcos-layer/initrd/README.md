# Initrd Customisation

You can use the Containerfile in this folder to create a custom
RHCOS image layer with Kata initrd containing CA certs and Registry auth file.

```sh
sudo podman build --build-arg RHCOS_COCO_IMAGE=<rhcos-coco-image> \
             --build-arg CA_CERTS_FILE=<ca-certs-file> \
             --build-arg REGISTRY_AUTH_FILE=<registry-auth-file> \
             -t <new-rhcos-coco-image> -f Containerfile .
```

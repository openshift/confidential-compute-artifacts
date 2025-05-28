# Initrd Customisation

You can use the Containerfile in this folder to create a custom
RHCOS image layer with Kata initrd containing CA certs and Registry auth file.

```sh
sudo podman build --build-arg RHCOSIMAGE=<rhcos-coco-image> \
             --build-arg CERTSFILE=<ca-certs-file> \
             --build-arg AUTHFILE=<registry-auth-file> \
             -t <new-rhcos-coco-image> -f Containerfile .
```

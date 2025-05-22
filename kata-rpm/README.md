## Create kata-containers rpm from a given srpm image

#### Run the following command to create a container image and extract the rpm into the current directory.
```sh
podman build --build-arg SRPM_IMAGE=<image.srpm> \
             --build-arg ORG=<rh-subsmgr-org> \
             --build-arg ACTIVATIONKEY=<rh-subsmgr-activationkey> \
             -t <build-image> -v $PWD:/host -f Containerfile .
```
Where:
```
    image-srpm                    = Source rpm of Kata-containers
    rh-subscription-org           = Organization's identifier in the Red Hat system
    rh-subscription-activationkey = Activation key to register systems for subscriptions
    build-image                   = Name (i.e. repo:tag) of the container image
```

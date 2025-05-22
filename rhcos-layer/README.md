# Creating RHCOS image Layer

This folder contains artifacts to create RHCOS image layer for SNP and TDX.

The RHCOS image layer contains the required kernel version, kata and qemu rpms.

## Download the required rpms

For SNP

```sh
cd rpms-download
./download-rpms.sh rhel-rpms.yaml
```

Copy the downloaded rpms to the `snp/rpms` folder

For TDX

```sh
cd rpms-download
./download-rpms.sh centos-rpms.yaml
```

Copy the downloaded rpms to the `tdx/rpms` folder

## Build the RHCOS layer

### Prerquisites

- make
- podman
- jq
- Depending on your build environment, you might have to spawn a root shell using `sudo` before executing
  the build commands

### With access to OCP cluster

You must have `oc` configured to work with the cluster.

Build for SNP

```sh
make TEE=snp build
```

To add the CA certs file (say ca-certs) and Registry auth file (say config.json) to the initrd,
run the following.

```sh
make TEE=snp CA_CERTS=ca-certs REG_AUTH=config.json build
```

Build for TDX

```sh
make TEE=tdx build
```

To add the CA certs file (say ca-certs) and Registry auth file (say config.json) to the initrd,
run the following.

```sh
make TEE=tdx CA_CERTS=ca-certs REG_AUTH=config.json build
```

### Without access to OCP cluster

You can also directly build using podman or docker.

- Download the OCP pull secret from console.redhat.com
- Run the following command:

This uses RHCOS base image for OCP 4.18

Build for SNP

```sh
podman build --authfile /tmp/pull-secret.json \
   --build-arg OCP_RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:bdaa82a5a1df84ee304cbf842c80278e2286fede509664c5f0cf9c93c0992658 \   
   -t snp-image -f snp/Containerfile .
```

Build for TDX

```sh
podman build --authfile /tmp/pull-secret.json \
   --build-arg OCP_RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:bdaa82a5a1df84ee304cbf842c80278e2286fede509664c5f0cf9c93c0992658 \
   -t tdx-image -f tdx/Containerfile .
```

If you want to add custom CA certs and registry auth file, then add the following arguments to
podman build `--build-arg CA_CERTS_FILE=ca-certs --build-arg REGISTRY_AUTH_FILE=config.json`

For example, the build for SNP will be the following

```sh
podman build --authfile /tmp/pull-secret.json \
   --build-arg OCP_RELEASE_IMAGE=quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:bdaa82a5a1df84ee304cbf842c80278e2286fede509664c5f0cf9c93c0992658 \
   --build-arg CA_CERTS_FILE=ca-certs --build-arg REGISTRY_AUTH_FILE=config.json \ 
   -t snp-image -f snp/Containerfile .
```

## Customising the Kata initrd in RHCOS image layer

In some cases, you might need to customise the Kata initrd in an existing RHCOS image layer.
Following the steps mentioned [here](initrd/README.md).
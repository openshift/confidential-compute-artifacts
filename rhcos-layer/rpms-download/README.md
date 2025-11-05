# Introduction

You can use the script to download the required rpms and copy it under {tdx,snp}/rpms folder

Further, you can copy the rpms into a container image and save it for versioning.

## Download the rpms and create a tar

Download the RHEL rpms used for SNP

```sh
./download-rpms.sh rhel-rpms.yaml
tar czvf snp-rpms-0.2.0.tar.gz *.rpm

Download the CentOS rpms used for TDX

```sh
./download-rpms.sh centos-rpms.yaml
tar czvf tdx-rpms-0.2.0.tar.gz *.rpm
```

## Add the rpms to a container image and save it to registry

Following is the Containerfile:

```sh
FROM registry.access.redhat.com/ubi10/ubi:latest

COPY tdx-rpms-0.2.0.tar.gz /tdx-rpms-0.2.0.tar.gz
COPY snp-rpms-0.2.0.tar.gz /snp-rpms-0.2.0.tar.gz
```

Set the container registry details

For example:

```sh
export RPM_IMAGE=quay.io/myuser/rhcos-layer/rpms:0.2.0
```

Build and push the image

```sh
podman build -t $RPM_IMAGE -f Containerfile .
podman push $RPM_IMAGE
```

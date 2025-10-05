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
- Create `org_secret.txt` and `key_secret.txt` file with your
ORG_ID and ACTIVATION_KEY for subscription

```bash
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```

### With access to OCP cluster

You must have `oc` configured to work with the cluster.

Build for SNP

```sh
make TEE=snp build
```

Build for TDX

```sh
make TEE=tdx build
```

### Without access to OCP cluster

You can also directly build using podman or docker.

- Download the OCP pull secret from console.redhat.com


Build for SNP

```sh
podman build --authfile ocp_pull_secret.json \
   --secret id=org,src=org_secret.txt --secret id=key,src=key_secret.txt \
   -t snp-image -f snp/Containerfile .
```

Build for TDX

```sh
podman build --authfile ocp_pull_secret.json \
   --secret id=org,src=org_secret.txt --secret id=key,src=key_secret.txt \
   -t tdx-image -f tdx/Containerfile .
```

You can also set the OCP_RELEASE_IMAGE and the PAYLOAD_IMAGE via build-arg

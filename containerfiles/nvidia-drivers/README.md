## Build image with signed nvidia drivers for RHEL 10

1. Update the arguments in the `argfile.conf` accordingly
2. Create the `org_secret.txt` and `key_secret.txt` with your ORG_ID and ACTIVATION_KEY for subscription
```
echo <ORG_ID> > org_secret.txt
echo <ACTIVATION_KEY> > key_secret.txt
```
3. Run the command below to build the nvidia-driver image
```
podman build . --no-cache \
    --build-arg-file=./argfile.conf \
    --secret id=org,src=org_secret.txt \
    --secret id=key,src=key_secret.txt \
    -t rhel10-nvidia-drivers:550.163.01
```

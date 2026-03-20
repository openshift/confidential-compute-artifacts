# luks-helper

A Go binary that runs as a native sidecar in a Kubernetes environment to
transparently provide LUKS-encrypted block storage to the main container.

## How it works

```
┌─ Pod (shareProcessNamespace: true ) ─────────────────────────────────────────┐
│                                                                              │
│  ┌─ format-disk (native sidecar) ────────────────────────────────────────┐   │
│  │  1. Check if PVC has a LUKS header                                    │   │
│  │     - Yes → open with passphrase                                      │   │
│  │     - No  → format with LUKS, then open                               │   │
│  │  2. mkfs.xfs on /dev/mapper/mapper  (skipped if already XFS)          │   │
│  │  3. mount /dev/mapper/mapper → /mnt/storage                           │   │
│  │  4. Write own PID to /dev/shm/luks-helper.pid                        │   │
│  │  5. Touch /mnt/storage/.ready                                         │   │
│  │  6. Block on SIGTERM (keeps mount alive for pod lifetime)             │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─ check-ready (init container) ────────────────────────────────────────┐   │
│  │  Reads PID from /dev/shm/luks-helper.pid                             │   │
│  │  Polls /proc/<pid>/root/mnt/storage/.ready                            │   │
│  │  Exits 0 when found → unblocks main container                         │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─ main container ──────────────────────────────────────────────────────┐   │
│  │  postStart: reads PID, symlinks /data → /proc/<pid>/root/mnt/storage  │   │
│  │  App reads/writes /data as if it were a normal mount                  │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────────────┘
```

The LUKS mount lives inside the sidecar's mount namespace. Other containers
access it via `/proc/<pid>/root` over the shared PID namespace — no extra
privileges or volume mounts required on the application side beyond what is
already needed for the Kata VM.

## Commands

### `format-disk`

Formats (if needed), opens, and mounts a LUKS-encrypted block device. Runs as a native sidecar (`restartPolicy: Always`).

```
luks-helper format-disk [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--device` | `/dev/block-device` | Raw block device to encrypt |
| `--mapper-name` | `mapper` | LUKS device-mapper name |
| `--mapper-path` | `/dev/mapper/mapper` | Resulting mapper device path |
| `--mount-target` | `/mnt/storage` | Where to mount the decrypted volume |
| `--pid-file` | `/dev/shm/luks-helper.pid` | Path to write own PID for handoff |

Requires the `PASS` environment variable (LUKS passphrase) to be set.

### `wait-ready`

Waits for `format-disk` to signal that the volume is mounted and ready. Runs as
a regular init container — exits 0 on success, exits 1 on timeout.

```
luks-helper wait-ready [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--mount-target` | `/mnt/storage` | Mount point to watch |
| `--pid-file` | `/dev/shm/luks-helper.pid` | PID file written by `format-disk` |
| `--max-wait` | `3600` | Timeout in seconds |

### `sleep`

Blocks forever without doing anything. Use this to debug a failing pod by
overriding the container command, then `kubectl exec` in to run `cryptsetup`,
`blkid`, or `mount` manually.

```
luks-helper sleep
```

```bash
# In pod.yaml, temporarily override:
command: ["/usr/local/bin/luks-helper", "sleep"]

# Then exec in:
kubectl exec -it <pod> -c format-disk -- /bin/sh
```

## Requirements

- `shareProcessNamespace: true` on the pod
- A CSI StorageClass that supports `volumeMode: Block` (e.g. `ocs-storagecluster-ceph-rbd` on ODF)
- A Secret containing the LUKS passphrase

## Build

```bash
go build -o luks-helper .
```

```bash
# Container image
docker build -t your-registry/luks-helper:latest .
```

## Deploy

See `pod.yaml` for a complete example including the Secret, PVC, and Pod spec.

The passphrase secret must exist before the pod starts:

```bash
kubectl create secret generic luks-passphrase --from-literal=passphrase='<your-passphrase>'
```


# Box

It's a package that contains a default configration and files for VM, like Vagrant's box. It is compressed by `tar` with `gzip`.

## `<BOX NAME>.xhv`

It contains the following files.

### config.yml
```yaml
# OS type of VM
type: linux

# Boot configurations
kernel: vmlinuz
initrd: initrd.gz
cmdline: earlyprintk=serial console=ttyS0 acpi=off

# SSH configurations
ssh:
  username:
  password:
  key-path:

# VM configurations for xhyve arguments
acpi: false
cpus: 1
memory: 1G
hdd:
  - "4,virtio-blk,disk.img"
net:
    - "2:0,virtio-net"
pci:
    - "0:0,hostbridge"
    - "31,lpc"
lpc:
    - "com1,stdio"
```

### A kernel image file (if type == linux)

### An initrd image file (if type == linux)

### HDD image files (if needed)

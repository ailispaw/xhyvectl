# Box

It's a package that contains a default configration and files for VM, like Vagrant's box. It is compressed by `tar` with `gzip`.

## `<BOX NAME>.box`

It contains the following files.

### config.yml
```yaml
type: linux

kernel: vmlinuz
initrd: initrd.gz
cmdline: earlyprintk=serial console=ttyS0 acpi=off

acpi: false
cpus: 1
memory: 1G
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

### CD image files (if needed)

### HDD image files (if needed)

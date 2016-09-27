# Commands

### `xhyvectl package <BOX NAME> <DIR> [--force|-f]`

Package `<DIR>/*` files into `<BOX NAME>.xhv`.

### `xhyvectl box list`

List installed boxes

### `xhyvectl box install <BOX NAME> <BOX FILE|URL> [--force|-f]`

Install a box file into your local HDD.  
It will be `~/.xhyvectl/boxes/<BOX NAME>.xhv`.

### `xhyvectl box update <BOX NAME> <BOX FILE|URL>`

Update/Replce the box with a box file.

### `xhyvectl box unintall <BOX NAME>...`

Remove box(es) from `~/.xhyvectl/boxes/`.

### `xhyvectl init <VM NAME> <BOX NAME> [--force|-f]`

Initialize a VM with `<BOX NAME>`.  
It will be placed at `~/.xhyvectl/vms/<VM NAME>/`.

### `xhyvectl up <VM NAME> [--cups|-c <# of CPUs>] [--memory|-m <Memory Size>]`

Boot up a VM with parameters for xhyve

### `xhyvectl status <VM NAME>`

Show status of the VM.

### `xhyvectl ssh <VM NAME> [flag]...`

SSH into the running VM.

### `xhyvectl ip <VM NAME>`

Get an IP address of the running VM.

### `xhyvectl halt <VM NAME>`

Stop the running VM.

### `xhyvectl restart <VM NAME>`

Restart the running VM.

### `xhyvectl destroy <VM NAME>`

Stop the running VM and remove a VM from `~/.xhyvectl/vms/`.

### `xhyvectl version`

Show the version information

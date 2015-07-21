# Commands

### `xhyvectl package <BOX NAME> <DIR>`

Package `<DIR>/*` files into `<BOX NAME>.xhv`.

### `xhyvectl install <BOX NAME> <BOX FILE|URL>`

Install a box file into your local HDD.  
It will be placed at `~/.xhyvectl/boxes/<BOX NAME>/`.

### `xhyvectl update <BOX NAME> <BOX FILE|URL>`

Update/Replce the box with a box file.

### `xhyvectl unintall <BOX NAME>`

Remove a box from `~/.xhyvectl/boxes/`.

### `xhyvectl init <VM NAME> <BOX NAME>`

Initialize a VM with `<BOX NAME>`.  
It will be placed at `~/.xhyvectl/vms/<VM NAME>/`.

### `xhyvectl up <VM NAME> [flag]...`

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

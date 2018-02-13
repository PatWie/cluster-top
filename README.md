# CLUSTER-TOP

The same as `top` but for multiple machines at the same time.

This use low-level C-code for efficiency for most operations. It only supports Linux (tested under Ubuntu)!

Output should be something like

```
+---------+---------------------------+-------+---------+-----------------+----------+
| Node    | RAM-Usage                 | Pid   | User    | Command         | CPU-Util |
+---------+---------------------------+-------+---------+-----------------+----------+
| node00  | 12349MiB / 32075MiB (37%) | 27133 | patwie  | sysbench        | 200%     |
|         |                           | 2546  | patwie  | audacity        | 8%       |
|         |                           | 3234  | patwie  | chrome          | 6%       |
|         |                           | 29589 | patwie  | plugin_host     | 4%       |
|         |                           | 30235 | patwie  | zsh             | 0%       |
|         |                           | 6147  | patwie  | chrome          | 0%       |
|         |                           | 25335 | patwie  | scdaemon        | 0%       |
|         |                           | 319   | patwie  | cluster-top-nod | 0%       |
|         |                           | 4873  | patwie  | chrome          | 0%       |
|         |                           | 29575 | patwie  | sublime_text    | 0%       |
+---------+---------------------------+-------+---------+-----------------+----------+
```

You might be interested as well in [cluster-smi](https://github.com/PatWie/cluster-smi) for GPUs.

Additional information are available, when using

```console
user@host $ cluster-top -h

Usage of cluster-top:
  -t  show time of events

```

## Monitoring Modes

This repository contains two versions: *cluster-top-local*, *cluster-top*.

### Local (cluster-top-local)

*cluster-top-local* is the same as *top*:


### Cluster (cluster-top)

*cluster-top* displays all information from *cluster-top-local* but for **multiple** machines at the same time.


On each machine you want to monitor you need to start *cluster-top-node*. They are sending information from the machine to a *cluster-top-router*, which further distributes these information to client (*cluster-top*) when requested.

You might be interested as well in [cluster-top](https://github.com/PatWie/cluster-top) for CPUS.

## Installation

### Requirements + Dependencies

- ZMQ (4.0.1)

Unfortunately, *ZMQ* can only be dynamically linked (`libzmq.so`) to this repository and you need to build it separately by

```bash
# compile ZMQ library for c++
cd /path/to/your_lib_folder
wget http:/files.patwie.com/mirror/zeromq-4.1.0-rc1.tar.gz
tar -xf zeromq-4.1.0-rc1.tar.gz
cd zeromq-4.1.0
./autogen.sh
./configure
./configure --prefix=/path/to/your_lib_folder/zeromq-4.1.0/dist
make
make install
```

Finally:

```
export PKG_CONFIG_PATH=/path/to/your_lib_folder/zeromq-4.1.0/dist/lib/pkgconfig/:$PKG_CONFIG_PATH
```

Edit the CFLAGS, LDFLAGS in file `nvvml/nvml.go` to match your setup.

### Compiling

You need to copy one config-file

```console
user@host $ cp config.example.go config.go
```

To obtain a portable small binary, I suggest to directly embed the configuration settings (ports, ip-addr) into the binary as compile-time constants. This way, the app is fully self-contained (excl. libzmq.so) and does not require any configuration-files. This can be done by editing `config.go`:

```go
...
c.RouterIp = "127.0.0.1"
c.Tick = 3
c.Timeout = 180
c.Ports.Nodes = "9080"
c.Ports.Clients = "9081"
...
```

Otherwise, you can specify the environment variable `CLUSTER_TOP_CONFIG_PATH` pointing to a yaml file (example in `cluster-top.example.yml`).

Then run

```bash
cd proc
go install
cd ..
make all
```


### Run

1. start `cluster-top-node` at different machines
2. start `cluster-top-router` at a specific machine (machine with ip-addr: `cluster_smi_router_ip`)
3. use `cluster-top` like `nvidia-top`

Make sure, the machines can communicate using the specifiec ports (e.g., `ufw allow 9080, 9081`)

### Use systemd

To ease the use of this app, I suggest to add the *cluster-top-node* into a systemd-service. An example config file can be found <a href="./docs/cluster-top-node.example.service">here</a>. The steps would be

```bash
# add new entry to systemd
sudo cp docs/cluster-top-node.example.service /etc/systemd/system/cluster-top-node.service
# edit the path to cluster-top-node
sudo nano /etc/systemd/system/cluster-top-node.service
# make sure you can start and stop the service (have a look at you cluster-top client)
sudo service cluster-top-node start
sudo service cluster-top-node stop
# register cluster-top-node to start on reboots
sudo systemctl enable cluster-top-node.service

# last, start the service
sudo service cluster-top-node start
```

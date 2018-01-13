# CLUSTER-TOP

The same as `top` but for multiple machines at the same time.

<p align="center"> <img src="./cluster-smi.jpg" width="100%"> </p>


Output should be something like

```
+---------+---------------------------+-------+-------+-----------------+----------+
| Node    | RAM-Usage                 | Pid   | User  | Command         | CPU-Util |
+---------+---------------------------+-------+-------+-----------------+----------+
| node00  | 12349MiB / 32075MiB (37%) | 27133 | user  | sysbench        | 200%     |
|         |                           | 2546  | user  | audacity        | 8%       |
|         |                           | 3234  | user  | chrome          | 6%       |
|         |                           | 29589 | user  | plugin_host     | 4%       |
|         |                           | 30235 | user  | zsh             | 0%       |
|         |                           | 6147  | user  | chrome          | 0%       |
|         |                           | 25335 | user  | scdaemon        | 0%       |
|         |                           | 319   | user  | cluster-top-nod | 0%       |
|         |                           | 4873  | user  | chrome          | 0%       |
|         |                           | 29575 | user  | sublime_text    | 0%       |
+---------+---------------------------+-------+-------+-----------------+----------+
```


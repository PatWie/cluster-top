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
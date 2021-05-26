# winCache

>  Everything can become cache

## [](#introduction)Introduction

一个用go语言书写的cachedb，操作类似redis. <br>
one cache using go.<br>


## [](#Description)Description
目前是通过一致性算法通过hash选择节点，<br>
Currently, nodes are selected through hash through consensus algorithm,<br>
再hash选择对应的group，<br>
Then hash select the corresponding group,<br>
然后再group中进行命中，<br>
Then hit in the group,<br>
而淘汰策略选用lru。<br>
The elimination strategy uses lru.<br>

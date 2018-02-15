# etcdcode
etcd相关功能代码
# 相关功能说明
## 1. etcdop
- etcdop功能类似与etcdctl相关的基本操作：get/getprefix/set/grantset(set并且设定一个超时)/watch
- 主要引用了myetcd

## 2. example
- 库的功能测试代码。现在有monitor/myrpc的相关功能测试代码

## 3. myetcd
- etcd的相关操作封装。这里创建etcd对象的时候，需要指定systemname（系统- 名称，即您所开发部件的系统），这是防止与其他系统混淆。
  选举，好像有问题（在多个etcd集群的时候）

## 4. monitor
- monitr主要用于服务发现。现在是指定需要监控的部件（进程名称）。
- 后期，应该会改成监控所有部件。但是，返回给外部的时候，进行过滤——需要监控哪个就哪个部件数据返回

## 5. myrpc
- 该库主要用于部件之间的通信。如果是内部部件，那么不需要证书；如果是外部部件，可以有证书。
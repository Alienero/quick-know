# 环境需求
MongoDB,etcd,Golang(>=1.4),redis

# 编译
```
go get github.com/Alienero/quick-know
go install  github.com/Alienero/quick-know/clients/cli
go install  github.com/Alienero/quick-know/config/shared
go install  github.com/Alienero/quick-know/comet
go install  github.com/Alienero/quick-know/comet
mv $GOPATH/src/github.com/Alienero/quick-know/config/shared/qk.conf $GOPATH/bin/qk.conf
```
启动MongoDB，etcd，redis，qk.conf中键入配置

```
#  自动配置
./shared -etcd=http://127.0.0.1:4001 -path=qk.conf
# 启动Comet节点
# 若etcd节点有多个，用逗号隔开.
./comet -etcd=http://127.0.0.1:4001 -rpc=127.0.0.1:8899 -tcp_listen=127.0.0.1:9001  -logtostderr=true
./comet -etcd=http://127.0.0.1:4001 -rpc=127.0.0.1:8999 -tcp_listen=127.0.0.1:9002  -logtostderr=true
# 启动Web节点
./web  -logtostderr=true -etcd=http://127.0.0.1:4001 -listen=127.0.0.1:9901
```

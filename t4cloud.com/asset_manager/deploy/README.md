# 部署步骤

## 准备部署工具
####（别准备了，我再这里放了两个已经生成好的）
``` base
在fabric路径下
make release
然后在
cd release/linux-amd64/bin/
就能找到构建好的所有工具

```

## 生成证书文件
``` bash
./cryptogen generate --config=./crypto-config.yaml
```

## 生成创世区块
``` bash
./configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block -channelID t4channel
```

## 生成通道的创世交易
**注意！此处的交易链ID和上方的系统链ID不能重名！**
**注意！因为我写JAVA的习惯，导致我使用了驼峰，不允许有大写！**
``` bash
./configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/t4cloud.tx -channelID t4channeltx
```

## 生成组织关于通道的锚节点（主节点）交易[这步骤可以不做]
``` bash
./configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org0MSPanchors.tx -channelID t4channel -asOrg Org0MSP
./configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID t4channel -asOrg Org1MSP
```

## 启动前查看docker有没有遗留，如果有的话需要清空
``` bash
停止所有容器
docker stop $(docker ps -a -q) 
remove删除所有容器
$docker  rm $(docker ps -a -q) 
docker container prune
docker system prune --all --force
docker rmi –f \[IMAGE ID\]
```

## 启动整个节点网络
``` bash
docker-compose up -d
```

## 查看节点是否启动正常
``` bash
docker logs [name]
```

## 启动正常，进入cli操作。配置节点通道等
``` bash
docker exec -it cli bash
```

## 查看节点通道
``` bash
peer channel list
```

## 创建通道
``` bash
peer channel create -o orderer.t4cloud.com:7050 -c t4channeltx -f /etc/hyperledger/config/t4cloud.tx
peer channel create -o orderer.t4cloud.com:7050 -c assetschannel -f /etc/hyperledger/config/assetschannel.tx
```

## 加入通道
``` bash
peer channel join -b t4channeltx.block
peer channel join -b assetschannel.block
```

## 设置主节点
``` bash
peer channel update -o orderer.t4cloud.com:7050 -c t4channeltx -f /etc/hyperledger/config/Org1MSPanchors.tx
```

## 链码安装
``` bash
peer chaincode install -n assets -v 1.0.0 -l golang -p github.com/chaincode/assetsExchange
```

## 链码实例化
``` bash
peer chaincode instantiate -o orderer.t4cloud.com:7050 -C assetschannel -n assets -l golang -v 1.0.0 -c '{"Args":["init"]}'
```

## 链码交互
``` bash
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userRegister", "user1", "user1"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["assetEnroll", "asset1", "asset1", "metadata", "user1"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userRegister", "user2", "user2"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["assetExchange", "user1", "asset1", "user2"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userDestroy", "user1"]}'
```

## 链码升级
``` bash
peer chaincode install -n assets -v 1.0.1 -l golang -p github.com/chaincode/assetsExchange
peer chaincode upgrade -C assetschannel -n assets -v 1.0.1 -c '{"Args":[""]}'
```



## 链码查询
``` bash
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryUser", "user1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAsset", "asset1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryUser", "user2"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1", "all"]}'
```

## 命令行模式的背书策略

EXPR(E[,E...])
EXPR = OR AND
E = EXPR(E[,E...])
MSP.ROLE
MSP 组织名 org0MSP org1MSP
ROLE admin member

OR('org0MSP.member','org1MSP.admin')

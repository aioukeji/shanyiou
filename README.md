## 运行方式

### 准备环境
首先按照[官方文档](https://hyperledger-fabric.readthedocs.io/en/release-2.0/prereqs.html)安装依赖。
Clone 代码:
```bash
git clone http://github.com/aioukeji/shanyiou
cd shanyiou
```   
之后获取 Fabric:
```bash
# 需要在 shanyiou 目录中运行
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s -- 2.0.1 1.4.6 0.4.18
cd fabric-samples
# 因为 Fabric Golang SDK 尚不支持 Fabric2.0 中的新 lifecycle
# 需要限制这个 feature
git apply ../fabric-samples.patch
cd test-network
# 启动 Fabric
./network.sh up createChannel -s couchdb
```

### 编译，启动
```bash
# 回到项目根目录
cd ../..
go build -mod vendor

# 启动 redis
brew services start redis # macos
# 或：sudo systemctl start redis # linux

./shanyiou
```

### 测试
```bash
# 回到项目根目录
cd example
yarn install
node client.js
```

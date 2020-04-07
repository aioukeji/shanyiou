## How to launch

### Setup Environment
First, install dependencies following [Fabric Docs](https://hyperledger-fabric.readthedocs.io/en/release-2.0/prereqs.html).   
Then, clone the code:
```bash
git clone http://github.com/aioukeji/shanyiou
cd shanyiou
```   
Then, install Fabric:
```bash
# run the following cmd in `shanyiou` folder
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s -- 2.0.1 1.4.6 0.4.18
cd fabric-samples
# Disable fabric chaincode lifecycle which in not support in Fabric Golang SDK yet
git apply ../fabric-samples.patch
cd test-network
# Launch Fabric network
./network.sh up createChannel -s couchdb
```

### Compile and launch
```bash
# return to project folder
cd ../..
go build -mod vendor

# Lanuch redis
brew services start redis # macos
# or: sudo systemctl start redis # linux

./shanyiou
```

### Test
```bash
# return to project folder
cd example
yarn install
node client.js
```

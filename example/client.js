const jayson = require("jayson/promise");
require("util").inspect.defaultOptions.depth = null;

const host = "127.0.0.1";
const url = `http://${host}:8910/jsonrpc`;
console.log("connecting", url);
const client = jayson.client.http(url);


async function listCharity() {
  try {
    console.log(await client.request("RpcStub.ListCharitySummary", [{}]));
    console.log(
      await client.request("RpcStub.ListCharityIncome", [
        { Offset: 0, Limit: 100 }
      ])
    );
    console.log(
      await client.request("RpcStub.ListCharityOutcome", [
        { Offset: 0, Limit: 100 }
      ])
    );
  } catch (e) {
    console.log("Exception:", e);
  }
}

async function testCharity() {
  try {
    console.log(await client.request("RpcStub.ListCharitySummary", [{}]));
    console.log(
      await client.request("RpcStub.ListCharityIncome", [
        { Offset: 100, Limit: 100 }
      ])
    );
    console.log(
      await client.request("RpcStub.BatchAddCharityIncome", [{
        Data: [
          { From: "成龙", Category: "资金", Amount: 10000, Detail: "收入项目" },
          { From: "李龙", Category: "口罩", Amount: 10000, Detail: "病毒防止" }
        ]
      }])
    );
    console.log(
      await client.request("RpcStub.AddCharityOutcome", [
        {
          To: "协和医院",
          Source: "s000000",
          Category: "资金",
          Amount: 1000,
          Detail: "支出项目"
        }
      ])
    );
    console.log(await client.request("RpcStub.ListCharitySummary", [{}]));
    console.log(
      await client.request("RpcStub.ListCharityIncome", [
        { Offset: 0, Limit: 100 }
      ])
    );
  } catch (e) {
    console.log("Exception:", e);
  }
}


async function main() {
  try {
    await testCharity();
    await listCharity();
  } catch (e) {
    console.log("Exception:", e);
  }
}
main();

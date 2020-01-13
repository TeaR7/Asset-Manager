package main

import (
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AssetManager struct{}

//------------------------------------------------------- 数据模型 -------------------------------------------------------

// 用户
type User struct {
	Id string `json:"id"`
	//用户名称
	Name string `json:"name"`
}

// 资产
type Asset struct {
	Id string `json:"id"`
	//资产名称
	Name string `json:"name"`
	//当前归属用户
	UserId string `json:"userId"`
	//可扩展的预留属性
	Metadata string `json:"metadata"`
}

// 资产变更记录
type AssetHistory struct {
	//资产ID
	AssetId string `json:"assetId"`
	// 变更前拥有者ID
	OriginUserId string `json:"originUserId"`
	// 变更后拥有者ID
	TargetUserId string `json:"targetUserId"`
}

func getKey(userId, keyType string) string {
	return fmt.Sprintf("%s_%s", keyType, userId)
}

//------------------------------------------------------- 用户相关 -------------------------------------------------------

// 用户注册
func userRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 2 {
		return shim.Error("参数数量异常")
	}
	//注意！由于区块链开发的特殊性，ID不能在链码内随机生成
	id := args[0]
	name := args[1]
	if name == "" || id == "" {
		return shim.Error("参数内容为空")
	}

	// 用户是否已注册
	if userBytes, err := stub.GetState(getKey(id, "user")); err == nil && len(userBytes) != 0 {
		return shim.Error("用户已存在（ID已存在）")
	}

	// 生成用户对象
	user := &User{
		Id:   id,
		Name: name,
	}

	// 序列化对象
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("json化异常！ %s", err))
	}

	if err := stub.PutState(getKey(id, "user"), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("用户生成失败 error %s", err))
	}

	// 成功返回
	return shim.Success(nil)
}

// 用户查询
func queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数验证
	if len(args) != 1 {
		return shim.Error("参数不足")
	}

	userId := args[0]
	if userId == "" {
		return shim.Error("用户ID为空")
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	userBytes, err := stub.GetState(getKey(userId, "user"))
	if err != nil || len(userBytes) == 0 {
		return shim.Error("用户不存在")
	}

	return shim.Success(userBytes)
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (c *AssetManager) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (c *AssetManager) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()

	switch funcName {
	case "userRegister":
		//用户注册
		return userRegister(stub, args)
	case "queryUser":
		//用户查询
		return queryUser(stub, args)
	default:
		return shim.Error(fmt.Sprintf("不支持的方法: %s", funcName))
	}

	// stub.SetEvent("name", []byte("data"))
}

func main() {
	err := shim.Start(new(AssetManager))
	if err != nil {
		fmt.Printf("Error starting AssertsExchange chaincode: %s", err)
	}
}

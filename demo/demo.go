package demo

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

/*
初始化sdk
输入：配置文件路径
输出：sdk对象
*/
func NewSDK(configPath string) (*fabsdk.FabricSDK, error) {
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create new SDK: %s")
	}
	return sdk, nil
}

/*
初始化资源管理客户端
输入：fabsdk、用户名、组织名
输出：资源管理客户端
*/
func newResClient(sdk *fabsdk.FabricSDK, userName, orgName string) (*resmgmt.Client, error) {
	clientProvider := sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	resMgmtClient, err := resmgmt.New(clientProvider)
	if err != nil {
		return nil, err
	}
	return resMgmtClient, err
}

/*
初始化账本客户端
输入：fabsdk、通道名、用户名、组织名
输出：资源管理客户端
*/
func newLedgerClient(sdk *fabsdk.FabricSDK, channelName, userName, orgName string) (*ledger.Client, error) {
	channelProvider := sdk.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return nil, err
	}
	return ledgerClient, nil
}

/*
查询指定的peer节点加入了哪些通道
输入：fsbsdk、peer名字、用户名、组织名
输出：包含所有通道名的字符串切片
*/
func QueryChannels(sdk *fabsdk.FabricSDK, peerName, userName, orgName string) ([]string, error) {
	resMgmtClient, err := newResClient(sdk, userName, orgName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create a resource management client:%s")
	}

	channelQueryResponse, err := resMgmtClient.QueryChannels(
		resmgmt.WithTargetEndpoints(peerName))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query channels: %s")
	}

	channels := []string{}
	for _, channel := range channelQueryResponse.Channels {
		channels = append(channels, channel.ChannelId)
	}
	return channels, nil
}

/*
查询指定通道的区块链信息，包括高度、当前区块哈希、前一个区块哈希
输入：fsbsdk、通道名、用户名、组织名
输出：通道的区块链信息
*/
func QueryBlockInfo(sdk *fabsdk.FabricSDK, channelName, userName, orgName string) (*fab.BlockchainInfoResponse, error) {
	ledgerClient, err := newLedgerClient(sdk, channelName, userName, orgName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create a ledger client:%s")
	}
	bci, err := ledgerClient.QueryInfo()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query for blockchain info: %s")
	}
	return bci, nil
}

/*
查询指定区块高度的区块信息
输入：fsbsdk、通道名、用户名、组织名、区块高度
输出：区块信息
*/
func QueryBlockByIndex(sdk *fabsdk.FabricSDK, channelName, userName, orgName string, index uint64) (*common.Block, error) {
	ledgerClient, err := newLedgerClient(sdk, channelName, userName, orgName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create a ledger client:%s")
	}
	block, err := ledgerClient.QueryBlock(index)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query a block by index:%s")
	}
	return block, nil
}

/*
查询指定区块哈希的区块信息
输入：fsbsdk、通道名、用户名、组织名、区块哈希
输出：区块信息
*/
func QueryBlockByHash(sdk *fabsdk.FabricSDK, channelName, userName, orgName string, hash []byte) (*common.Block, error) {
	ledgerClient, err := newLedgerClient(sdk, channelName, userName, orgName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create a ledger client:%s")
	}
	block, err := ledgerClient.QueryBlockByHash(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query a block by hash:%s")
	}
	return block, nil
}

/*
查询指定交易id的交易信息
输入：fsbsdk、通道名、用户名、组织名、交易id
输出：交易信息
*/
func QueryTxByID(sdk *fabsdk.FabricSDK, channelName, userName, orgName string, txID fab.TransactionID) (*peer.ProcessedTransaction, error) {
	ledgerClient, err := newLedgerClient(sdk, channelName, userName, orgName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create a ledger client:%s")
	}
	tx, err := ledgerClient.QueryTransaction(txID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query a block by transactionID:%s")
	}
	return tx, nil
}

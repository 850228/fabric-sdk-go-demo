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

// const (
// 	orgName  = "Org1"
// 	orgAdmin = "Admin"
// )

func NewSDK(configPath string) (*fabsdk.FabricSDK, error) {
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create new SDK: %s")
	}
	return sdk, nil
}

func newResClient(sdk *fabsdk.FabricSDK, userName, orgName string) (*resmgmt.Client, error) {
	clientProvider := sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	resMgmtClient, err := resmgmt.New(clientProvider)
	if err != nil {
		return nil, err
	}
	return resMgmtClient, err
}

func newLedgerClient(sdk *fabsdk.FabricSDK, channelName, userName, orgName string) (*ledger.Client, error) {
	channelProvider := sdk.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return nil, err
	}
	return ledgerClient, nil
}

// 查询指定peer节点加入了哪些通道
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

// Query blockchain info, including Height,CurrentBlockHash,PreviousBlockHash
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

func QueryBlockByID(sdk *fabsdk.FabricSDK, channelName, userName, orgName string, txID fab.TransactionID) (*peer.ProcessedTransaction, error) {
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

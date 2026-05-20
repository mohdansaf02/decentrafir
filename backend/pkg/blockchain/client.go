package blockchain

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractABI = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"string","name":"firId","type":"string"},{"internalType":"bytes32","name":"firHash","type":"bytes32"}],"name":"createFIR","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"firId","type":"string"},{"internalType":"enum FIRManagement.FIRStatus","name":"newStatus","type":"uint8"}],"name":"updateFIRStatus","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`

type Client struct {
	client   *ethclient.Client
	contract common.Address
	abi      abi.ABI
	key      *bind.TransactOpts
	chainID  *big.Int
	privKey  *ecdsa.PrivateKey
}

func NewClient(rpcURL, contractAddr, privateKeyHex string) (*Client, error) {
	if contractAddr == "" || privateKeyHex == "" {
		return &Client{}, nil // offline mode
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:   client,
		contract: common.HexToAddress(contractAddr),
		abi:      parsedABI,
		key:      auth,
		chainID:  chainID,
		privKey:  key,
	}, nil
}

func MetadataHash(data string) string {
	sum := sha256.Sum256([]byte(data))
	return "0x" + hex.EncodeToString(sum[:])
}

func (c *Client) CreateFIROnChain(ctx context.Context, firID, metadata string) (string, uint64, error) {
	if c.client == nil {
		return "0xmock" + firID, 0, nil
	}
	hash := common.HexToHash(strings.TrimPrefix(MetadataHash(metadata), "0x"))
	data, err := c.abi.Pack("createFIR", firID, hash)
	if err != nil {
		return "", 0, err
	}
	nonce, err := c.client.PendingNonceAt(ctx, c.key.From)
	if err != nil {
		return "", 0, err
	}
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", 0, err
	}
	tx := types.NewTransaction(nonce, c.contract, big.NewInt(0), 300000, gasPrice, data)
	signer := types.LatestSignerForChainID(c.chainID)
	signed, err := types.SignTx(tx, signer, c.privKey)
	if err != nil {
		return "", 0, err
	}
	if err := c.client.SendTransaction(ctx, signed); err != nil {
		return "", 0, err
	}
	receipt, err := bind.WaitMined(ctx, c.client, signed)
	if err != nil {
		return signed.Hash().Hex(), 0, err
	}
	return signed.Hash().Hex(), receipt.BlockNumber.Uint64(), nil
}

func (c *Client) UpdateStatusOnChain(ctx context.Context, firID string, status uint8) (string, error) {
	if c.client == nil {
		return "0xmock-update-" + firID, nil
	}
	data, err := c.abi.Pack("updateFIRStatus", firID, status)
	if err != nil {
		return "", err
	}
	nonce, err := c.client.PendingNonceAt(ctx, c.key.From)
	if err != nil {
		return "", err
	}
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	tx := types.NewTransaction(nonce, c.contract, big.NewInt(0), 300000, gasPrice, data)
	signer := types.LatestSignerForChainID(c.chainID)
	signed, err := types.SignTx(tx, signer, c.privKey)
	if err != nil {
		return "", err
	}
	if err := c.client.SendTransaction(ctx, signed); err != nil {
		return "", err
	}
	_, err = bind.WaitMined(ctx, c.client, signed)
	if err != nil {
		return signed.Hash().Hex(), err
	}
	return signed.Hash().Hex(), nil
}

func StatusToChain(status string) uint8 {
	switch status {
	case "under_review":
		return 1
	case "approved":
		return 2
	case "rejected":
		return 3
	case "closed":
		return 4
	default:
		return 0
	}
}

func (c *Client) Enabled() bool {
	return c.client != nil
}

func (c *Client) ContractAddress() string {
	if c.contract == (common.Address{}) {
		return ""
	}
	return c.contract.Hex()
}

// VerifyWalletSignature validates an Ethereum signed message (EIP-191 personal_sign).
func VerifyWalletSignature(address, message, signatureHex string) bool {
	sig := common.FromHex(signatureHex)
	if len(sig) != 65 {
		return false
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefix))
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false
	}
	recovered := crypto.PubkeyToAddress(*pubKey)
	return strings.EqualFold(recovered.Hex(), address)
}

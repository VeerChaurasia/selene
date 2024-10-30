package evm

import (
	"math/big"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewEnv(t *testing.T) {
	env := NewEnv()
	assert.NotNil(t, env, "NewEnv should not return nil")
}

func TestNewSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)

	assert.NotNil(t, sig)
	assert.Equal(t, uint8(1), sig.V)
	assert.Equal(t, r, sig.R)
	assert.Equal(t, s, sig.S)
}

func TestToRawSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)

	rawSig := sig.ToRawSignature()
	assert.Equal(t, 64, len(rawSig))
	assert.Equal(t, r.Bytes(), rawSig[32-len(r.Bytes()):32])
	assert.Equal(t, s.Bytes(), rawSig[64-len(s.Bytes()):])
}

func TestFromRawSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)
	rawSig := sig.ToRawSignature()

	reconstructedSig, err := FromRawSignature(rawSig, sig.V)
	assert.NoError(t, err)
	assert.Equal(t, sig.V, reconstructedSig.V)
	assert.Equal(t, r, reconstructedSig.R)
	assert.Equal(t, s, reconstructedSig.S)
}

func TestCfgEnvStruct(t *testing.T) {
	cfg := CfgEnv{
		ChainID: 1234,
		DisableBaseFee: true,
		MemoryLimit: 2048,
	}
	assert.Equal(t, uint64(1234), cfg.ChainID)
	assert.True(t, cfg.DisableBaseFee)
	assert.Equal(t, uint64(2048), cfg.MemoryLimit)
}

func TestTxKind(t *testing.T) {
	txKind := TxKind{Type: Call2}
	assert.Equal(t, Call2, txKind.Type)
	assert.Nil(t, txKind.Address)
}

func TestOptionalNonce(t *testing.T) {
	var nonceValue uint64 = 42
	nonce := OptionalNonce{Nonce: &nonceValue}
	assert.Equal(t, uint64(42), *nonce.Nonce)
	
	nonce = OptionalNonce{}
	assert.Nil(t, nonce.Nonce)
}

func TestAuthorizationStructs(t *testing.T) {
	chainID := ChainID(12345)
	auth := Authorization{
		ChainID: chainID,
		Address: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		Nonce: OptionalNonce{},
	}
	assert.Equal(t, chainID, auth.ChainID)
	assert.Equal(t, Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}, auth.Address)
	assert.Nil(t, auth.Nonce.Nonce)

	signedAuth := SignedAuthorization{Inner: auth}
	assert.Equal(t, chainID, signedAuth.Inner.ChainID)

	recoveredAuth := RecoveredAuthorization{Inner: auth}
	assert.Equal(t, chainID, recoveredAuth.Inner.ChainID)
	assert.Nil(t, recoveredAuth.Authority)
}

func TestBlobExcessGasAndPrice(t *testing.T) {
	blob := BlobExcessGasAndPrice{
		ExcessGas:   1000,
		BlobGasPrice: 2000,
	}
	assert.Equal(t, uint64(1000), blob.ExcessGas)
	assert.Equal(t, uint64(2000), blob.BlobGasPrice)
}

func TestBlockEnvStruct(t *testing.T) {
	blockEnv := BlockEnv{
		Number:    U256(big.NewInt(1234)),
		Timestamp: U256(big.NewInt(5678)),
		Coinbase: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		GasLimit: U256(big.NewInt(1234)),
		BaseFee: U256(big.NewInt(1234)),
		Difficulty: U256(big.NewInt(1234)),
		Prevrandao: &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
		BlobExcessGasAndPrice: &BlobExcessGasAndPrice{
			ExcessGas: 90,
			BlobGasPrice: 80,
		},
	}
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.Number)
	assert.Equal(t, U256(big.NewInt(5678)), blockEnv.Timestamp)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.GasLimit)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.BaseFee)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.Difficulty)
	assert.Equal(t, &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee}, blockEnv.Prevrandao)
	assert.Equal(t, &BlobExcessGasAndPrice{ExcessGas: 90, BlobGasPrice: 80}, blockEnv.BlobExcessGasAndPrice)
}

func TestAccessListItem(t *testing.T) {
	accessList := AccessListItem{
		Address: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		StorageKeys: []B256{{0xaa, 0xbb, 0xcc, 0xdd, 0xee}},
	}
	assert.Equal(t, Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}, accessList.Address)
	assert.Equal(t, []B256{{0xaa, 0xbb, 0xcc, 0xdd, 0xee}}, accessList.StorageKeys)
}


package evm

import (
	"encoding/json"
	"testing"
	"math/big"
	"github.com/stretchr/testify/assert"
)

func TestCallInputsNew(t *testing.T) {
	// Test cases for CallInputs.New
	tests := []struct {
		name     string
		txEnv    *TxEnv
		gasLimit uint64
		want     *CallInputs
	}{
		{
			name: "valid call transaction",
			txEnv: &TxEnv{
				TransactTo: TransactTo{
					Type:    Call2,
					Address: &Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				},
				Data:   Bytes{1, 2, 3},
				Caller: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				Value:  U256(big.NewInt(123)),
			},
			gasLimit: 1000,
			want: &CallInputs{
				Input:              Bytes{1, 2, 3},
				GasLimit:           1000,
				TargetAddress:      Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				BytecodeAddress:    Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				Caller:             Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				Value:              CallValue{ValueType: "transfer", Amount: U256(big.NewInt(123))},
				Scheme:             ICall,
				IsStatic:           false,
				IsEof:              false,
				ReturnMemoryOffset: Range{Start: 0, End: 0},
			},
		},
		{
			name: "non-call transaction",
			txEnv: &TxEnv{
				TransactTo: TransactTo{
					Type: Create2,
				},
			},
			gasLimit: 1000,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CallInputs{}
			got := ci.New(tt.txEnv, tt.gasLimit)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateInputsNew(t *testing.T) {
	tests := []struct {
		name     string
		txEnv    *TxEnv
		gasLimit uint64
		want     *CreateInputs
	}{
		{
			name: "valid create transaction",
			txEnv: &TxEnv{
				TransactTo: TransactTo{
					Type: Create2,
				},
				Data:   Bytes{1, 2, 3},
				Caller: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				Value:  U256(big.NewInt(123)),
			},
			gasLimit: 1000,
			want: &CreateInputs{
				Caller: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
				Scheme: CreateScheme{
					SchemeType: SchemeTypeCreate,
				},
				Value:    U256(big.NewInt(123)),
				InitCode: Bytes{1, 2, 3},
				GasLimit: 1000,
			},
		},
		{
			name: "non-create transaction",
			txEnv: &TxEnv{
				TransactTo: TransactTo{
					Type: Call2,
				},
			},
			gasLimit: 1000,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &CreateInputs{}
			got := ci.New(tt.txEnv, tt.gasLimit)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEOFCreateInputsNewTx(t *testing.T) {
	tx := &TxEnv{
		Caller: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		Value:  U256(big.NewInt(456)),
		Data:   Bytes{4, 5, 6},
	}
	gasLimit := uint64(2000)

	want := EOFCreateInputs{
		Caller:   Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		Value:    U256(big.NewInt(456)),
		GasLimit: 2000,
		Kind: EOFCreateKind{
			Kind: TxK,
			Data: Bytes{4, 5, 6},
		},
	}

	got := EOFCreateInputs{}.NewTx(tx, gasLimit)
	assert.Equal(t, want, got)
}

func TestCallValueJSON(t *testing.T) {
	tests := []struct {
		name      string
		callValue CallValue
		wantJSON  string
	}{
		{
			name: "transfer value",
			callValue: CallValue{
				ValueType: "transfer",
				Amount:    U256(big.NewInt(123)),
			},
			wantJSON: `{"value_type":"transfer","amount":123}`,
		},
		{
			name: "apparent value",
			callValue: CallValue{
				ValueType: "apparent",
				Amount:    U256(big.NewInt(456)),
			},
			wantJSON: `{"value_type":"apparent","amount":456}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.callValue)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(bytes))

			// Test unmarshaling
			var decoded CallValue
			err = json.Unmarshal(bytes, &decoded)
			assert.NoError(t, err)
			assert.Equal(t, tt.callValue, decoded)
		})
	}
}

func TestCallValueHelperFunctions(t *testing.T) {
	amount := U256(big.NewInt(789))
	
	transferValue := Transfer(amount)
	assert.Equal(t, CallValue{
		ValueType: "transfer",
		Amount:    amount,
	}, transferValue)

	apparentValue := Apparent(amount)
	assert.Equal(t, CallValue{
		ValueType: "apparent",
		Amount:    amount,
	}, apparentValue)
}
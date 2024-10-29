package evm

import (
    "encoding/json"
    "math/big"
    "reflect"
    "testing"
	
)

// Test Account struct
func TestAccount(t *testing.T) {
    t.Run("zero value initialization", func(t *testing.T) {
        account := Account{}
        
        if !reflect.DeepEqual(account.Info, AccountInfo{}) {
            t.Errorf("Expected empty AccountInfo, got %v", account.Info)
        }
        
        if account.Storage != nil {
            t.Errorf("Expected nil Storage, got %v", account.Storage)
        }
        
        if account.Status != 0 {
            t.Errorf("Expected Status 0, got %v", account.Status)
        }
    })

    t.Run("initialization with values", func(t *testing.T) {
        code := Bytecode{
            Kind:      LegacyRawKind,
            LegacyRaw: []byte{1, 2, 3},
        }
        
        info := AccountInfo{
            Balance:  big.NewInt(100),
            Nonce:    1,
            CodeHash: B256{1, 2, 3},
            Code:     &code,
        }
        
        storage := make(EvmStorage)
        storage[big.NewInt(1)] = EvmStorageSlot{
            OriginalValue: big.NewInt(10),
            PresentValue:  big.NewInt(20),
            IsCold:        true,
        }

        account := Account{
            Info:    info,
            Storage: storage,
            Status:  Loaded,
        }

        if !reflect.DeepEqual(account.Info, info) {
            t.Errorf("Expected Info %v, got %v", info, account.Info)
        }

        if !reflect.DeepEqual(account.Storage, storage) {
            t.Errorf("Expected Storage %v, got %v", storage, account.Storage)
        }

        if account.Status != Loaded {
            t.Errorf("Expected Status Loaded, got %v", account.Status)
        }
    })
}

// Test AccountInfo struct and its constructor
func TestAccountInfo(t *testing.T) {
    t.Run("NewAccountInfo constructor", func(t *testing.T) {
        balance := big.NewInt(100)
        nonce := uint64(1)
        codeHash := B256{1, 2, 3}
        code := Bytecode{
            Kind:      LegacyRawKind,
            LegacyRaw: []byte{1, 2, 3},
        }

        info := NewAccountInfo(balance, nonce, codeHash, code)

        if info.Balance.Cmp(balance) != 0 {
            t.Errorf("Expected Balance %v, got %v", balance, info.Balance)
        }
        if info.Nonce != nonce {
            t.Errorf("Expected Nonce %v, got %v", nonce, info.Nonce)
        }
        if !reflect.DeepEqual(info.CodeHash, codeHash) {
            t.Errorf("Expected CodeHash %v, got %v", codeHash, info.CodeHash)
        }
        if !reflect.DeepEqual(*info.Code, code) {
            t.Errorf("Expected Code %v, got %v", code, *info.Code)
        }
    })
}

// Test EvmStorage
func TestEvmStorage(t *testing.T) {
    t.Run("storage operations", func(t *testing.T) {
        storage := make(EvmStorage)
        
        // Test setting and getting values
        key := big.NewInt(1)
        slot := EvmStorageSlot{
            OriginalValue: big.NewInt(10),
            PresentValue:  big.NewInt(20),
            IsCold:        true,
        }
        
        storage[key] = slot
        
        retrieved, exists := storage[key]
        if !exists {
            t.Error("Expected storage slot to exist")
        }
        
        if !reflect.DeepEqual(retrieved, slot) {
            t.Errorf("Expected slot %v, got %v", slot, retrieved)
        }
    })
}

// Test AccountStatus constants and operations
func TestAccountStatus(t *testing.T) {
    tests := []struct {
        name   string
        status AccountStatus
        value  uint8
    }{
        {"Loaded", Loaded, 0b00000000},
        {"Created", Created, 0b00000001},
        {"SelfDestructed", SelfDestructed, 0b00000010},
        {"Touched", Touched, 0b00000100},
        {"LoadedAsNotExisting", LoadedAsNotExisting, 0b0001000},
        {"Cold", Cold, 0b0010000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if uint8(tt.status) != tt.value {
                t.Errorf("%s: expected %b, got %b", tt.name, tt.value, uint8(tt.status))
            }
        })
    }
}

// Test Bytecode
func TestBytecode(t *testing.T) {
    t.Run("legacy raw bytecode", func(t *testing.T) {
        code := Bytecode{
            Kind:      LegacyRawKind,
            LegacyRaw: []byte{1, 2, 3},
        }

        if code.Kind != LegacyRawKind {
            t.Errorf("Expected Kind LegacyRawKind, got %v", code.Kind)
        }
        if !reflect.DeepEqual(code.LegacyRaw, []byte{1, 2, 3}) {
            t.Errorf("Expected LegacyRaw [1 2 3], got %v", code.LegacyRaw)
        }
    })

    t.Run("legacy analyzed bytecode", func(t *testing.T) {
        analyzed := &LegacyAnalyzedBytecode{
            Bytecode:    []byte{1, 2, 3},
            OriginalLen: 3,
            JumpTable:   JumpTable{BitVector: &Bitvector{bits: []uint8{1}, size: 8}},
        }
        
        code := Bytecode{
            Kind:           LegacyAnalyzedKind,
            LegacyAnalyzed: analyzed,
        }

        if code.Kind != LegacyAnalyzedKind {
            t.Errorf("Expected Kind LegacyAnalyzedKind, got %v", code.Kind)
        }
        if !reflect.DeepEqual(code.LegacyAnalyzed, analyzed) {
            t.Errorf("Expected LegacyAnalyzed %v, got %v", analyzed, code.LegacyAnalyzed)
        }
    })
}

// Test EOF related structs
func TestEof(t *testing.T) {
    t.Run("eof header", func(t *testing.T) {
        header := EofHeader{
            TypesSize:         2,
            CodeSizes:         []uint16{100, 200},
            ContainerSizes:    []uint16{300, 400},
            DataSize:          500,
            SumCodeSizes:      300,
            SumContainerSizes: 700,
        }

        // Test JSON marshaling/unmarshaling
        data, err := json.Marshal(header)
        if err != nil {
            t.Fatalf("Failed to marshal EofHeader: %v", err)
        }

        var decoded EofHeader
        err = json.Unmarshal(data, &decoded)
        if err != nil {
            t.Fatalf("Failed to unmarshal EofHeader: %v", err)
        }

        if !reflect.DeepEqual(header, decoded) {
            t.Errorf("Expected header %+v, got %+v", header, decoded)
        }
    })

    t.Run("eof body", func(t *testing.T) {
        body := EofBody{
            TypesSection: []TypesSection{{
                Inputs:       1,
                Outputs:      2,
                MaxStackSize: 1024,
            }},
            CodeSection:      []Bytes{{1, 2, 3}},
            ContainerSection: []Bytes{{4, 5, 6}},
            DataSection:      Bytes{7, 8, 9},
            IsDataFilled:     true,
        }

        // Test JSON marshaling/unmarshaling
        data, err := json.Marshal(body)
        if err != nil {
            t.Fatalf("Failed to marshal EofBody: %v", err)
        }

        var decoded EofBody
        err = json.Unmarshal(data, &decoded)
        if err != nil {
            t.Fatalf("Failed to unmarshal EofBody: %v", err)
        }

        if !reflect.DeepEqual(body, decoded) {
            t.Errorf("Expected body %+v, got %+v", body, decoded)
        }
    })
}
func TestBitvector(t *testing.T) {
    t.Run("bitvector initialization", func(t *testing.T) {
        bv := Bitvector{
            bits: []uint8{0b10101010},
            size: 8,
        }

        if len(bv.bits) != 1 {
            t.Errorf("Expected bits length 1, got %d", len(bv.bits))
        }
        if bv.size != 8 {
            t.Errorf("Expected size 8, got %d", bv.size)
        }
    })
}

// Test JumpTable
func TestJumpTable(t *testing.T) {
    t.Run("jumptable initialization", func(t *testing.T) {
        jt := JumpTable{
            BitVector: &Bitvector{
                bits: []uint8{0b11110000},
                size: 8,
            },
        }

        if jt.BitVector == nil {
            t.Error("Expected non-nil BitVector")
        }
    })

    t.Run("jumptable sync.Once", func(t *testing.T) {
        jt := JumpTable{}
        initialized := false

        // Test lazy initialization
        jt.once.Do(func() {
            initialized = true
        })

        if !initialized {
            t.Error("sync.Once did not execute initialization")
        }
    })
}


func TestOpcode(t *testing.T) {
    t.Run("opcode structure", func(t *testing.T) {
        createdAddress := Address{Addr: [20]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}} 
        opcode := Opcode{
            InitCode: Eof{Header: EofHeader{
                    TypesSize: 1,
                    CodeSizes: []uint16{100},
                }, Body: EofBody{
                    TypesSection: []TypesSection{{
                        Inputs:       1,
                        Outputs:      1,
                        MaxStackSize: 100,
                    }},
                }, Raw: Bytes{1, 2, 3}},
            Input:          Bytes{4, 5, 6},
            CreatedAddress: createdAddress,
        }

        // Tests for opcode structure
        if len(opcode.Input) != 3 {
            t.Errorf("Expected Input length 3, got %d", len(opcode.Input))
        }

        if reflect.DeepEqual(opcode.InitCode, Eof{}) {
            t.Error("Expected non-empty InitCode")
        }

        if !reflect.DeepEqual(opcode.CreatedAddress, createdAddress) {
            t.Errorf("Expected CreatedAddress %v, got %v", createdAddress, opcode.CreatedAddress)
        }
    })
}



// Test TypesSection
func TestTypesSection(t *testing.T) {
    t.Run("types section values", func(t *testing.T) {
        ts := TypesSection{
            Inputs:       2,
            Outputs:      3,
            MaxStackSize: 1024,
        }

        if ts.Inputs != 2 {
            t.Errorf("Expected Inputs 2, got %d", ts.Inputs)
        }
        if ts.Outputs != 3 {
            t.Errorf("Expected Outputs 3, got %d", ts.Outputs)
        }
        if ts.MaxStackSize != 1024 {
            t.Errorf("Expected MaxStackSize 1024, got %d", ts.MaxStackSize)
        }
    })
}

//Not done- Test EOFCreateKind JSON marshaling/unmarshaling

func TestEvmStorageSlot(t *testing.T) {
    t.Run("storage slot values", func(t *testing.T) {
        slot := EvmStorageSlot{
            OriginalValue: big.NewInt(100),
            PresentValue:  big.NewInt(200),
            IsCold:        true,
        }

        if slot.OriginalValue.Cmp(big.NewInt(100)) != 0 {
            t.Errorf("Expected OriginalValue 100, got %v", slot.OriginalValue)
        }
        if slot.PresentValue.Cmp(big.NewInt(200)) != 0 {
            t.Errorf("Expected PresentValue 200, got %v", slot.PresentValue)
        }
        if !slot.IsCold {
            t.Error("Expected IsCold to be true")
        }
    })

    t.Run("storage slot modifications", func(t *testing.T) {
        slot := EvmStorageSlot{
            OriginalValue: big.NewInt(100),
            PresentValue:  big.NewInt(100),
            IsCold:        true,
        }

        // Modify present value
        slot.PresentValue = big.NewInt(150)
        
        if slot.PresentValue.Cmp(big.NewInt(150)) != 0 {
            t.Error("Failed to modify PresentValue")
        }
        if slot.OriginalValue.Cmp(big.NewInt(100)) != 0 {
            t.Error("OriginalValue should not change")
        }
    })
}

// Integration test for complex scenarios
func TestComplexScenarios(t *testing.T) {
    t.Run("account with storage modifications", func(t *testing.T) {
        // Create initial account state
        account := Account{
            Info: NewAccountInfo(
                big.NewInt(1000),
                1,
                B256{1},
                Bytecode{Kind: LegacyRawKind, LegacyRaw: []byte{1, 2, 3}},
            ),
            Storage: make(EvmStorage),
            Status:  Created,
        }

        // Add some storage
        key := big.NewInt(1)
        account.Storage[key] = EvmStorageSlot{
            OriginalValue: big.NewInt(100),
            PresentValue:  big.NewInt(100),
            IsCold:        true,
        }

        // Modify storage
        slot := account.Storage[key]
        slot.PresentValue = big.NewInt(200)
        slot.IsCold = false
        account.Storage[key] = slot

        // Verify all changes
        if account.Storage[key].PresentValue.Cmp(big.NewInt(200)) != 0 {
            t.Error("Storage modification failed")
        }
        if account.Storage[key].OriginalValue.Cmp(big.NewInt(100)) != 0 {
            t.Error("Original value should not change")
        }
        if account.Storage[key].IsCold {
            t.Error("IsCold should be false")
        }
    })

    t.Run("account status transitions", func(t *testing.T) {
        account := Account{Status: Created}

        // Test status transitions
        transitions := []AccountStatus{
            Touched,
            SelfDestructed,
            Cold,
        }

        for _, newStatus := range transitions {
            account.Status = newStatus
            if account.Status != newStatus {
                t.Errorf("Failed to transition to status %v", newStatus)
            }
        }
    })
}
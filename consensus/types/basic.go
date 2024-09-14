package types

import (
	"encoding/hex"
	"fmt"
)

// ByteList represents a list of bytes with a maximum capacity.
type ByteList struct {
	data    []byte
	maxSize int
}

// NewByteList creates a new ByteList with the given maximum size.
func NewByteList(maxSize int) *ByteList {
	return &ByteList{
		data:    make([]byte, 0, maxSize),
		maxSize: maxSize,
	}
}

// FromBytes creates a ByteList from a byte slice.
func FromBytes(data []byte, maxSize int) (*ByteList, error) {
	if len(data) > maxSize {
		return nil, fmt.Errorf("data exceeds maximum size of %d", maxSize)
	}
	return &ByteList{
		data:    append([]byte(nil), data...),
		maxSize: maxSize,
	}, nil
}

// Append adds bytes to the ByteList.
func (bl *ByteList) Append(data ...byte) error {
	if len(bl.data)+len(data) > bl.maxSize {
		return fmt.Errorf("appending would exceed maximum size of %d", bl.maxSize)
	}
	bl.data = append(bl.data, data...)
	return nil
}

// Bytes returns the underlying byte slice.
func (bl *ByteList) Bytes() []byte {
	return bl.data
}

// Len returns the current length of the ByteList.
func (bl *ByteList) Len() int {
	return len(bl.data)
}

// MaxSize returns the maximum size of the ByteList.
func (bl *ByteList) MaxSize() int {
	return bl.maxSize
}

// String returns a hexadecimal representation of the ByteList.
func (bl *ByteList) String() string {
	return "0x" + hex.EncodeToString(bl.data)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (bl *ByteList) MarshalText() ([]byte, error) {
	return []byte(bl.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (bl *ByteList) UnmarshalText(text []byte) error {
	if len(text) >= 2 && text[0] == '0' && (text[1] == 'x' || text[1] == 'X') {
		text = text[2:]
	}
	data, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(data) > bl.maxSize {
		return fmt.Errorf("data exceeds maximum size of %d", bl.maxSize)
	}
	bl.data = data
	return nil
}

// ByteVector represents a fixed-size vector of bytes.
type ByteVector[N int] struct {
	data [N]byte
}

// NewByteVector creates a new ByteVector from a slice of bytes.
func NewByteVector[N int](data []byte) (ByteVector[N], error) {
	if len(data) != N {
		return ByteVector[N]{}, fmt.Errorf("invalid data length: expected %d, got %d", N, len(data))
	}
	var bv ByteVector[N]
	copy(bv.data[:], data)
	return bv, nil
}

// Bytes returns the underlying byte array as a slice.
func (bv ByteVector[N]) Bytes() []byte {
	return bv.data[:]
}

// String returns a hexadecimal representation of the ByteVector.
func (bv ByteVector[N]) String() string {
	return "0x" + hex.EncodeToString(bv.data[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (bv ByteVector[N]) MarshalText() ([]byte, error) {
	return []byte(bv.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (bv *ByteVector[N]) UnmarshalText(text []byte) error {
	if len(text) >= 2 && text[0] == '0' && (text[1] == 'x' || text[1] == 'X') {
		text = text[2:]
	}
	data, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(data) != N {
		return fmt.Errorf("invalid data length: expected %d, got %d", N, len(data))
	}
	copy(bv.data[:], data)
	return nil
}

type Address = [20] byte;
type Bytes32 =[32] byte;
type LogsBloom =[256] byte;
type BLSPubKey =[48] byte;
type SignatureBytes= [96] byte;
type Transaction =ByteList;//1073741824
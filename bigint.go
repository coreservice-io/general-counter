package general_counter

import (
	"fmt"
	"math/big"

	"database/sql/driver"
)

type BigInt big.Int

func NewBigInt(value int64) *BigInt {
	x := big.NewInt(value)
	return (*BigInt)(x)
}

func NewBigIntFromString(s string) (*BigInt, bool) {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, ok
	}

	return (*BigInt)(x), ok
}

func (b *BigInt) Sign() int {
	return (*big.Int)(b).Sign()
}

func (b *BigInt) Neg() *BigInt {
	x := new(big.Int).Neg((*big.Int)(b))
	return (*BigInt)(x)
}

func (b *BigInt) Abs() *BigInt {
	x := new(big.Int).Abs((*big.Int)(b))
	return (*BigInt)(x)
}

func (b BigInt) String() string {
	return (*big.Int)(&b).String()
}

func (b BigInt) Value() (driver.Value, error) {
	return b.String(), nil
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}

	switch t := value.(type) {
	case int64:
		val, _ := value.(int64)
		(*big.Int)(b).SetInt64(val)
	case []uint8:
		str := string(value.([]uint8))
		_, ok := (*big.Int)(b).SetString(str, 10)
		if !ok {
			return fmt.Errorf("failed to load value from str: %v", str)
		}
	default:
		return fmt.Errorf("could not scan type %T into BigInt", t)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (b *BigInt) MarshalJSON() ([]byte, error) {
	return (*big.Int)(b).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *BigInt) UnmarshalJSON(text []byte) error {
	return (*big.Int)(b).UnmarshalJSON(text)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (b *BigInt) MarshalText() (text []byte, err error) {
	return (*big.Int)(b).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigInt) UnmarshalText(text []byte) error {
	return (*big.Int)(b).UnmarshalText(text)
}

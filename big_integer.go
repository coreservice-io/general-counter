package general_counter

import (
	"fmt"
	"math/big"

	"database/sql/driver"
)

type BigInteger big.Int

func NewBigInteger(val int64) *BigInteger {
	x := big.NewInt(val)
	return (*BigInteger)(x)
}

func NewBigIntegerFromString(val string) (*BigInteger, bool) {
	x, ok := new(big.Int).SetString(val, 10)
	if !ok {
		return nil, ok
	}

	return (*BigInteger)(x), ok
}

func NewBigIntegerFromBigInt(val *big.Int) *BigInteger {
	x := new(big.Int).Set(val)
	return (*BigInteger)(x)
}

func (b *BigInteger) BigInt() *big.Int {
	x := new(big.Int).Set((*big.Int)(b))
	return (*big.Int)(x)
}

func (b BigInteger) String() string {
	return (*big.Int)(&b).String()
}

func (b *BigInteger) Sign() int {
	return (*big.Int)(b).Sign()
}

func (b *BigInteger) Neg() *BigInteger {
	x := new(big.Int).Neg((*big.Int)(b))
	return (*BigInteger)(x)
}

func (b *BigInteger) Abs() *BigInteger {
	x := new(big.Int).Abs((*big.Int)(b))
	return (*BigInteger)(x)
}

func (b BigInteger) Value() (driver.Value, error) {
	return b.String(), nil
}

func (b *BigInteger) Scan(value interface{}) error {
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
		return fmt.Errorf("could not scan type %T into BigInteger", t)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (b *BigInteger) MarshalJSON() ([]byte, error) {
	return (*big.Int)(b).MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *BigInteger) UnmarshalJSON(text []byte) error {
	return (*big.Int)(b).UnmarshalJSON(text)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (b *BigInteger) MarshalText() (text []byte, err error) {
	return (*big.Int)(b).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigInteger) UnmarshalText(text []byte) error {
	return (*big.Int)(b).UnmarshalText(text)
}

// Package byteConverter is used to convert bytes to different datatypes.
package byteConverter

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type conversionType int

// conversion types that can be received from messages
const (
	conversionUint      conversionType = iota // 0
	conversionInt                             // 1
	conversionFloat                           // 2
	conversionString                          // 3
	conversionHexString                       // 4
	conversionBool                            // 5
)

// ByteConverter is used to convert slices of byte to other data types, that are
// then afterwards again converted to string for further usage.
type ByteConverter interface {
	ConvertSingleValue(payload []byte, conversion int) (string, error)
}

type byteConverter struct {
}

// New Creates a new MessageConverter
func New() ByteConverter {
	return new(byteConverter)
}

// Converts a single value that was received from a message to a string
// representation of the received value. The conversion int refers to the type of
// the value that has been received. The types that can be input are uint(0),
// int(1), float(2), string(3), hexString(4) or bool(5).
// Returns an error when an invalid type is chosen.
func (bc *byteConverter) ConvertSingleValue(payload []byte, conversion int) (string, error) {
	var result string
	var err error
	switch conversionType(conversion) {
	case conversionUint:
		result, err = convertUint(payload)
	case conversionInt:
		result, err = convertInt(payload)
	case conversionFloat:
		result, err = convertFloat(payload)
	case conversionString:
		result, err = convertString(payload)
	case conversionHexString:
		result, err = convertHexString(payload)
	case conversionBool:
		result, err = convertBool(payload)
	default:
		err = errors.New("Invalid Conversion type")
		return "", err
	}
	return result, err
}

// Conversion method for the uint type, which type of uint(8, 16, 32, 64) will
// be determined by the length of the byte slice. Returns a string representation
// of the resulting uint.
func convertUint(payload []byte) (string, error) {
	var result string
	var err error
	switch len(payload) {

	case 1: //uint8
		return strconv.FormatUint(uint64(payload[0]), 10), nil
	case 2: //uint16
		return strconv.FormatUint(uint64(binary.BigEndian.Uint16(payload)), 10), nil
	case 4: //uint32
		return strconv.FormatUint(uint64(binary.BigEndian.Uint32(payload)), 10), nil
	case 8: //uint64
		return strconv.FormatUint(binary.BigEndian.Uint64(payload), 10), nil
	default:
		result = ""
		err = errors.New("illegal length of payload for a uint type")
	}
	return result, err
}

// Conversion method for the int type, which type of int(8, 16, 32, 64) will
// be determined by the length of the byte slice. Returns a string representation
// of the resulting int.
func convertInt(payload []byte) (string, error) {

	var err error
	buf := bytes.NewReader(payload)
	switch len(payload) {
	case 1: //int8
		var integer int8
		err = binary.Read(buf, binary.BigEndian, &integer)
		return fmt.Sprintf("%d", integer), err
	case 2: //int16
		var integer int16
		err = binary.Read(buf, binary.BigEndian, &integer)
		return fmt.Sprintf("%d", integer), err
	case 4: //int32
		var integer int32
		err = binary.Read(buf, binary.BigEndian, &integer)
		return fmt.Sprintf("%d", integer), err
	case 8: //int64
		var integer int64
		err = binary.Read(buf, binary.BigEndian, &integer)
		return fmt.Sprintf("%d", integer), err
	default:
		err = errors.New("illegal length of payload for an int type")
		return "", err
	}
}

// Conversion method for the float type, which type of float(32, 64) will
// be determined by the length of the byte slice. Returns a string representation
// of the resulting float.
func convertFloat(payload []byte) (string, error) {
	var err error
	buf := bytes.NewReader(payload)
	switch len(payload) {
	case 4: //float32
		var float float32
		err = binary.Read(buf, binary.BigEndian, &float)
		return fmt.Sprintf("%f", float), err
	case 8: //float64
		var float float64
		err = binary.Read(buf, binary.BigEndian, &float)
		return fmt.Sprintf("%f", float), err
	default:
		err = errors.New("illegal length of payload for a float type")
		return "", err
	}
}

// Conversion method for the string type.
// A byte is replaced with the corresponding value in the ascii table.
func convertString(payload []byte) (string, error) {
	return string(payload), nil
}

// Conversion method for hex string, bytes that are being put in are converted
// to hex and returned:
// 0xAB is converted to AB
// Returned values are upper case.
func convertHexString(payload []byte) (string, error) {
	return strings.ToUpper(hex.EncodeToString(payload)), nil
}

// Conversion method for the bool type.
// It is expected that the length of payload is never greater than  one.
// Result will be either "true" or "false" as a string.
func convertBool(payload []byte) (string, error) {
	if len(payload) > 1 {
		return "", errors.New("invalid payload length")
	}
	if payload[0] > 0 {
		return "true", nil
	}
	return "false", nil
}

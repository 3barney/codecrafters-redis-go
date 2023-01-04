package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	SimpleString ValueType = '+'
	BulkString   ValueType = '$'
	Array        ValueType = '*'
)

type ValueType byte

type Value struct {
	valueType ValueType
	bytes     []byte
	array     []Value
}

// String converts Value to a string.
//
// If Value cannot be converted, an empty string is returned.
func (v Value) String() string {
	if v.valueType == BulkString || v.valueType == SimpleString {
		return string(v.bytes)
	}

	return ""
}

// Array converts Value to an array.
//
// If Value cannot be converted, an empty array is returned.
func (v Value) Array() []Value {
	if v.valueType == Array {
		return v.array
	}

	return []Value{}
}

func DecodeRESP(bufferDataStream *bufio.Reader) (Value, error) {
	byteData, err := bufferDataStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(byteData) {
	case "+":
		return DecodeSimpleString(bufferDataStream)
	case "$":
		return decodeBulkString(bufferDataStream)
	case "*":
		return DecodeArray(bufferDataStream)
	}
	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(byteData))
}

func DecodeArray(data *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilEnd(data)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(data)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		valueType: Array,
		array:     array,
	}, nil
}

func decodeBulkString(data *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilEnd(data)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(data, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		valueType: BulkString,
		bytes:     readBytes[:count],
	}, nil

}

func DecodeSimpleString(data *bufio.Reader) (Value, error) {
	readBytes, err := readUntilEnd(data)
	if err != nil {
		return Value{}, err
	}

	return Value{
		valueType: SimpleString,
		bytes:     readBytes,
	}, nil
}

func readUntilEnd(data *bufio.Reader) ([]byte, error) {
	byteData := []byte{}

	for {
		b, err := data.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		byteData = append(byteData, b...)
		if len(byteData) >= 2 && byteData[len(byteData)-2] == '\r' {
			break
		}
	}
	return byteData[:len(byteData)-2], nil

}

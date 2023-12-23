// Copyright 2019-2023 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package utils provides some utilities which might be useful specifically for GTP(or other telco protocols).
package utils

import (
	"encoding/binary"
	"encoding/hex"
)

// StrToSwappedBytes returns swapped bits from a byte.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is the hex character(0-f) to fill the last digit when
// handling a odd number. "f" is used In most cases.
func StrToSwappedBytes(s, filler string) ([]byte, error) {
	var raw []byte
	var err error
	if len(s)%2 == 0 {
		raw, err = hex.DecodeString(s)
	} else {
		raw, err = hex.DecodeString(s + filler)
	}
	if err != nil {
		return nil, err
	}

	return swap(raw), nil
}

// SwappedBytesToStr decodes raw swapped bytes into string.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is to decide whether to cut the last digit or not.
func SwappedBytesToStr(raw []byte, cutLastDigit bool) string {
	if len(raw) == 0 {
		return ""
	}

	s := hex.EncodeToString(swap(raw))
	if cutLastDigit {
		s = s[:len(s)-1]
	}

	return s
}

func swap(raw []byte) []byte {
	swapped := make([]byte, len(raw))
	for n := range raw {
		t := ((raw[n] >> 4) & 0xf) + ((raw[n] << 4) & 0xf0)
		swapped[n] = t
	}
	return swapped
}

// Uint24To32 converts 24bits-length []byte value into the uint32 with 8bits of zeros as prefix.
// This function is used for the fields with 3 octets.
func Uint24To32(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// Uint32To24 converts the uint32 value into 24bits-length []byte. The values in 25-32 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint32To24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// Uint40To64 converts 40bits-length []byte value into the uint64 with 8bits of zeros as prefix.
// This function is used for the fields with 3 octets.
func Uint40To64(b []byte) uint64 {
	if len(b) != 5 {
		return 0
	}
	return uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4])
}

// Uint64To40 converts the uint64 value into 40bits-length []byte. The values in 25-64 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint64To40(n uint64) []byte {
	return []byte{uint8(n >> 32), uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// EncodePLMN encodes MCC and MNC as BCD-encoded bytes.
func EncodePLMN(mcc, mnc string) ([]byte, error) {
	c, err := StrToSwappedBytes(mcc, "f")
	if err != nil {
		return nil, err
	}
	n, err := StrToSwappedBytes(mnc, "f")
	if err != nil {
		return nil, err
	}

	// 2-digit
	b := make([]byte, 3)
	if len(mnc) == 2 {
		b = append(c, n...)
		return b, nil
	}

	// 3-digit
	b[0] = c[0]
	b[1] = (c[1] & 0x0f) | (n[1] << 4 & 0xf0)
	b[2] = n[0]

	return b, nil
}

var usePLMNMemoization bool
var memoizedMNCs [0xFFFF + 1]string
var memoizedMCCs [0xFFFF + 1]string

// DecodeMCC decodes BCD-encoded MCC as it occurs in CGI/SAI/RAI.
func DecodeMNC(b []byte) string {
	if usePLMNMemoization {
		return memoizedMNCs[binary.BigEndian.Uint16(b)]
	}
	raw := hex.EncodeToString(b)
	if raw[0] == 'f' {
		return string([]byte{raw[3], raw[2]})
	}
	return string([]byte{raw[3], raw[2], raw[0]})
}

// DecodeMNC decodes BCD-encoded MNC as it occurs in CGI/SAI/RAI
func DecodeMCC(b []byte) string {
	if usePLMNMemoization {
		return memoizedMCCs[binary.BigEndian.Uint16(b)]
	}
	raw := hex.EncodeToString(b)
	return string([]byte{raw[1], raw[0], raw[3]})
}

// DecodePLMN decodes BCD-encoded bytes into MCC and MNC.
func DecodePLMN(b []byte) (mcc, mnc string, err error) {
	mcc = DecodeMCC(b[0:2])
	mnc = DecodeMNC(b[1:3])
	return
}

// ParseECI decodes ECI uint32 into e-NodeB ID and Cell ID.
func ParseECI(eci uint32) (enbID uint32, cellID uint8, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, eci)
	cellID = buf[3]
	enbID = binary.BigEndian.Uint32([]byte{0, buf[0], buf[1], buf[2]})
	return
}

func EnablePLMNMemoization() {
	var i uint16
	for i = 0; i < 0xFFFF; i++ {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, i)
		memoizedMCCs[i] = DecodeMCC(buf)
		memoizedMNCs[i] = DecodeMNC(buf)
	}
	i = 0xFFFF
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	memoizedMCCs[i] = DecodeMCC(buf)
	memoizedMNCs[i] = DecodeMNC(buf)
	usePLMNMemoization = true
}

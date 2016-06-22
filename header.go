package yar

import (
	"bytes"
	"encoding/binary"
)

const (
	ProtocolLength = 82
	PackagerLength = 8
)

type ErrorType int

const (
	MagicNumber = 0x80DFEC60
)

const (
	ERR_OKEY           ErrorType = 0x0
	ERR_PACKAGER       ErrorType = 0x1
	ERR_PROTOCOL       ErrorType = 0x2
	ERR_REQUEST        ErrorType = 0x4
	ERR_OUTPUT         ErrorType = 0x8
	ERR_TRANSPORT      ErrorType = 0x10
	ERR_FORBIDDEN      ErrorType = 0x20
	ERR_EXCEPTION      ErrorType = 0x40
	ERR_EMPTY_RESPONSE ErrorType = 0x80
)

type Header struct {
	Id          uint32
	Version     uint16
	MagicNumber uint32
	Reserved    uint32
	Provider    [28]byte
	Encrypt     uint32
	Token       [32]byte
	BodyLength  uint32
	Packager    [8]byte
}

func NewHeader() *Header {

	proto := new(Header)
	proto.MagicNumber = MagicNumber
	return proto

}

func NewHeaderWithBytes(payload *bytes.Buffer) *Header {

	p := NewHeader()
	p.Init(payload)
	return p
}

func (self *Header) Init(payload *bytes.Buffer) bool {

	binary.Read(payload, binary.BigEndian, &self.Id)
	binary.Read(payload, binary.BigEndian, &self.Version)
	binary.Read(payload, binary.BigEndian, &self.MagicNumber)
	binary.Read(payload, binary.BigEndian, &self.Reserved)
	binary.Read(payload, binary.BigEndian, &self.Provider)
	binary.Read(payload, binary.BigEndian, &self.Encrypt)
	binary.Read(payload, binary.BigEndian, &self.Token)
	binary.Read(payload, binary.BigEndian, &self.BodyLength)
	binary.Read(payload, binary.BigEndian, &self.Packager)
	return true
}

func (self *Header) Bytes() *bytes.Buffer {

	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, self)

	if err != nil {
		return nil
	}
	return buffer
}

package yar

import (
	"encoding/binary"
	"bytes"
)


type Protocol struct {
	
	id 			uint32
	version		uint16
	magic_num	uint32
	reserved   	uint32
	provider 	[32]byte
	token      	[32]byte
	body_len 	uint32 		
	
}

func ProtocolNew() *Protocol {
	
	proto := new(Protocol)
	
	return proto
	
}


func (self *Protocol)Init(payload *bytes.Buffer) bool {
	
	binary.Read(payload,binary.LittleEndian,&self.id)
	binary.Read(payload,binary.LittleEndian,&self.version)
	binary.Read(payload,binary.LittleEndian,&self.magic_num)
	binary.Read(payload,binary.LittleEndian,&self.reserved)
	binary.Read(payload,binary.LittleEndian,&self.provider)
	binary.Read(payload,binary.LittleEndian,&self.token)
	binary.Read(payload,binary.LittleEndian,&self.body_len)
		
	return true
}


func (self *Protocol) Bytes() *bytes.Buffer {
	
	buffer := new(bytes.Buffer)
	
	err := binary.Write(buffer, binary.LittleEndian, self)
	
	if err != nil {
		return nil
	}
	
	return buffer
	
}


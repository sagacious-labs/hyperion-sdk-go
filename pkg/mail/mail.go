// mail package contains helper to deal with hyperion "mails"
package mail

import (
	"encoding/binary"
	"io"
)

// MailType is an alias of uint8 which denotes type of Hyperion mail
// object
//
// This type is crucial to inform hyperion whether the mail is supposed
// to be treated as log or data
type MailType = uint8

const (
	// LOG - inidicates mail is of type log
	LOG MailType = 0

	// DATA - indicates mail is of type data
	DATA MailType = 1

	// mailTypeSize is the the size (in bytes) of the mail type
	// field in the Mail object
	mailTypeSize uint = 1

	// mailPayloadSize is the the size (in bytes) of the mail's
	// payload size's size field in the Mail object
	mailPayloadSize uint = 8
)

// Mail struct represents the data format which is supported
// by hyperion for communication
//
// Structure of the data (TLV based)
//   --------------------------------------
//  | typ (u8) | size (u64) | data (bytes) |
//   --------------------------------------
type Mail struct {
	// typ indicates the type of the mail
	typ MailType

	// size is the size of the payload that is attached
	// in this frame
	size uint64

	// data is any data that needs to be transported to
	// or from Hyperion
	data []byte
}

// New returns an instance of Mail struct
func New(typ MailType, data []byte) Mail {
	return Mail{
		typ:  typ,
		data: data,
		size: uint64(len(data)),
	}
}

// Encode converts the mail object into a byte slice
func (m Mail) Encode() (byt []byte) {
	// Encode type of the mail
	byt = append(byt, byte(m.typ))

	// Encode size of the data
	sbyt := make([]byte, 8)
	binary.LittleEndian.PutUint64(sbyt, m.size)
	byt = append(byt, sbyt...)

	// Encode the data
	byt = append(byt, m.data...)

	return
}

// EncodeTo takes a IO writer and writes the encoded
// mail struct into the given writer
func (m Mail) EncodeTo(dst io.Writer) error {
	data := m.Encode()
	written := 0

	for {
		n, err := dst.Write(data[written:])
		if err != nil {
			return err
		}

		written += n
		if written == len(data) {
			break
		}
	}

	return nil
}

// Decode takes in byte and decodes it into
// the current mail object
func (m *Mail) Decode(byt []byte) {
	// Get type of the mail
	m.typ = MailType(byt[0])

	// Get size of the payload
	size := byt[mailTypeSize : mailTypeSize+mailPayloadSize]
	m.size = binary.LittleEndian.Uint64(size)

	// Store the payload
	m.data = byt[mailTypeSize+mailPayloadSize : m.size]
}

// DecodeFrom takes a reader and will read the bytes into the mail object
func (m *Mail) DecodeFrom(src io.Reader) error {
	buffer := make([]byte, 128)
	data := []byte{}

	for {
		n, err := src.Read(buffer)
		if err != nil {
			return err
		}

		if n == 0 {
			break
		}

		data = append(data, buffer[:n]...)

		m.typ = MailType(data[0])

		if n < int(mailTypeSize)+int(mailPayloadSize) {
			continue
		}
		m.size = binary.LittleEndian.Uint64(data[mailTypeSize : mailTypeSize+mailPayloadSize])

		if n < int(m.size)+int(mailTypeSize)+int(mailPayloadSize) {
			continue
		}
		m.data = append(m.data, data[mailTypeSize+mailPayloadSize:m.size]...)

		break
	}

	return nil
}

// GetType returns type of the mail
func (m *Mail) GetType() MailType {
	return m.typ
}

// SetType sets type of the mail
func (m *Mail) SetType(typ MailType) {
	m.typ = typ
}

// GetType returns type of the mail
func (m *Mail) GetSize() MailType {
	return m.typ
}

// GetData returns data of the mail
func (m *Mail) GetData() []byte {
	return m.data
}

// SetData sets type of the mail
//
// This method also updates the size attribute
// hence this method should be used to update
// the mail data
func (m *Mail) SetData(data []byte) {
	m.size = uint64(len(data))
	m.data = data
}

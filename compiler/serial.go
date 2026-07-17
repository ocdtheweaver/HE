// Package compiler: .hbc binary format serialiser/deserialiser.
//
// Format (little-endian, length-prefixed):
//   Header:  "HE\x00\x01"  (magic + version byte)
//   Chunk:   see encodeChunk
//
// Chunk encoding:
//   name      string  (uint16 len + utf8)
//   constants uint16 count; each: tag byte + payload
//   names     uint16 count; each: string
//   subchunks uint16 count; each: MethodChunk
//   instrs    uint32 count; each: Instruction
//
// Instruction encoding:
//   opcode  byte
//   operand byte (tag) + variable payload

package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

const (
	hbcMagic   = "HE"
	hbcVersion = byte(1)
)

// ── Encode ────────────────────────────────────────────────────────────────────

func EncodeChunk(c *Chunk) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(hbcMagic)
	buf.WriteByte(hbcVersion)
	if err := encodeChunk(&buf, c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeChunk(buf *bytes.Buffer, c *Chunk) error {
	writeString(buf, c.Name)

	// Constants
	writeU16(buf, uint16(len(c.Constants)))
	for _, cv := range c.Constants {
		if err := encodeConst(buf, cv); err != nil {
			return err
		}
	}

	// Names
	writeU16(buf, uint16(len(c.Names)))
	for _, n := range c.Names {
		writeString(buf, n)
	}

	// SubChunks (method bodies)
	writeU16(buf, uint16(len(c.SubChunks)))
	for _, mc := range c.SubChunks {
		if err := encodeMethodChunk(buf, mc); err != nil {
			return err
		}
	}

	// Instructions
	writeU32(buf, uint32(len(c.Instructions)))
	for _, ins := range c.Instructions {
		if err := encodeInstruction(buf, ins); err != nil {
			return err
		}
	}
	return nil
}

func encodeMethodChunk(buf *bytes.Buffer, mc *MethodChunk) error {
	writeString(buf, mc.Name)
	writeString(buf, mc.ProtectTag) // "" for unprotected
	writeU16(buf, uint16(len(mc.Params)))
	for _, p := range mc.Params {
		writeString(buf, p)
	}
	writeU32(buf, uint32(len(mc.Instructions)))
	for _, ins := range mc.Instructions {
		if err := encodeInstruction(buf, ins); err != nil {
			return err
		}
	}
	return nil
}

func encodeConst(buf *bytes.Buffer, v interface{}) error {
	switch cv := v.(type) {
	case float64:
		buf.WriteByte('n')
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(cv))
		buf.Write(b)
	case string:
		buf.WriteByte('s')
		writeString(buf, cv)
	case bool:
		buf.WriteByte('b')
		if cv {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}
	case nil:
		buf.WriteByte('0')
	default:
		return fmt.Errorf("unserialisable constant type %T", v)
	}
	return nil
}

func encodeInstruction(buf *bytes.Buffer, ins Instruction) error {
	buf.WriteByte(byte(ins.Op))
	return encodeOperand(buf, ins.Operand)
}

func encodeOperand(buf *bytes.Buffer, op interface{}) error {
	switch v := op.(type) {
	case nil:
		buf.WriteByte('0')
	case int:
		buf.WriteByte('i')
		writeI32(buf, int32(v))
	case float64:
		buf.WriteByte('n')
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(v))
		buf.Write(b)
	case bool:
		buf.WriteByte('b')
		if v {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}
	case string:
		buf.WriteByte('s')
		writeString(buf, v)
	case [2]int:
		buf.WriteByte('2')
		writeI32(buf, int32(v[0]))
		writeI32(buf, int32(v[1]))
	case [3]int:
		buf.WriteByte('3')
		writeI32(buf, int32(v[0]))
		writeI32(buf, int32(v[1]))
		writeI32(buf, int32(v[2]))
	default:
		buf.WriteByte('0') // unknown → nil
	}
	return nil
}

// ── Decode ────────────────────────────────────────────────────────────────────

func DecodeChunk(data []byte) (*Chunk, error) {
	r := bytes.NewReader(data)

	// Header
	magic := make([]byte, 2)
	if _, err := r.Read(magic); err != nil {
		return nil, fmt.Errorf("hbc: can't read magic: %v", err)
	}
	if string(magic) != hbcMagic {
		return nil, fmt.Errorf("hbc: not a .hbc file (wrong magic)")
	}
	ver, err := r.ReadByte()
	if err != nil || ver != hbcVersion {
		return nil, fmt.Errorf("hbc: unsupported version %d", ver)
	}

	return decodeChunk(r)
}

func decodeChunk(r *bytes.Reader) (*Chunk, error) {
	c := &Chunk{}
	var err error

	c.Name, err = readString(r)
	if err != nil {
		return nil, err
	}

	// Constants
	nConst, _ := readU16(r)
	c.Constants = make([]interface{}, nConst)
	for i := range c.Constants {
		c.Constants[i], err = decodeConst(r)
		if err != nil {
			return nil, err
		}
	}

	// Names
	nNames, _ := readU16(r)
	c.Names = make([]string, nNames)
	for i := range c.Names {
		c.Names[i], err = readString(r)
		if err != nil {
			return nil, err
		}
	}

	// SubChunks
	nSub, _ := readU16(r)
	c.SubChunks = make([]*MethodChunk, nSub)
	for i := range c.SubChunks {
		c.SubChunks[i], err = decodeMethodChunk(r)
		if err != nil {
			return nil, err
		}
	}

	// Instructions
	nIns, _ := readU32(r)
	c.Instructions = make([]Instruction, nIns)
	for i := range c.Instructions {
		c.Instructions[i], err = decodeInstruction(r)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func decodeMethodChunk(r *bytes.Reader) (*MethodChunk, error) {
	mc := &MethodChunk{}
	var err error
	mc.Name, err = readString(r)
	if err != nil {
		return nil, err
	}
	mc.ProtectTag, err = readString(r)
	if err != nil {
		return nil, err
	}
	nParams, _ := readU16(r)
	mc.Params = make([]string, nParams)
	for i := range mc.Params {
		mc.Params[i], err = readString(r)
		if err != nil {
			return nil, err
		}
	}
	nIns, _ := readU32(r)
	mc.Instructions = make([]Instruction, nIns)
	for i := range mc.Instructions {
		mc.Instructions[i], err = decodeInstruction(r)
		if err != nil {
			return nil, err
		}
	}
	return mc, nil
}

func decodeConst(r *bytes.Reader) (interface{}, error) {
	tag, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch tag {
	case 'n':
		b := make([]byte, 8)
		r.Read(b)
		bits := binary.LittleEndian.Uint64(b)
		return math.Float64frombits(bits), nil
	case 's':
		return readString(r)
	case 'b':
		b, _ := r.ReadByte()
		return b != 0, nil
	case '0':
		return nil, nil
	}
	return nil, fmt.Errorf("unknown constant tag %q", tag)
}

func decodeInstruction(r *bytes.Reader) (Instruction, error) {
	opByte, err := r.ReadByte()
	if err != nil {
		return Instruction{}, err
	}
	op := Opcode(opByte)
	operand, err := decodeOperand(r)
	if err != nil {
		return Instruction{}, err
	}
	return Instruction{Op: op, Operand: operand}, nil
}

func decodeOperand(r *bytes.Reader) (interface{}, error) {
	tag, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch tag {
	case '0':
		return nil, nil
	case 'i':
		v, _ := readI32(r)
		return int(v), nil
	case 'n':
		b := make([]byte, 8)
		r.Read(b)
		bits := binary.LittleEndian.Uint64(b)
		return math.Float64frombits(bits), nil
	case 'b':
		b, _ := r.ReadByte()
		return b != 0, nil
	case 's':
		return readString(r)
	case '2':
		a, _ := readI32(r)
		b, _ := readI32(r)
		return [2]int{int(a), int(b)}, nil
	case '3':
		a, _ := readI32(r)
		b, _ := readI32(r)
		c, _ := readI32(r)
		return [3]int{int(a), int(b), int(c)}, nil
	}
	return nil, fmt.Errorf("unknown operand tag %q", tag)
}

// ── Wire helpers ──────────────────────────────────────────────────────────────

func writeU16(buf *bytes.Buffer, v uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	buf.Write(b)
}

func writeU32(buf *bytes.Buffer, v uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	buf.Write(b)
}

func writeI32(buf *bytes.Buffer, v int32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	buf.Write(b)
}

func writeString(buf *bytes.Buffer, s string) {
	writeU16(buf, uint16(len(s)))
	buf.WriteString(s)
}

func readU16(r *bytes.Reader) (uint16, error) {
	b := make([]byte, 2)
	_, err := r.Read(b)
	return binary.LittleEndian.Uint16(b), err
}

func readU32(r *bytes.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := r.Read(b)
	return binary.LittleEndian.Uint32(b), err
}

func readI32(r *bytes.Reader) (int32, error) {
	b := make([]byte, 4)
	_, err := r.Read(b)
	return int32(binary.LittleEndian.Uint32(b)), err
}

func readString(r *bytes.Reader) (string, error) {
	n, err := readU16(r)
	if err != nil {
		return "", err
	}
	b := make([]byte, n)
	_, err = r.Read(b)
	return string(b), err
}

package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	constMSBByte   = 0x80
	constNum1Byte  = 0
	constNum2Bytes = 1
	constNum3Bytes = 2
	constNum4Bytes = 3

	dateTypeSure   = 0x00
	dateTypeAbout  = 0x01
	dateTypeMaybe  = 0x02
	dateTypeBefore = 0x03
	dateTypeAfter  = 0x04
	// dateTypeOr
	// dateTypeBetween
	// dateTypeBasedOnAge

	errOpenFile = "could not open file: %w"
	errReadFile = "could not read file: %w"

	constUint32Size = 4

	constDaySepMonthSepLength = 4
	constDaySep               = 0x10
	constMonthSep             = 0x18
)

var errCouldNotParseNumber = errors.New("could not parse number")

func FileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

//nolint
func DumpByteSlice(b []byte) {
	var a [16]byte
	n := (len(b) + 15) &^ 15
	for i := 0; i < n; i++ {
		if i%16 == 0 {
			fmt.Printf("%08x", i)
		}
		if i%8 == 0 {
			fmt.Print(" ")
		}
		if i < len(b) {
			fmt.Printf(" %02x", b[i])
		} else {
			fmt.Print("   ")
		}
		if i >= len(b) {
			a[i%16] = ' '
		} else if b[i] < 32 || b[i] > 126 {
			a[i%16] = '.'
		} else {
			a[i%16] = b[i]
		}
		if i%16 == 15 {
			fmt.Printf("  %s\n", string(a[:]))
		}
	}
}

func ReadVarInt(r io.Reader) (int, error) {
	buf := make([]byte, constUint32Size)

	i := 0
	for i < constUint32Size {
		if err := binary.Read(r, binary.BigEndian, &buf[i]); err != nil {
			return 0, fmt.Errorf(errReadFile, err)
		}

		if buf[i]&constMSBByte != constMSBByte {
			break
		}

		i++
	}

	switch i {
	case constNum1Byte:
		return int(buf[0]), nil
	case constNum2Bytes:
		return int(buf[1])<<7 + int(buf[0])&0x7f, nil
	case constNum3Bytes:
		return int(buf[2])<<7 + int(buf[1])<<7 + int(buf[0])&0x7f, nil
	case constNum4Bytes:
		return int(buf[3])<<7 + int(buf[2])<<7 + int(buf[1])<<7 + int(buf[0])&0x7f, nil
	}

	return 0, errCouldNotParseNumber
}

func ReadDate(r io.Reader) (day, month, year int, err error) {
	daySepMonthSep := make([]byte, constDaySepMonthSepLength)
	if err = binary.Read(r, binary.BigEndian, &daySepMonthSep); err != nil {
		return 0, 0, 0, fmt.Errorf(errReadFile, err)
	}

	day = int(daySepMonthSep[0])

	if daySepMonthSep[1] != constDaySep {
		return day, 0, 0, fmt.Errorf("invalid day separator: %#02x", daySepMonthSep[1])
	}

	month = int(daySepMonthSep[2])

	if daySepMonthSep[3] != constMonthSep {
		return day, month, 0, fmt.Errorf("invalid month separator: %#02x", daySepMonthSep[3])
	}

	year, err = ReadVarInt(r)
	if err != nil {
		return day, month, year, fmt.Errorf("could not get year: %w", err)
	}

	return day, month, year, nil
}

func ReadString(r io.Reader) (string, error) {
	var textLength byte
	if err := binary.Read(r, binary.BigEndian, &textLength); err != nil {
		return "", fmt.Errorf(errReadFile, err)
	}

	text := make([]byte, textLength)
	if err := binary.Read(r, binary.BigEndian, &text); err != nil {
		return "", fmt.Errorf(errReadFile, err)
	}

	return string(text), nil
}

func ReadBytes(r io.Reader) ([]byte, error) {
	var recordLength byte
	if err := binary.Read(r, binary.BigEndian, &recordLength); err != nil {
		return nil, fmt.Errorf(errReadFile, err)
	}

	record := make([]byte, recordLength)
	if err := binary.Read(r, binary.BigEndian, &record); err != nil {
		return record, fmt.Errorf(errReadFile, err)
	}

	return record, nil
}

func ReadIdx(path string) (idx []uint, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return idx, fmt.Errorf(errOpenFile, err)
	}

	defer fi.Close()

	for {
		var offset uint32

		err = binary.Read(fi, binary.BigEndian, &offset)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return idx, fmt.Errorf(errReadFile, err)
		}

		idx = append(idx, uint(offset))
	}

	return idx, nil
}

func DateType(b byte) string {
	switch b {
	case dateTypeSure:
		return "="
	case dateTypeAbout:
		return "~="
	case dateTypeMaybe:
		return "?"
	case dateTypeBefore:
		return "<"
	case dateTypeAfter:
		return ">"
	default:
		return "unknow date type"
	}
}

package database

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
)

type BaseInfo struct {
	NbPersons uint32
	Sosa      uint32
	RootSosa  uint32
	Timestamp int64
}

func ReadInfoBase(path string) (*BaseInfo, error) {
	fileFullPath := filepath.Join(path, "pb_base_info.dat")
	if !utils.FileExists(fileFullPath) {
		return nil, fmt.Errorf("%w: %s", utils.ErrFileMissing, "pb_base_info.dat")
	}

	f, err := os.Open(fileFullPath)
	if err != nil {
		return nil, fmt.Errorf("base info file could not be opened: %w", err)
	}

	defer f.Close()

	buf := make([]byte, utils.ConstUint32Bytes)

	if _, err = io.ReadFull(f, buf); err != nil {
		return nil, fmt.Errorf(utils.ErrRead, err)
	}

	b := BaseInfo{}
	b.NbPersons = binary.BigEndian.Uint32(buf)

	if _, err = io.ReadFull(f, buf); err != nil {
		return &b, fmt.Errorf(utils.ErrRead, err)
	}

	b.Sosa = binary.BigEndian.Uint32(buf)

	var unknown byte
	if err = binary.Read(f, binary.BigEndian, &unknown); err != nil {
		return &b, fmt.Errorf(utils.ErrRead, err)
	}

	if _, err = io.ReadFull(f, buf); err != nil {
		return &b, fmt.Errorf(utils.ErrRead, err)
	}

	b.RootSosa = binary.BigEndian.Uint32(buf)

	if _, err = io.ReadFull(f, buf); err != nil {
		return &b, fmt.Errorf(utils.ErrRead, err)
	}

	timestampLength := binary.BigEndian.Uint32(buf)
	buf = make([]byte, timestampLength)

	if _, err = io.ReadFull(f, buf); err != nil {
		return &b, fmt.Errorf(utils.ErrRead, err)
	}

	timestamp, err := strconv.ParseInt(string(buf), utils.ConstDecBase, 0)
	if err != nil {
		return &b, fmt.Errorf("base info timestamp could not be parsed: %w", err)
	}

	b.Timestamp = timestamp

	return &b, nil
}

func GetInfoBaseString(b *BaseInfo) string {
	return fmt.Sprintf("base info: NbPersons=%d, Sosa=%d, RootSosa=%d, Date=%s\n\n",
		b.NbPersons, b.Sosa, b.RootSosa, time.Unix(b.Timestamp, 0))
}

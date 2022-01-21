package info

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
	"github.com/lunixbochs/struc"
)

const (
	errOpen      = "base info file could not be opened: %w"
	errUnpack    = "base info file could not be unpacked: %w"
	errTimestamp = "base info timestamp could not be parsed: %w"
)

var errMissing = errors.New("base info data file (pb_base_info.dat) missing")

// DumpStruct is the struct to unpack the baseinfo binary database.
type DumpStruct struct {
	NbPersons       uint   `struc:"big,uint32"`
	Sosa            uint   `struc:"big,uint32"`
	Unknown1        byte   `struc:"big,byte"` // FIX: no idea of the role of that field.
	RootSosa        uint   `struc:"big,uint32"`
	TimestampLength uint   `struc:"big,uint32"`
	Timestamp       string `struc:"big,[]byte,sizefrom=TimestampLength"`
}

// Info is the struct to externally expose what we extracted from the baseinfo binary database.
type BaseInfo struct {
	path      string
	nbPersons uint
	sosa      uint
	rootSosa  uint
	timestamp int64
}

func New(path string) (*BaseInfo, error) {
	fileFullPath := filepath.Join(path, "pb_base_info.dat")
	if !utils.FileExists(fileFullPath) {
		return nil, errMissing
	}

	return &BaseInfo{path: fileFullPath}, nil
}

func (b *BaseInfo) ReadInfoBase() error {
	f, err := os.Open(b.path)
	if err != nil {
		return fmt.Errorf(errOpen, err)
	}

	defer f.Close()

	dump := &DumpStruct{}

	err = struc.Unpack(f, dump)
	if err != nil {
		return fmt.Errorf(errUnpack, err)
	}

	b.nbPersons = dump.NbPersons
	b.sosa = dump.Sosa
	b.rootSosa = dump.RootSosa

	timestamp, err := strconv.Atoi(dump.Timestamp)
	if err != nil {
		return fmt.Errorf(errTimestamp, err)
	}

	b.timestamp = int64(timestamp)

	return nil
}

func (b *BaseInfo) GetNbPersons() uint {
	return b.nbPersons
}

func (b *BaseInfo) GetSosa() uint {
	return b.sosa
}

func (b *BaseInfo) GetRootSosa() uint {
	return b.rootSosa
}

func (b *BaseInfo) GetTimestamp() int64 {
	return b.timestamp
}

func (b *BaseInfo) GetInfoBaseString() string {
	return fmt.Sprintf("base info: NbPersons=%d, Sosa=%d, RootSosa=%d, Date=%s\n\n",
		b.nbPersons, b.sosa, b.rootSosa, time.Unix(b.timestamp, 0))
}

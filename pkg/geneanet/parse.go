package geneanet

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	errFileMissing   = "file %s missing"
	errOpenFile      = "can not open file: %w"
	errReadFile      = "can not read file: %w"
	errReadTimestamp = "can not convert timestamp: %w"
)

func fileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

type HeaderBaseInfo struct {
	NbPersons       uint32
	Sosa            uint32
	Unknown         [5]byte
	TimestampLength uint32
	Timestamp       [10]byte
}

func (g *Geneanet) ParseInfoBase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	fillFullPath := filepath.Join(g.outputDir, "pb_base_info.dat")
	if !fileExists(fillFullPath) {
		return fmt.Errorf(errFileMissing, fillFullPath)
	}

	f, err := os.Open(fillFullPath)
	if err != nil {
		return fmt.Errorf(errOpenFile, err)
	}

	defer f.Close()

	h := HeaderBaseInfo{}
	err = binary.Read(f, binary.BigEndian, &h)

	if err != nil {
		return fmt.Errorf(errReadFile, err)
	}

	g.nbPersons = h.NbPersons
	g.sosa = h.Sosa

	timestamp, err := strconv.ParseInt(string(h.Timestamp[:]), 10, 64)
	if err != nil {
		return fmt.Errorf(errReadTimestamp, err)
	}

	g.timestamp = timestamp

	return nil
}

func (g *Geneanet) GetNbPersons() uint32 {
	return g.nbPersons
}

func (g *Geneanet) GetSosa() uint32 {
	return g.sosa
}

func (g *Geneanet) GetTimestamp() int64 {
	return g.timestamp
}

package database

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/trois-six/geneparse/pkg/geneanet/utils"
)

type Database interface {
	CheckPath() error
	ReadIdx() error
	ReadData() error
	ReadIdxNote() error
	ReadDataNote() error
	GetIdx() []uint32
	GetData() [][]byte
	GetIdxNote() []uint32
	GetNotes() [][]utils.NoteWithTag
	Unmarshal() error
}

type commonDatabase struct {
	path             string
	baseFilePrefix   string
	idxFullPath      string
	dataFullPath     string
	idxNoteFullPath  string
	dataNoteFullPath string
	idx              []uint32
	data             [][]byte
	idxNote          []uint32
	dataNote         [][]utils.NoteWithTag
}

func (d *commonDatabase) CheckPath() (err error) {
	d.idxFullPath = filepath.Join(d.path, d.baseFilePrefix+".inx")
	if !utils.FileExists(d.idxFullPath) {
		return fmt.Errorf("%w: %s", utils.ErrFileMissing, d.idxFullPath)
	}

	d.dataFullPath = filepath.Join(d.path, d.baseFilePrefix+".dat")
	if !utils.FileExists(d.dataFullPath) {
		return fmt.Errorf("%w: %s", utils.ErrFileMissing, d.dataFullPath)
	}

	d.idxNoteFullPath = filepath.Join(d.path, d.baseFilePrefix+"_note.inx")
	if !utils.FileExists(d.idxNoteFullPath) {
		return fmt.Errorf("%w: %s", utils.ErrFileMissing, d.idxNoteFullPath)
	}

	d.dataNoteFullPath = filepath.Join(d.path, d.baseFilePrefix+"_note.dat")
	if !utils.FileExists(d.dataNoteFullPath) {
		return fmt.Errorf("%w: %s", utils.ErrFileMissing, d.dataNoteFullPath)
	}

	return nil
}

func (d *commonDatabase) ReadIdx() (err error) {
	d.idx, err = utils.ReadIdx(d.idxFullPath)
	if err != nil {
		return fmt.Errorf("failed reading database index file: %w", err)
	}

	return nil
}

func (d *commonDatabase) ReadData() (err error) {
	fb, err := os.Open(d.dataFullPath)
	if err != nil {
		return fmt.Errorf("failed opening database data file: %w", err)
	}

	defer fb.Close()

	buf := make([]byte, utils.ConstUint32Bytes)
	if _, err := io.ReadFull(fb, buf); err != nil {
		return fmt.Errorf(utils.ErrRead, err)
	}

	size := binary.BigEndian.Uint32(buf)

	var sizeRead uint32

	d.data = make([][]byte, 0, len(d.idx))

	for range d.idx {
		dataSize, data, err := utils.ReadBytes(fb)
		if err != nil {
			return fmt.Errorf(utils.ErrRead, err)
		}

		sizeRead += dataSize

		d.data = append(d.data, data)
	}

	if sizeRead != size {
		return fmt.Errorf("%w: %s", utils.ErrFileMalFormatted, d.dataFullPath)
	}

	return nil
}

func (d *commonDatabase) ReadIdxNote() (err error) {
	d.idxNote, err = utils.ReadIdx(d.idxNoteFullPath)
	if err != nil {
		return fmt.Errorf("failed reading database note index file: %w", err)
	}

	return nil
}

func (d *commonDatabase) ReadDataNote() (err error) {
	fb, err := os.Open(d.dataNoteFullPath)
	if err != nil {
		return fmt.Errorf("failed opening database data note file: %w", err)
	}

	defer fb.Close()

	buf := make([]byte, utils.ConstUint32Bytes)
	if _, err := io.ReadFull(fb, buf); err != nil {
		return fmt.Errorf(utils.ErrRead, err)
	}

	d.dataNote = make([][]utils.NoteWithTag, len(d.idxNote))

	for k := range d.idxNote {
		if d.idxNote[k] == 0 {
			continue
		}

		_, data, err := utils.ReadBytes(fb)
		if err != nil {
			return fmt.Errorf(utils.ErrRead, err)
		}

		d.dataNote[k] = utils.ExplodeNote(string(data))
	}

	return nil
}

func (d *commonDatabase) GetIdx() []uint32 {
	return d.idx
}

func (d *commonDatabase) GetData() [][]byte {
	return d.data
}

func (d *commonDatabase) GetIdxNote() []uint32 {
	return d.idxNote
}

func (d *commonDatabase) GetNotes() [][]utils.NoteWithTag {
	return d.dataNote
}

func PopulateDatabases(databases []Database) error {
	for _, db := range databases {
		if err := db.CheckPath(); err != nil {
			return fmt.Errorf("failed checking databases: %w", err)
		}

		if err := db.ReadIdx(); err != nil {
			return fmt.Errorf("failed reading index from database: %w", err)
		}

		if err := db.ReadData(); err != nil {
			return fmt.Errorf("failed reading data from database: %w", err)
		}

		if err := db.Unmarshal(); err != nil {
			return fmt.Errorf("failed unmarshaling data from database: %w", err)
		}

		if err := db.ReadIdxNote(); err != nil {
			return fmt.Errorf("failed reading note index from database: %w", err)
		}

		if err := db.ReadDataNote(); err != nil {
			return fmt.Errorf("failed reading note data from database: %w", err)
		}
	}

	return nil
}

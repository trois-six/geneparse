package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/elliotchance/gedcom"
)

const (
	ErrRead       = "could not read file: %w"
	ErrParseInput = "could not parse input: %w"

	ConstUint32Bytes = 4
	ConstDecBase     = 10
	constNoteMaxLen  = 71
)

var (
	ErrFileMissing      = errors.New("file missing")
	ErrFileMalFormatted = errors.New("file malformatted")
	ErrDirDoesNotExist  = errors.New("directory does not exist")
	ErrDirMustBeADir    = errors.New("directory must be a directory")
)

func FileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func ReadBytes(r io.Reader) (uint32, []byte, error) {
	buf := make([]byte, ConstUint32Bytes)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, nil, fmt.Errorf(ErrRead, err)
	}

	size := binary.BigEndian.Uint32(buf)

	data := make([]byte, size)
	if err := binary.Read(r, binary.BigEndian, &data); err != nil {
		return 0, data, fmt.Errorf(ErrRead, err)
	}

	return ConstUint32Bytes + size, data, nil
}

func PointerStr(prefix string, id int32) string {
	return prefix + strconv.FormatInt(int64(id), ConstDecBase)
}

func ReadIdx(fileFullPath string) (idx []uint32, err error) {
	fi, err := os.Open(fileFullPath)
	if err != nil {
		return idx, fmt.Errorf("failed opening database index file: %w", err)
	}

	defer fi.Close()

	buf := make([]byte, ConstUint32Bytes)

	for {
		_, err = io.ReadFull(fi, buf)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return idx, fmt.Errorf(ErrRead, err)
		}

		idx = append(idx, binary.BigEndian.Uint32(buf))
	}

	return idx, nil
}

type NoteWithTag struct {
	tag  gedcom.Tag
	note string
}

func ExplodeNote(note string) []NoteWithTag {
	m := regexp.MustCompile("(<[b|B][r|R][ ]?[/]?>[\n]?|\n)")
	splitted := m.Split(note, -1)

	notesWithTag := make([]NoteWithTag, len(splitted))
	notesWithTag[0] = NoteWithTag{
		tag:  gedcom.TagNote,
		note: splitted[0],
	}

	if len(splitted) > 0 {
		for k, note := range splitted[1:] {
			notesWithTag[k+1] = NoteWithTag{
				tag:  gedcom.TagContinued,
				note: note,
			}
		}
	}

	final := make([]NoteWithTag, 0, len(splitted))

	for _, note := range notesWithTag {
		if len(note.note) == 0 {
			final = append(final, note)

			continue
		}

		var chuncks []NoteWithTag

		runes := []rune(note.note)

		for i := 0; i < len(runes); i += constNoteMaxLen {
			end := i + constNoteMaxLen
			if end > len(runes) {
				end = len(runes)
			}

			chuncks = append(chuncks, NoteWithTag{
				tag: func() gedcom.Tag {
					if i > 0 {
						return gedcom.TagConcatenation
					}

					return note.tag
				}(),
				note: string(runes[i:end]),
			})
		}

		final = append(final, chuncks...)
	}

	return final
}

func (n *NoteWithTag) GetTag() gedcom.Tag {
	return n.tag
}

func (n *NoteWithTag) GetNote() string {
	return n.note
}

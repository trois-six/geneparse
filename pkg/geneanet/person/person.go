package person

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
	"github.com/lunixbochs/struc"
)

const (
	constRecordJob             = 0x01
	constRecordID              = 0x08
	constRecordSex             = 0x10
	constRecordLastName        = 0x1a
	constRecordChild           = 0x1a
	constRecordFirstName       = 0x22
	constUnknown7              = 0x22
	constRecordExtendedInfo    = 0x28
	constMarriageSrc           = 0x32
	constRecordSrc             = 0x3a
	constRecordBaptismSrc      = 0x41
	constRecordNickName        = 0x42
	constUnknown1              = 0x48
	constRecordOtherFirstName  = 0x4a
	constRecordOtherLastName   = 0x52
	constRecordDateInfo        = 0x58
	constRecordBirth           = 0x62
	constRecordBirthLocation   = 0x6a
	constRecordBirthNote       = 0x72
	constRecordBaptism         = 0x7a
	constUnknown2              = 0x80
	constRecordBaptismLocation = 0x82
	constUnknown3              = 0x8a
	constRecordDeath           = 0x92
	constRecordDeathLocation   = 0x9a
	constRecordDeathNote       = 0xa2
	constUnknown4              = 0xa8
	constUnknown8              = 0xca
	constUnknown5              = 0xf0
	constUnknown6              = 0xf8

	constHexBase    = 16
	constUint32Size = 4

	errUnpack = "base person file could not be unpacked: %w"
	errSeek   = "base person file could not be seeked: %w"
	errRead   = "base person file could not be read: %w"
)

var (
	errIndexMissing = errors.New("base person index file (pb_base_person.inx) missing")
	errDataMissing  = errors.New("base person data file (pb_base_person.dat) missing")
)

func ReadPersonIdx(path string) ([]uint, error) {
	fileFullPath := filepath.Join(path, "pb_base_person.inx")
	if !utils.FileExists(fileFullPath) {
		return []uint{}, errIndexMissing
	}

	return utils.ReadIdx(fileFullPath)
}

type DumpStruct struct {
	Size int    `struc:"big,uint32"`
	Data []byte `struc:"big,[]byte,sizefrom=Size"`
}

type BasePerson struct {
	path string
	size uint32
	data []byte
}

type Person struct {
	id         uint
	sex        bool
	firstNames []string
	lastNames  []string
	nickNames  []string
	job        string
}

func New(path string) (*BasePerson, error) {
	fileFullPath := filepath.Join(path, "pb_base_person.dat")
	if !utils.FileExists(fileFullPath) {
		return nil, errDataMissing
	}

	return &BasePerson{path: fileFullPath}, nil
}

func (b *BasePerson) ReadPersonBase() error {
	data, err := ioutil.ReadFile(b.path)
	if err != nil {
		return fmt.Errorf(errRead, err)
	}

	r := bytes.NewReader(data)

	if err := binary.Read(r, binary.BigEndian, &b.size); err != nil {
		return fmt.Errorf(errRead, err)
	}

	b.data = make([]byte, r.Len())
	if err := binary.Read(r, binary.BigEndian, &b.data); err != nil {
		return fmt.Errorf(errRead, err)
	}

	return nil
}

func (b *BasePerson) GetPersonBaseString() string {
	return fmt.Sprintf("base person: Size=%d\n\n", b.size)
}

func ReadDateEvent(b []byte) (string, error) {
	if b[6] != 0x08 {
		return "", fmt.Errorf("invalid date magic: %#02x", b[7])
	}

	r := bytes.NewReader(b[7:])

	day, month, year, err := utils.ReadDate(r)
	if err != nil {
		return "", fmt.Errorf("could not get date: %w", err)
	}

	var magicSep byte
	if err := binary.Read(r, binary.BigEndian, &magicSep); err != nil {
		return "", fmt.Errorf("could not read dated event magic separator: %w", err)
	}

	if magicSep == 0x20 {
		return fmt.Sprintf("%s%d/%d/%d", utils.DateType(b[3]), day, month, year), nil
	}

	return "", errors.New("could not find location separator")
}

func (b *BasePerson) GetPerson(id, offset uint) (*Person, error) {
	rBase := bytes.NewReader(b.data)

	log.Printf("Person id=%d offset: %08x", id, offset)

	if _, err := rBase.Seek(int64(offset), io.SeekCurrent); err != nil {
		return nil, fmt.Errorf(errSeek, err)
	}

	dump := &DumpStruct{}
	if err := struc.Unpack(rBase, dump); err != nil {
		return nil, fmt.Errorf(errUnpack, err)
	}

	r := bytes.NewReader(dump.Data)

	p := Person{}
	p.firstNames = []string{"NN"}
	p.lastNames = []string{"NN"}

	var recordType, tmp byte

	var stop bool
	for !stop {
		if err := binary.Read(r, binary.BigEndian, &recordType); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf(errRead, err)
		}

		switch recordType {
		case constRecordID: // id
			log.Printf("read id: %#02x", constRecordID)

			recordID, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.id = uint(recordID)

			log.Printf("read id: %d", p.id)
		case constRecordSex: // sex
			log.Printf("read sex: %#02x", constRecordSex)

			sex, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, err
			}

			p.sex = sex != 0
		case constRecordLastName: // last name
			log.Printf("read last name: %#02x", constRecordLastName)

			lastName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.lastNames[0] = string(lastName)
		case constRecordFirstName: // first name
			log.Printf("read first name: %#02x", constRecordFirstName)

			firstName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.firstNames[0] = string(firstName)
		case constRecordOtherLastName: // other last name
			log.Printf("read other last name: %#02x", constRecordOtherLastName)

			lastName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.lastNames = append(p.lastNames, string(lastName))
		case constRecordOtherFirstName: // other first name
			log.Printf("read other first name: %#02x", constRecordOtherFirstName)

			firstName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.firstNames = append(p.firstNames, string(firstName))
		case constRecordNickName: // nick names
			log.Printf("read nick name: %#02x", constRecordNickName)

			nickName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.nickNames = append(p.nickNames, string(nickName))
		case constRecordJob: // job
			log.Printf("read job: %#02x", constRecordJob)

			job, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.job = job
		case constRecordExtendedInfo: // extended info
			log.Printf("read extended info: %#02x", constRecordExtendedInfo)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("sep extended infos: %#02x", tmp)

			stop = true
		default:
			log.Printf("read unknown field %#02x, stop", recordType)
		}
	}

	stop = false

	for !stop {
		if err := binary.Read(r, binary.BigEndian, &recordType); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf(errRead, err)
		}

		switch recordType {
		case constRecordOtherLastName: // other last name
			log.Printf("read other last name: %#02x", constRecordOtherLastName)

			lastName, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.lastNames = append(p.lastNames, string(lastName))
		case constRecordJob: // job
			log.Printf("read job: %#02x", constRecordJob)

			job, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			p.job = job
		case constRecordID: // ??
			log.Printf("read id: %#02x", constRecordID)

			idNumber, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, err
			}

			log.Printf("read id number: %#v", idNumber)
		case constRecordChild: // child
			log.Printf("read child record: %#02x", constRecordChild)

			childRecord, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			recordDate, err := ReadDateEvent(childRecord)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("child record date: %s", recordDate)
		case constRecordDateInfo: // date info
			log.Printf("read date info: %#02x", constRecordDateInfo)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("sep date infos: %#02x", tmp)
		case constRecordBirth: // birth
			log.Printf("read birth record: %#02x", constRecordBirth)

			birthRecord, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			recordDate, err := ReadDateEvent(birthRecord)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("birth record date: %s", recordDate)
		case constRecordBaptism: // baptism
			log.Printf("read baptism record: %#02x", constRecordBaptism)

			baptismRecord, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			recordDate, err := ReadDateEvent(baptismRecord)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("baptism record date: %s", recordDate)
		case constRecordDeath: // death
			log.Printf("read death record: %#02x", constRecordDeath)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("sep death record: %#02x", tmp)

			deathRecord, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			recordDate, err := ReadDateEvent(deathRecord)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("death record date: %s", recordDate)
		case constRecordBirthLocation: // birth location
			log.Printf("read birth location: %#02x", constRecordBirthLocation)

			birthLocation, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("birth location: %s", birthLocation)
		case constRecordDeathLocation: // death location
			log.Printf("read death location: %#02x", constRecordDeathLocation)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("sep death location: %#02x", tmp)

			deathLocation, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("death location: %s", deathLocation)
		case constRecordBirthNote: // birth note
			log.Printf("read birth note: %#02x", constRecordBirthNote)

			birthNote, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("birth note: %s", birthNote)
		case constRecordDeathNote: // death note
			log.Printf("read death note: %#02x", constRecordDeathNote)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("sep death note: %#02x", tmp)

			deathNote, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("death note: %s", deathNote)
		case constUnknown1: // unknown1
			log.Printf("read unknown1: %#02x", constUnknown1)

			unknown1, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("unknown1: %#v", unknown1)
		case constUnknown2: // unknown2
			log.Printf("read unknown2: %#02x", constUnknown2)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			unknown2, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("unknown2: %#v, %#v", tmp, unknown2)
		case constUnknown3: // unknown3
			// log.Printf("sep unknown3: %#02x", tmp)
			log.Printf("read unknown3: %#02x", constUnknown3)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			unknown3, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("unknown3: %#v, %#v", tmp, unknown3)
		case constUnknown4: // unknown4
			log.Printf("read unknown4: %#02x", constUnknown4)

			unknown4, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("unknown4: %#v", unknown4)
		case constUnknown5: // unknown5
			log.Printf("read unknown5: %#02x", constUnknown5)

			unknown5, err := utils.ReadBytes(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("unknown5: %#v", unknown5)
		case constUnknown6: // unknown6
			log.Printf("read unknown6: %#02x", constUnknown6)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			unknown6, err := utils.ReadVarInt(r)
			if err != nil {
				return nil, err
			}

			log.Printf("unknown6: %#v, %#v", tmp, unknown6)
		case constUnknown7: // unknown7
			log.Printf("read unknown7: %#02x", constUnknown7)

			event, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read unknown7: %s", event)
		case constRecordSrc: // source of record
			log.Printf("read record source: %#02x", constRecordSrc)

			source, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read record source: %s", source)
		case constMarriageSrc: // source of marriage
			log.Printf("read marriage source: %#02x", constMarriageSrc)

			marriageSource, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read marriage source: %s", marriageSource)
		case constUnknown8: // unknown8
			log.Printf("read unknown8: %#02x", constUnknown8)
		case constRecordBaptismLocation: // baptism location
			log.Printf("read baptism location: %#02x", constRecordBaptismLocation)

			if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read baptism location sep: %#v", tmp)

			baptismLocation, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read baptism location: %s", baptismLocation)
		case constRecordBaptismSrc: // source of baptism
			log.Printf("read baptism source: %#02x", constRecordBaptismSrc)

			baptismSource, err := utils.ReadString(r)
			if err != nil {
				return nil, fmt.Errorf(errRead, err)
			}

			log.Printf("read baptism source: %s", baptismSource)
		default:
			log.Printf("read unknown field %#02x, stop", recordType)

			stop = true
		}
	}

	offsetBase := uint(float64((offset+constUint32Size)/constHexBase) * constHexBase)
	offsetAfterBase := (offset + constUint32Size) % constHexBase
	log.Printf("id=%d offset=%08x+%x firstNames=%#v lastNames=%#v nickNames=%#v sex=%v",
		id, offsetBase, offsetAfterBase, p.firstNames, p.lastNames, p.nickNames, p.sex)

	if stop {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf(errRead, err)
		}

		pData := make([]byte, r.Len())
		if err := binary.Read(r, binary.BigEndian, &pData); err != nil {
			return nil, fmt.Errorf(errRead, err)
		}

		utils.DumpByteSlice(pData)
	}

	return nil, nil
}

/*

Père DEAD
a8 01
   01
f0 01
   00
80 02
   00 8a
02 13
   08 00 1a 0f 10 00 18 00 22 09 08 0a 10 0a 18 c6 0f 20 00
8a 02
   04 08
32 48 01
8a 02 13
08 02
   1a 0f 10 00 18 00 22 09 08 00 10 00 18 da 0f 20 00

Mère BEEF
a8 01
   00
f0 01
   00
80 02
   00 8a
02 13
   08 00 1a 0f 10 00 18 00 22 09 08 0b 10 0b 18 c7 0f 20 00
8a 02
   04 08
32 48 00

Fils DEAD
Birth
   62 0f 10 00 18 00 22 09 08 02 10 02 18 d2 0f 20 00
???
a8 01 00
f0 01 00
f8 01 00
8a 02 13
08 00
   1a 0f 10 00 18 00 22 09 08 02 10 02 18 d2 0f 20 00



Florence ALFONSI
a8 01 01
f0 01 00
f8 01 02
80 02 01
8a 02 4c
08 00
   1a 0f 10 00 18 00 22 09 08 1b 10 06 18 af 0f 20 00



Jeanne ERRAUD
a8 01 03
f0 01 01
f8 01 a9 01
80 02 8f 04
8a 02 99 01
08 00
   1a 0f 10 00 18 00 22 09 08 02 10 03 18 f5 0c 20 00

48 84 08
8a 02 9d 01 08 02
   1a 0f 10 00 18 00 7H 22 09 08 06 10 05 18 97 0d 20 00
*/

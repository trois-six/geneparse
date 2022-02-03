package geneanet

import (
	"fmt"

	"github.com/Trois-Six/geneparse/pkg/geneanet/database"
	"github.com/Trois-Six/geneparse/pkg/geneanet/gengedcom"
	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
)

type Geneanet struct {
	path string

	nbPersons uint32
	sosa      uint32
	rootSosa  uint32
	timestamp int64
}

func New(path string) (*Geneanet, error) {
	if !utils.FileExists(path) {
		return nil, fmt.Errorf("%w: %s", utils.ErrDirDoesNotExist, path)
	}

	return &Geneanet{path: path}, nil
}

func (g *Geneanet) Parse() error {
	info, err := database.ReadInfoBase(g.path)
	if err != nil {
		return fmt.Errorf("could not read base info: %w", err)
	}

	g.nbPersons = info.NbPersons
	g.sosa = info.Sosa
	g.rootSosa = info.RootSosa
	g.timestamp = info.Timestamp

	person := database.NewPerson(g.path)
	family := database.NewFamily(g.path)

	if err = database.PopulateDatabases([]database.Database{person, family}); err != nil {
		return fmt.Errorf("databases populate failed: %w", err)
	}

	genGedcom := gengedcom.New("test")
	if err = genGedcom.Write("test",
		person.GetPersons(),
		family.GetFamilies(),
		person.GetNotes(),
		family.GetNotes(),
	); err != nil {
		return fmt.Errorf("could not write gedcom: %w", err)
	}

	return nil
}

package geneanet

import (
	"fmt"
	"log"

	"github.com/Trois-Six/geneparse/pkg/geneanet/info"
	"github.com/Trois-Six/geneparse/pkg/geneanet/person"
	"github.com/Trois-Six/geneparse/pkg/geneanet/utils"
)

type Geneanet struct {
	path      string
	personIdx []uint
	nbPersons uint
	sosa      uint
	rootSosa  uint
	timestamp int64
}

func New(path string) (*Geneanet, error) {
	if !utils.FileExists(path) {
		return nil, fmt.Errorf("could not create Geneanet, directory doesn't exist: %s", path)
	}

	return &Geneanet{path: path}, nil
}

func (g *Geneanet) Parse() error {
	info, err := info.New(g.path)
	if err != nil {
		return fmt.Errorf("could not initialize base info: %w", err)
	}

	if err = info.ReadInfoBase(); err != nil {
		return fmt.Errorf("could not read base info: %w", err)
	}

	g.nbPersons = info.GetNbPersons()
	g.sosa = info.GetSosa()
	g.rootSosa = info.GetRootSosa()
	g.timestamp = info.GetTimestamp()

	log.Print(info.GetInfoBaseString())

	if g.personIdx, err = person.ReadPersonIdx(g.path); err != nil {
		return fmt.Errorf("could not read person index: %w", err)
	}

	person, err := person.New(g.path)
	if err != nil {
		return fmt.Errorf("could not initialize base person: %w", err)
	}

	if err = person.ReadPersonBase(); err != nil {
		return fmt.Errorf("could not read base person: %w", err)
	}

	log.Print(person.GetPersonBaseString())

	// for i := uint(0); i < 200; i++ {
	// 	log.Printf("idx=%08d, offset=%08x", i, g.personidx[i])
	// }

	for i := uint(0); i < 200; i++ {
		if _, err := person.GetPerson(i, g.personIdx[i]); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	// for i := uint(1000); i < 1030; i++ {
	// 	if _, err := person.GetPerson(i, g.personIdx[i]); err != nil {
	// 		return err
	// 	}
	// 	fmt.Printf("\n")
	// }

	// if _, err := person.GetPerson(0, g.personIdx[0]); err != nil {
	// 	return err
	// }

	// if _, err := g.GetPerson(8294); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(8295); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(8296); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(8392); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(8505); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(8506); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(1); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(g.GetRootSosa()); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(0x2144); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(39); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(53); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(126); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(1987); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(1); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	// if _, err := g.GetPerson(34); err != nil {
	// 	return fmt.Errorf(errFailedGetPerson, err)
	// }

	return nil
}

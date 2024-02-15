package database

import (
	"fmt"

	"github.com/trois-six/geneparse/pkg/geneanet/api"
	"google.golang.org/protobuf/proto"
)

type Person struct {
	persons []*api.Person
	commonDatabase
}

const personBasePrefix = "pb_base_person"

func NewPerson(path string) *Person {
	return &Person{
		commonDatabase: commonDatabase{
			path:           path,
			baseFilePrefix: personBasePrefix,
		},
	}
}

func (p *Person) Unmarshal() error {
	p.persons = make([]*api.Person, 0, len(p.idx))

	for _, personByte := range p.data {
		person := new(api.Person)

		if err := proto.Unmarshal(personByte, person); err != nil {
			return fmt.Errorf("failed parsing database person data file: %w", err)
		}

		p.persons = append(p.persons, person)
	}

	return nil
}

func (p *Person) GetPersons() []*api.Person {
	return p.persons
}

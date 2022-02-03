package database

import (
	"fmt"

	"github.com/Trois-Six/geneparse/pkg/geneanet/api"
	"google.golang.org/protobuf/proto"
)

type Family struct {
	families []*api.Family
	commonDatabase
}

const familyBasePrefix = "pb_base_family"

func NewFamily(path string) *Family {
	return &Family{
		commonDatabase: commonDatabase{
			path:           path,
			baseFilePrefix: familyBasePrefix,
		},
	}
}

func (f *Family) Unmarshal() error {
	f.families = make([]*api.Family, 0, len(f.idx))

	for _, familyByte := range f.data {
		family := new(api.Family)

		if err := proto.Unmarshal(familyByte, family); err != nil {
			return fmt.Errorf("failed parsing database family data file: %w", err)
		}

		f.families = append(f.families, family)
	}

	return nil
}

func (f *Family) GetFamilies() []*api.Family {
	return f.families
}

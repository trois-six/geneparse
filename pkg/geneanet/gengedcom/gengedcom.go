package gengedcom

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trois-six/geneparse/pkg/geneanet/api"
	"github.com/trois-six/geneparse/pkg/geneanet/utils"
	"github.com/elliotchance/gedcom"
)

// mapEventNameTagName commes from api.proto and
// https://github.com/geneweb/geneweb/blob/master/bin/gwb2ged/gwb2gedLib.ml.
var mapPrecisionDateString = map[api.Precision]string{ // nolint:gochecknoglobals
	api.Precision_SURE:   "",
	api.Precision_ABOUT:  "ABT ",
	api.Precision_MAYBE:  "EST ",
	api.Precision_BEFORE: "BEF ",
	api.Precision_AFTER:  "AFT ",
}

// mapEventNameTagName commes from api.proto and
// https://github.com/geneweb/geneweb/blob/master/bin/gwb2ged/gwb2gedLib.ml.
var mapEventNameTagName = map[api.EventName]gedcom.Tag{ // nolint:gochecknoglobals
	api.EventName_EPERS_BIRTH:                   gedcom.TagBirth,
	api.EventName_EPERS_BAPTISM:                 gedcom.TagBaptism,
	api.EventName_EPERS_DEATH:                   gedcom.TagDeath,
	api.EventName_EPERS_BURIAL:                  gedcom.TagBurial,
	api.EventName_EPERS_CREMATION:               gedcom.TagCremation,
	api.EventName_EPERS_ACCOMPLISHMENT:          gedcom.TagFromString("Accomplishment"),
	api.EventName_EPERS_ACQUISITION:             gedcom.TagFromString("Acquisition"),
	api.EventName_EPERS_ADHESION:                gedcom.TagFromString("Membership"),
	api.EventName_EPERS_BAPTISMLDS:              gedcom.TagLDSBaptism,
	api.EventName_EPERS_BARMITZVAH:              gedcom.TagBarMitzvah,
	api.EventName_EPERS_BATMITZVAH:              gedcom.TagBasMitzvah,
	api.EventName_EPERS_BENEDICTION:             gedcom.TagBlessing,
	api.EventName_EPERS_CHANGENAME:              gedcom.TagChange,
	api.EventName_EPERS_CIRCUMCISION:            gedcom.TagFromString("Circumcision"),
	api.EventName_EPERS_CONFIRMATION:            gedcom.TagConfirmation,
	api.EventName_EPERS_CONFIRMATIONLDS:         gedcom.TagLDSConfirmation,
	api.EventName_EPERS_DECORATION:              gedcom.TagFromString("Award"),
	api.EventName_EPERS_DEMOBILISATIONMILITAIRE: gedcom.TagFromString("Military discharge"),
	api.EventName_EPERS_DIPLOMA:                 gedcom.TagFromString("Degree"),
	api.EventName_EPERS_DISTINCTION:             gedcom.TagFromString("Distinction"),
	api.EventName_EPERS_DOTATION:                gedcom.TagEndowment,
	api.EventName_EPERS_DOTATIONLDS:             gedcom.TagFromString("DotationLDS"),
	api.EventName_EPERS_EDUCATION:               gedcom.TagEducation,
	api.EventName_EPERS_ELECTION:                gedcom.TagFromString("Election"),
	api.EventName_EPERS_EMIGRATION:              gedcom.TagEmigration,
	api.EventName_EPERS_EXCOMMUNICATION:         gedcom.TagFromString("Excommunication"),
	api.EventName_EPERS_FAMILYLINKLDS:           gedcom.TagFromString("Family link LDS"),
	api.EventName_EPERS_FIRSTCOMMUNION:          gedcom.TagFirstCommunion,
	api.EventName_EPERS_FUNERAL:                 gedcom.TagFromString("Funeral"),
	api.EventName_EPERS_GRADUATE:                gedcom.TagGraduation,
	api.EventName_EPERS_HOSPITALISATION:         gedcom.TagFromString("Hospitalization"), //nolint:misspell
	api.EventName_EPERS_ILLNESS:                 gedcom.TagFromString("Illness"),
	api.EventName_EPERS_IMMIGRATION:             gedcom.TagImmigration,
	api.EventName_EPERS_LISTEPASSENGER:          gedcom.TagFromString("Passenger list"),
	api.EventName_EPERS_MILITARYDISTINCTION:     gedcom.TagFromString("Military distinction"),
	api.EventName_EPERS_MILITARYPROMOTION:       gedcom.TagFromString("Military promotion"),
	api.EventName_EPERS_MILITARYSERVICE:         gedcom.TagFromString("Military service"),
	api.EventName_EPERS_MOBILISATIONMILITAIRE:   gedcom.TagFromString("Military mobilization"),
	api.EventName_EPERS_NATURALISATION:          gedcom.TagNaturalization, //nolint:misspell
	api.EventName_EPERS_OCCUPATION:              gedcom.TagOccupation,
	api.EventName_EPERS_ORDINATION:              gedcom.TagOrdination,
	api.EventName_EPERS_PROPERTY:                gedcom.TagProperty,
	api.EventName_EPERS_RECENSEMENT:             gedcom.TagCensus,
	api.EventName_EPERS_RESIDENCE:               gedcom.TagResidence,
	api.EventName_EPERS_RETIRED:                 gedcom.TagRetirement,
	api.EventName_EPERS_SCELLENTCHILDLDS:        gedcom.TagSealingChild,
	api.EventName_EPERS_SCELLENTPARENTLDS:       gedcom.TagFromString("Scellent parent LDS"),
	api.EventName_EPERS_SCELLENTSPOUSELDS:       gedcom.TagSealingSpouse,
	api.EventName_EPERS_VENTEBIEN:               gedcom.TagFromString("Property sale"),
	api.EventName_EPERS_WILL:                    gedcom.TagWill,

	api.EventName_EFAM_MARRIAGE:          gedcom.TagMarriage,
	api.EventName_EFAM_NO_MARRIAGE:       gedcom.TagFromString("unmarried"),
	api.EventName_EFAM_NO_MENTION:        gedcom.TagFromString("nomen"),
	api.EventName_EFAM_ENGAGE:            gedcom.TagEngagement,
	api.EventName_EFAM_DIVORCE:           gedcom.TagDivorce,
	api.EventName_EFAM_SEPARATED:         gedcom.TagFromString("SEP"),
	api.EventName_EFAM_ANNULATION:        gedcom.TagAnnulment,
	api.EventName_EFAM_MARRIAGE_BANN:     gedcom.TagMarriageBann,
	api.EventName_EFAM_MARRIAGE_CONTRACT: gedcom.TagMarriageContract,
	api.EventName_EFAM_MARRIAGE_LICENSE:  gedcom.TagMarriageLicence,
	api.EventName_EFAM_PACS:              gedcom.TagFromString("pacs"),
	api.EventName_EFAM_RESIDENCE:         gedcom.TagFromString("residence"),
}

// mapMarriageTypeTagName commes from api.proto and
// https://github.com/geneweb/geneweb/blob/master/bin/gwb2ged/gwb2gedLib.ml.
var mapMarriageTypeTagName = map[api.MarriageType]gedcom.Tag{ // nolint:gochecknoglobals
	api.MarriageType_MARRIED:                    gedcom.TagMarriage,
	api.MarriageType_NOT_MARRIED:                gedcom.TagFromString("unmarried"),
	api.MarriageType_ENGAGED:                    gedcom.TagEngagement,
	api.MarriageType_NO_SEXES_CHECK_NOT_MARRIED: gedcom.TagFromString("unmarried"),
	api.MarriageType_NO_MENTION:                 gedcom.TagFromString("nomen"),
	api.MarriageType_NO_SEXES_CHECK_MARRIED:     gedcom.TagMarriage,
	api.MarriageType_MARRIAGE_BANN:              gedcom.TagMarriageBann,
	api.MarriageType_MARRIAGE_CONTRACT:          gedcom.TagMarriageContract,
	api.MarriageType_MARRIAGE_LICENSE:           gedcom.TagMarriageLicence,
	api.MarriageType_PACS:                       gedcom.TagFromString("pacs"),
	api.MarriageType_RESIDENCE:                  gedcom.TagFromString("residence"),
}

type GenGedcom struct {
	path string
}

func New(path string) GenGedcom {
	return GenGedcom{
		path: path,
	}
}

func getEmptyDocument(name string) *gedcom.Document {
	currentTime := time.Now()
	doc := gedcom.NewDocument()
	doc.HasBOM = true
	doc.AddNode(gedcom.NewNode(gedcom.TagHeader, "", "",
		gedcom.NewNode(gedcom.TagGedcomInformation, "", "",
			gedcom.NewNode(gedcom.TagVersion, "5.5.5", ""),
			gedcom.NewNode(gedcom.TagFormat, "LINEAGE-LINKED", "",
				gedcom.NewNode(gedcom.TagVersion, "5.5.5", ""),
			),
		),
		gedcom.NewNode(gedcom.TagCharacterSet, "UTF-8", ""),
		gedcom.NewNode(gedcom.TagSource, "Geneanet", "",
			gedcom.NewNode(gedcom.TagName, "GEDCOM computed file from a Geneanet Tree", ""),
			gedcom.NewNode(gedcom.TagVersion, "5.5.5", ""),
			gedcom.NewNode(gedcom.TagCorporate, "geneparse", "",
				gedcom.NewNode(gedcom.TagWWW, "https://github.com/trois-six/geneparse", ""),
			),
		),
		gedcom.NewNode(gedcom.TagDate, currentTime.Format("02 Jan 2006"), "",
			gedcom.NewNode(gedcom.TagTime, currentTime.Format("15:04:05"), ""),
		),
		gedcom.NewNode(gedcom.TagFile, name+".ged", ""),
	))

	return doc
}

func getName(person *api.Person) gedcom.Node {
	if len(person.FirstnameAliases) > 0 {
		firstNameAliases := strings.Join(person.GetFirstnameAliases(), ",")

		return gedcom.NewNameNode(fmt.Sprintf(
			"\"%s\" %s /%s/",
			firstNameAliases, person.GetFirstname(), person.GetLastname(),
		))
	}

	return gedcom.NewNameNode(fmt.Sprintf(
		"%s /%s/",
		person.GetFirstname(), person.GetLastname(),
	))
}

func getSex(sex api.Sex) string {
	switch sex {
	case api.Sex_MALE:
		return "M"
	case api.Sex_FEMALE:
		return "F"
	case api.Sex_UNKNOWN:
		return "U"
	default:
		return "U"
	}
}

// TODO: manage other Date types than GREGORIAN.
func getDmy(date *api.Dmy) string {
	var dateString string

	if date.Day != nil && date.GetDay() != 0 {
		dateString += strconv.FormatInt(int64(date.GetDay()), utils.ConstDecBase) + " "
	}

	if date.Month != nil && date.GetMonth() != 0 {
		dateString += strings.ToUpper(time.Month(date.GetMonth()).String()[0:3]) + " "
	}

	// TODO: how to manage year 0?? and years are uint32... how did they manage negative years?
	if date.Year != nil {
		dateString += strconv.FormatInt(int64(date.GetYear()), utils.ConstDecBase) + " "
	}

	return strings.TrimRight(dateString, " ")
}

func getDate(date *api.Date) string {
	var dateString string

	switch prec := date.GetPrec(); {
	case prec >= api.Precision_SURE && prec <= api.Precision_AFTER:
		dateString = mapPrecisionDateString[prec] + getDmy(date.GetDmy())
	case prec == api.Precision_ORYEAR:
		dateString = "FROM " + getDmy(date.GetDmy())
		if date.Dmy2 != nil {
			dateString += " TO " + getDmy(date.GetDmy2())
		}
	case prec == api.Precision_YEARINT:
		dateString = "BET " + getDmy(date.GetDmy())
		if date.Dmy2 != nil {
			dateString += " AND " + getDmy(date.GetDmy2())
		}
	}

	return dateString
}

func getTitle(title *api.Title) gedcom.Node {
	t := gedcom.NewNode(gedcom.TagTitle, title.GetTitle()+", "+title.GetFief(), "")

	if title.GetDateBegin() != nil || title.GetDateEnd() != nil {
		dateString := ""
		if title.GetDateBegin() != nil {
			dateString = "FROM " + getDate(title.GetDateBegin()) + " "
		}

		if title.GetDateEnd() != nil {
			dateString += "TO " + getDate(title.GetDateEnd()) + " "
		}

		t.AddNode(gedcom.NewNode(gedcom.TagDate, strings.TrimRight(dateString, " "), ""))
	}

	return t
}

func getEvent(event *api.Event) gedcom.Node {
	t := gedcom.NewNode(mapEventNameTagName[event.GetName()], "", "")

	if event.Date != nil {
		t.AddNode(gedcom.NewDateNode(getDate(event.GetDate())))
	}

	if event.Place != nil {
		t.AddNode(gedcom.NewPlaceNode(event.GetPlace()))
	}

	if event.Src != nil {
		t.AddNode(gedcom.NewSourceNode(event.GetSrc(), ""))
	}

	return t
}

func getMarriageEvent(family *api.Family) []gedcom.Node {
	var nodes []gedcom.Node

	if family.GetMarriageType() != api.MarriageType_NOT_MARRIED {
		t := gedcom.NewNode(mapMarriageTypeTagName[family.GetMarriageType()], "", "")

		if family.MarriageDate != nil {
			t.AddNode(gedcom.NewDateNode(getDate(family.GetMarriageDate())))
		}

		if family.MarriagePlace != nil {
			t.AddNode(gedcom.NewPlaceNode(family.GetMarriagePlace()))
		}

		if family.MarriageSrc != nil {
			t.AddNode(gedcom.NewSourceNode(family.GetMarriageSrc(), ""))
		}

		nodes = append(nodes, t)
	}

	if family.GetDivorceType() == api.DivorceType_DIVORCED {
		t := gedcom.NewNode(gedcom.TagDivorce, "", "")

		if family.DivorceDate != nil {
			t.AddNode(gedcom.NewDateNode(getDate(family.GetDivorceDate())))
		}

		nodes = append(nodes, t)
	}

	return nodes
}

func createFullIndividualNodesAndFamilyNodes( //nolint:funlen,gocognit,gocyclo,cyclop
	persons []*api.Person,
	personsNotes [][]utils.NoteWithTag,
	doc,
	docFamilies *gedcom.Document) {
	for i := 0; i < len(persons); i++ {
		indiNode := doc.AddIndividual(utils.PointerStr("I", persons[i].GetIndex()+1))
		indiNode.AddNode(getName(persons[i]))

		if persons[i].Sex != nil {
			if sex := getSex(persons[i].GetSex()); sex != "U" {
				indiNode.SetSex(sex)
			}
		}

		if len(persons[i].Aliases) > 0 {
			indiNode.AddName(strings.Join(persons[i].GetAliases(), ","))
		}

		if len(persons[i].Qualifiers) > 0 {
			indiNode.AddNode(gedcom.NewNode(gedcom.TagNickname, strings.Join(persons[i].GetQualifiers(), ","), ""))
		}

		if len(persons[i].SurnameAliases) > 0 {
			indiNode.AddNode(gedcom.NewNode(gedcom.TagSurname, strings.Join(persons[i].GetSurnameAliases(), ","), ""))
		}

		if persons[i].Occupation != nil {
			indiNode.AddNode(gedcom.NewNode(gedcom.TagOccupation, persons[i].GetOccupation(), ""))
		}

		if persons[i].Psources != nil {
			indiNode.AddNode(gedcom.NewNode(gedcom.TagSource, persons[i].GetPsources(), ""))
		}

		for _, title := range persons[i].GetTitles() {
			if t := getTitle(title); t != nil {
				indiNode.AddNode(t)
			}
		}

		if persons[i].Parents != nil {
			familyID := utils.PointerStr("F", persons[i].GetParents())
			indiNode.AddNode(gedcom.NewNode(gedcom.TagFamilyChild, "@"+familyID+"@", ""))

			if docFamilies.Families().ByPointer(familyID) == nil {
				docFamilies.AddFamily(familyID)
			}
		}

		for _, family := range persons[i].GetFamilies() {
			familyID := utils.PointerStr("F", family)
			indiNode.AddNode(gedcom.NewNode(gedcom.TagFamilySpouse, "@"+familyID+"@", ""))

			if docFamilies.Families().ByPointer(familyID) == nil {
				docFamilies.AddFamily(familyID)
			}
		}

		for _, event := range persons[i].GetEvents() {
			// events >= 50 are related to families, not individuals
			if event.GetName() < api.EventName_EFAM_MARRIAGE {
				if e := getEvent(event); e != nil && e.Tag().IsKnown() {
					indiNode.AddNode(e)
				}
			}
		}

		if personsNotes[i] != nil {
			noteNode := gedcom.NewNode(personsNotes[i][0].GetTag(), personsNotes[i][0].GetNote(), "")

			if len(personsNotes[i]) > 0 {
				for _, note := range personsNotes[i][1:] {
					noteNode.AddNode(gedcom.NewNode(note.GetTag(), note.GetNote(), ""))
				}
			}

			indiNode.AddNode(noteNode)
		}
	}
}

func fillFamilies( //nolint:cyclop
	families []*api.Family,
	familiesNotes [][]utils.NoteWithTag,
	doc *gedcom.Document) error {
	for k, familyNode := range doc.Families() {
		familyIdx, err := strconv.ParseInt(familyNode.Pointer()[1:], utils.ConstDecBase, 0)
		if err != nil {
			return fmt.Errorf("could not parse family ID: %w", err)
		}

		family := families[familyIdx]

		for _, event := range getMarriageEvent(family) {
			if event.Tag().IsKnown() {
				familyNode.AddNode(event)
			}
		}

		if family.Fsources != nil {
			familyNode.AddNode(gedcom.NewSourceNode(family.GetFsources(), ""))
		}

		if family.Father != nil {
			husbandID := utils.PointerStr("I", family.GetFather()+1)
			familyNode.SetHusband(doc.Individuals().ByPointer(husbandID))
		}

		if family.Mother != nil {
			wifeID := utils.PointerStr("I", family.GetMother()+1)
			familyNode.SetWife(doc.Individuals().ByPointer(wifeID))
		}

		for _, child := range family.GetChildren() {
			childID := utils.PointerStr("I", child+1)
			familyNode.AddChild(doc.Individuals().ByPointer(childID))
		}

		if familiesNotes[k] != nil {
			noteNode := gedcom.NewNode(familiesNotes[k][0].GetTag(), familiesNotes[k][0].GetNote(), "")

			if len(familiesNotes[k]) > 0 {
				for _, note := range familiesNotes[k][1:] {
					noteNode.AddNode(gedcom.NewNode(note.GetTag(), note.GetNote(), ""))
				}
			}

			familyNode.AddNode(noteNode)
		}
	}

	return nil
}

func (g *GenGedcom) Write(
	name string,
	persons []*api.Person,
	families []*api.Family,
	personsNotes, familiesNotes [][]utils.NoteWithTag) error {
	doc := getEmptyDocument(name)
	docFamilies := gedcom.NewDocument()

	createFullIndividualNodesAndFamilyNodes(persons, personsNotes, doc, docFamilies)

	for _, fam := range docFamilies.Nodes() {
		doc.AddNode(fam)
	}

	if err := fillFamilies(families, familiesNotes, doc); err != nil {
		return fmt.Errorf("failed to fill family nodes: %w", err)
	}

	doc.AddNode(gedcom.NewNode(gedcom.TagTrailer, "", ""))

	// log.Printf("%+v", doc)

	f, err := os.Create("output/" + name + ".ged")
	if err != nil {
		return fmt.Errorf("could not open gedcom file for writing: %w", err)
	}

	defer f.Close()

	enc := gedcom.NewEncoder(f, doc)
	if err = enc.Encode(); err != nil {
		return fmt.Errorf("error writing gedcom file: %w", err)
	}

	return nil
}

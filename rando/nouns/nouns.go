package nouns

import (
    "errors"
    "fmt"
    "github.com/arch4ngel/go-commoners/rando/shared"
    "golang.org/x/exp/slices"
)

const (
    ProperNounType NounType = "proper"
    CommonNounType NounType = "common"
)

var (
    nounTypes = []NounType{ProperNounType, CommonNounType}
    Values    = values{}
)

type NounType string

// values is a structure that binds noun types together.
type values struct {
    // Proper values.
    Proper []string
    // Common values.
    Common []string
}

// All aggregates each noun into a slice.
func All() (all *[]string) {
    all = new([]string)
    slices.Insert(*all, 0, Values.Proper...)
    slices.Insert(*all, 0, Values.Common...)
    return all
}

// Get returns a pointer to a noun, the type of which is
// determined by t.
func Get[S shared.StringType](t S) (v *string, err error) {
    var nT NounType
    if nT, err = ValidateType(t); err == nil {
        switch nT {
        case CommonNounType:
            v = shared.RandVal(&Values.Common)
        case ProperNounType:
            v = shared.RandVal(&Values.Proper)
        default:
            err = errors.New("noun type must be either common or proper")
        }
    }
    return v, err
}

// GetProper returns a pointer to a random proper noun.
func GetProper() (*string, error) {
    return Get(ProperNounType)
}

// GetCommon returns a pointer to a common noun.
func GetCommon() (*string, error) {
    return Get(CommonNounType)
}

func GetTypes() []NounType {
    return nounTypes[:]
}

func ValidType[S shared.StringType](nT S) bool {
    return slices.Index(nounTypes, NounType(nT)) != -1
}

func ValidateType[S shared.StringType](nT S) (NounType, error) {
    var err error
    if !ValidType(nT) {
        err = errors.New(fmt.Sprintf("invalid noun type: %v", nT))
    }
    return NounType(nT), err
}

package plate

import (
	"errors"
	"fmt"
)

type Plate struct {
	Alphabet    string
	StartNumber int8
	EndNumber   int8
	RegionCode  int8
}

const (
	Concatenated        = `12ب34599`
	Dashed              = `12-ب-345-99`
	Spaced              = `12 ب 345 99`
	NoAlphabet          = `1234599`
	NoAlphabetSpaced    = `12 345 99`
	NoAlphabetDashed    = `12-345-99`
	CodedAlphabet       = `120134599`
	CodedAlphabetDashed = `12-01-345-99`
	CodedAlphabetSpaced = `12 01 345 99`
)

var alphabetCodeMap = map[string]string{
	"الف":     "01",
	"ب":       "02",
	"پ":       "03",
	"ت":       "04",
	"ث":       "05",
	"ج":       "06",
	"چ":       "07",
	"ح":       "08",
	"خ":       "09",
	"د":       "10",
	"ذ":       "11",
	"ر":       "12",
	"ز":       "13",
	"ژ":       "14",
	"س":       "15",
	"ش":       "16",
	"ص":       "17",
	"ض":       "18",
	"ط":       "19",
	"ظ":       "20",
	"ع":       "21",
	"غ":       "22",
	"ف":       "23",
	"ق":       "24",
	"ک":       "25",
	"گ":       "26",
	"ل":       "27",
	"م":       "28",
	"ن":       "29",
	"و":       "30",
	"ه":       "31",
	"ی":       "32",
	"معلولین": "33",
	"تشریفات": "34",
}

func (p *Plate) Format(layout string) (string, error) {
	text := ""
	switch layout {
	case Concatenated:
		text = fmt.Sprintf("%v%v%v%v", p.StartNumber, p.Alphabet, p.EndNumber, p.RegionCode)
	case Dashed:
		text = fmt.Sprintf("%v-%v-%v-%v", p.StartNumber, p.Alphabet, p.EndNumber, p.RegionCode)
	case Spaced:
		text = fmt.Sprintf("%v %v %v %v", p.StartNumber, p.Alphabet, p.EndNumber, p.RegionCode)
	case NoAlphabet:
		text = fmt.Sprintf("%v%v%v", p.StartNumber, p.EndNumber, p.RegionCode)
	case NoAlphabetSpaced:
		text = fmt.Sprintf("%v %v %v", p.StartNumber, p.EndNumber, p.RegionCode)
	case NoAlphabetDashed:
		text = fmt.Sprintf("%v-%v-%v", p.StartNumber, p.EndNumber, p.RegionCode)
	case CodedAlphabet:
		text = fmt.Sprintf("%v%v%v%v", p.StartNumber, alphabetCodeMap[p.Alphabet], p.EndNumber, p.RegionCode)
	case CodedAlphabetSpaced:
		text = fmt.Sprintf("%v %v %v %v", p.StartNumber, alphabetCodeMap[p.Alphabet], p.EndNumber, p.RegionCode)
	case CodedAlphabetDashed:
		text = fmt.Sprintf("%v-%v-%v-%v", p.StartNumber, alphabetCodeMap[p.Alphabet], p.EndNumber, p.RegionCode)
	default:
		return "", errors.New("invalid layout")
	}
	return text, nil
}

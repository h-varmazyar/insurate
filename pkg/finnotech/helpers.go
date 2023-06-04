package finnotech

import (
	"fmt"
	"github.com/h-varmazyar/insurate/entity"
)

func generatePlateCode(plate *entity.Plate) string {
	alphabetCode := ""
	switch plate.Alphabet {
	case "الف":
		alphabetCode = "01"
	case "ب":
		alphabetCode = "02"
	case "پ":
		alphabetCode = "03"
	case "ت":
		alphabetCode = "04"
	case "ث":
		alphabetCode = "05"
	case "ج":
		alphabetCode = "06"
	case "چ":
		alphabetCode = "07"
	case "ح":
		alphabetCode = "08"
	case "خ":
		alphabetCode = "09"
	case "د":
		alphabetCode = "10"
	case "ذ":
		alphabetCode = "11"
	case "ر":
		alphabetCode = "12"
	case "ز":
		alphabetCode = "13"
	case "ژ":
		alphabetCode = "14"
	case "س":
		alphabetCode = "15"
	case "ش":
		alphabetCode = "16"
	case "ص":
		alphabetCode = "17"
	case "ض":
		alphabetCode = "18"
	case "ط":
		alphabetCode = "19"
	case "ظ":
		alphabetCode = "20"
	case "ع":
		alphabetCode = "21"
	case "غ":
		alphabetCode = "22"
	case "ف":
		alphabetCode = "23"
	case "ق":
		alphabetCode = "24"
	case "ک":
		alphabetCode = "25"
	case "گ":
		alphabetCode = "26"
	case "ل":
		alphabetCode = "27"
	case "م":
		alphabetCode = "28"
	case "ن":
		alphabetCode = "29"
	case "و":
		alphabetCode = "30"
	case "ه":
		alphabetCode = "31"
	case "ی":
		alphabetCode = "32"
	case "معلولین":
		alphabetCode = "33"
	case "تشریفات":
		alphabetCode = "34"
	default:
		return ""
	}
	return fmt.Sprintf("%v%v%v%v", plate.StartNumber, alphabetCode, plate.EndNumber, plate.RegionCode)
}

package startups

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

func registerTag() {
	govalidator.ParamTagMap["gt"] = func(value string, params ...string) bool {
		if len(params) == 1 {
			valueNum, _ := govalidator.ToFloat(value)
			minNum, _ := govalidator.ToFloat(params[0])
			return valueNum > minNum
		}
		return false
	}
	govalidator.ParamTagRegexMap["gt"] = regexp.MustCompile("^gt\\((\\d+)\\)$")
}

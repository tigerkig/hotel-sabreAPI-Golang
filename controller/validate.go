package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/context"
)

// ValidateExists - Validates all elements exists in map. It might be null or blank but key is exists
func ValidateExists(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		}
	}
	return flag
}

// ValidateNotNull - Validates all elements exists in map. Also it should not be blank
func ValidateNotNull(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			if data[col] == "" {
				flag = 0
			}
		}
	}
	return flag
}

// ValidateDateIfExists - Validates all elements exists in map as well as it should be in date format
func ValidateDateIfExists(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			if data[col] != "" {
				_, err := time.Parse("2006-01-02 15:04:05", data[col].(string))
				if err != nil {
					flag = 0
				}
			}
		}
	}
	return flag
}

// ValidateDateOnlyIfExists - Validates date only if exists
func ValidateDateOnlyIfExists(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			if data[col] != "" {
				_, err := time.Parse("2006-01-02", data[col].(string))
				if err != nil {
					flag = 0
				}
			}
		}
	}
	return flag
}

// ValidateDateOnly - Validates date strictly
func ValidateDateOnly(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; ok {
			_, err := time.Parse("2006-01-02", data[col].(string))
			if err != nil {
				flag = 0
			}
		} else {
			flag = 0
		}
	}
	return flag
}

// ValidateEmail - Validates string type is email or not
func ValidateEmail(email string) int {
	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(email))
	if matched == false {
		return 0
	}
	return 1
}

// ValidateMaxLength - Validates string's max length.
func ValidateMaxLength(str string, length int) int {
	if len(str) <= length {
		return 1
	}
	return 0
}

// ValidateMinLength - Validates string's min length.
func ValidateMinLength(str string, length int) int {
	if len(str) >= length {
		return 1
	}
	return 0
}

// ValidateNotNullAndFloat - Validates all elements exists in map. Also it should not be blank as well as it can be assertion as float64
func ValidateNotNullAndFloat(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			_, ok := data[col].(float64)
			if data[col] == "" || !ok {
				flag = 0
			}
		}
	}
	return flag
}

// ValidateNotNullAndString - Validates all elements exists in map. Also it should not be blank as well as it can be assertion as string
func ValidateNotNullAndString(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			_, ok := data[col].(string)
			if data[col] == "" || !ok {
				flag = 0
			}
		}
	}
	return flag
}

// ValidateNotNullStructString - Validates all elements exists in struct and check its not null.
func ValidateNotNullStructString(params ...string) int {
	flag := 1
	for _, col := range params {
		if col == "" {
			flag = 0
		}
	}
	return flag
}

// ValidateNotNullStructFloat - Validates all elements exists in struct and check its not null.
func ValidateNotNullStructFloat(params ...float64) int {
	flag := 1
	for _, col := range params {
		if fmt.Sprintf("%f", col) == "" {
			flag = 0
		}
	}
	return flag
}

// ValidateNotNullAndLangArr - Validates language array and default language value is not null
func ValidateNotNullAndLangArr(r *http.Request, data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			langMap, ok := data[col].(map[string]interface{})
			if !ok {
				flag = 0
			} else {
				lang := context.Get(r, "lang").(string)
				_, ok = langMap[lang].(string)
				if langMap[lang].(string) == "" || !ok {
					flag = 0
				}
			}
		}
	}
	return flag
}

// ValidateExistsAndLangArr - Validates language array and default language value is exists
func ValidateExistsAndLangArr(r *http.Request, data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			langMap, ok := data[col].(map[string]interface{})
			if !ok {
				flag = 0
			} else {
				lang := context.Get(r, "lang").(string)
				_, ok = langMap[lang].(string)
				if !ok {
					flag = 0
				}
			}
		}
	}
	return flag
}

// ValidateTimeOnly - Validates time strictly
func ValidateTimeOnly(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; ok {
			_, err := time.Parse("15:04:05", data[col].(string))
			if err != nil {
				flag = 0
			}
		} else {
			flag = 0
		}
	}
	return flag
}

// ValidateNotNullStrArr - Validates string array Created By Harshit
func ValidateNotNullStrArr(data map[string]interface{}, validates []string) int {
	flag := 1
	for _, col := range validates {
		if _, ok := data[col]; !ok {
			flag = 0
		} else {
			_, ok := data[col].(string)
			if data[col] == "" || !ok {
				flag = 0
			}
		}
	}
	return flag
}

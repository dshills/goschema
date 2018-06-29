package main

import (
	"strings"
	"unicode"
)

func numToWord(str string) string {
	if !unicode.IsDigit(rune(str[0])) {
		return str
	}
	ns := ""
	for i, r := range str {
		if unicode.IsDigit(r) {
			switch r {
			case '0':
				ns += "zero"
			case '1':
				ns += "one"
			case '2':
				ns += "two"
			case '3':
				ns += "three"
			case '4':
				ns += "four"
			case '5':
				ns += "five"
			case '6':
				ns += "six"
			case '7':
				ns += "seven"
			case '8':
				ns += "eight"
			case '9':
				ns += "nine"
			}
		} else {
			ns += str[i:]
			return ns
		}
	}
	return ns
}

func goName(str string) string {
	str = strings.ToLower(str)
	str = numToWord(str)
	if str == "dob" {
		str = "DOB"
	}
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Replace(str, " id", "ID", -1)
	str = strings.Replace(str, "url", "URL", -1)
	str = strings.Replace(str, "uid", "UID", -1)
	str = strings.Replace(str, " api", "API", -1)
	str = strings.Replace(str, " ip", "IP", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

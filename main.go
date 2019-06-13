package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const profile = "data/profile.json"

type Section struct {
	Length int `json:"length"`
	CharacterSet string `json:"character-set"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	Content string `json:"content"`
}

type Profile struct {
	Count int `json:"count"`
	Separator string `json:"separator"`
	Sections []*Section `json:"sections"`
	Output string `json:"output"`
	Characters map[string]string `json:"characters"`
}

func main(){
	file, err := os.Open(profile)
	if nil != err {
		panic(err.Error())
	}
	var settingsPtr = new(*Profile)
	err = json.NewDecoder(file).Decode(&settingsPtr)
	var settings = *settingsPtr
	file.Close()
	count := settings.Count
	separator := settings.Separator
	characters := settings.Characters
	outputFile, err := os.Create(settings.Output)
	if nil != err {
		panic(err)
	}
	defer outputFile.Close()
	codes := make(map[string]int)
	var charsets = make([]string, len(settings.Sections))
	for ;len(codes) < count; {
		var code string
		for idx, section := range settings.Sections {
			if section.Content == "" {
				code += section.Prefix
				charset := charsets[idx]
				if charset == "" {
					chrs := strings.Split(section.CharacterSet, "-")
					for _, chrn := range chrs {
						chr,ok := characters[chrn]
						if ok {
							charset += chr
						} else {
							fmt.Printf("ERROR: character set %v not provided!!\n", chrn)
						}
					}
					charsets[idx] = charset
				}
				code += RandStr(
					charset,
					section.Length - len(section.Prefix) - len(section.Suffix))
				code += section.Suffix
			} else {
				code += section.Content
			}
			if idx != (len(settings.Sections) - 1) {
				code += separator
			}
		}
		fmt.Println(code)
		codes[code] = 1
	}
	for idx := range codes {
		outputFile.WriteString(idx)
		outputFile.WriteString("\r\n")
	}
}

func RandStr(charset string, length int) string  {
	scope := len(charset)
	if 0 == scope {
		return ""
	}
	var res strings.Builder
	for ;res.Len() < length; {
		nextIdx := rand.Intn(scope)
		res.WriteByte(charset[nextIdx])
	}
	return res.String()
}
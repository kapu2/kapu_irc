package main

import (
	"fmt"
)

type Tag struct {
	clientPrefix bool   // optional
	vendor       []rune // optional, not sure what it does
	key          []rune
	escapedValue []rune // optional

}
type Source struct {
	sourceName []rune // nickname or sourcename, only nickname may have user or host fields
	user       []rune // optional
	host       []rune // optional
}

type IRCMessage struct {
	tags       []Tag
	source     Source
	command    []rune
	parameters [][]rune
}

func FindRunePos(line []rune, toFind rune) int {
	pos := -1
	for i, r := range line {
		if toFind == r {
			pos = i
			break
		}
	}
	if pos != -1 {
		return pos
	} else {
		return pos
	}
}

func ParseIRCMessage(line []rune) (IRCMessage, error) {
	msg := IRCMessage{}
	var tags []rune
	var source []rune
	var commandAndParameters []rune

	if len(line) > 0 && line[0] == '@' {
		pos := FindRunePos(line, ' ')
		tags = line[0:pos]
		// pos + 1, so we eat away space
		line = line[pos+1:]

	}
	if len(line) > 0 && line[0] == ':' {
		pos := FindRunePos(line, ' ')
		source = line[0:pos]
		// pos + 1, so we eat away space
		line = line[pos+1:]

	}
	if len(line) > 0 {
		commandAndParameters = line
	}
	if len(tags) > 0 {
		parsedTags, err := ParseTags(tags)
		if err != nil {
			fmt.Print(err)
		}
		msg.tags = parsedTags
	}
	if len(source) > 0 {
		parsedSource, err := ParseSource(source)
		if err != nil {
			fmt.Print(err)
		}
		msg.source = parsedSource
	}

	splitCommandAndParameters := Split(commandAndParameters, ' ')
	msg.command = splitCommandAndParameters[0]

	var parameters []rune
	for i, v := range splitCommandAndParameters {
		if i != 0 {
			parameters = append(parameters, v...)
			parameters = append(parameters, ' ')
		}
	}
	parameters = parameters[0 : len(parameters)-1] // remove trailing space
	msg.parameters = ParseParameters(parameters)
	return msg, nil
}

func Split[K comparable](toSplit []K, delimiter K) [][]K {
	var ret [][]K
	start := 0
	i := 0
	var v K
	for i, v = range toSplit {
		if v == delimiter {
			ret = append(ret, toSplit[start:i])
			start = i + 1
		}
	}
	if len(toSplit) > 0 {
		remaining := toSplit[start : i+1]
		if len(remaining) > 0 {
			ret = append(ret, remaining)
		}
	}

	return ret
}

func MultiDelimiterSplit[K comparable](toSplit []K, delimiters []K) [][]K {
	var ret [][]K
	start := 0
	i := 0
	var v K
	for i, v = range toSplit {
		for _, delimiter := range delimiters {
			if v == delimiter {
				ret = append(ret, toSplit[start:i])
				start = i + 1
				break
			}
		}
	}
	remaining := toSplit[start : i+1]
	if len(remaining) > 0 {
		ret = append(ret, remaining)
	}
	return ret
}

func ParseParameters(parameters []rune) [][]rune {
	prevWhiteSpace := true
	addFinalParam := false
	var finalParam []rune
	for i, v := range parameters {
		if v == ' ' {
			prevWhiteSpace = true
		} else if v == ':' && prevWhiteSpace {
			if i != len(parameters)-1 && i != 0 {
				finalParam = parameters[i+1:]
				// removing the : and final parameter
				parameters = parameters[0:i]
				addFinalParam = true
				break
			} else if i != len(parameters)-1 {
				finalParam = parameters[i+1:]
				// we only have final parameter, so clear parameters
				parameters = []rune("")
				addFinalParam = true
				break
			} else {
				finalParam = []rune("")
				// removing final lonely :
				parameters = parameters[0:i]
				addFinalParam = false
				break
			}
		} else {
			prevWhiteSpace = false
		}
	}
	parsedParams := Split(parameters, ' ')
	if addFinalParam {
		parsedParams = append(parsedParams, finalParam)
	}
	return parsedParams
}

func ParseSource(source []rune) (Source, error) {
	var ret Source
	if len(source) > 1 && source[0] == ':' {
		// eat away :
		source = source[1:]
	} else {
		return ret, fmt.Errorf("error: source only consisting of : and nothing else")
	}
	sourceNameUserAndHost := Split(source, '@')
	if len(sourceNameUserAndHost) > 1 {
		ret.host = sourceNameUserAndHost[1]
	}
	sourceNameAndUser := Split(sourceNameUserAndHost[0], '!')
	ret.sourceName = sourceNameAndUser[0]
	if len(sourceNameAndUser) == 2 {
		ret.user = sourceNameAndUser[1]
	}
	return ret, nil
}

func ParseTag(tag []rune) Tag {
	ret := Tag{}

	if len(tag) > 0 && tag[0] == '+' {
		ret.clientPrefix = true
		tag = tag[1:]
	}
	vendorKeyAndVal := Split(tag, '=')
	vendorKey := vendorKeyAndVal[0]
	if len(vendorKeyAndVal) == 2 {
		ret.escapedValue = vendorKeyAndVal[1]
	} else if FindRunePos(tag, '=') != -1 {
		ret.escapedValue = []rune("")
	}
	vendorAndKey := Split(vendorKey, '/')
	if len(vendorAndKey) == 2 {
		ret.vendor = vendorAndKey[0]
		ret.key = vendorAndKey[1]
	} else {
		ret.key = vendorAndKey[0]
	}
	return ret
}
func ParseTags(tags []rune) ([]Tag, error) {
	var ret []Tag
	if len(tags) > 1 && tags[0] == '@' {
		// eat away @
		tags = tags[1:]
	} else {
		return ret, fmt.Errorf("error: tag only consisting of @ and nothing else")
	}
	tagList := Split(tags, ';')

	for _, v := range tagList {
		ret = append(ret, ParseTag(v))
	}
	return ret, nil
}

func IRCMessageToString(msg IRCMessage) []rune {
	var ret []rune
	tags := TagsToString(msg.tags)
	source := SourceToString(msg.source)
	command := msg.command
	params := ParametersToString(msg.parameters)
	ret = append(ret, tags...)
	ret = append(ret, source...)
	ret = append(ret, command...)
	ret = append(ret, []rune(" ")...)
	ret = append(ret, params...)
	return ret
}

func TagsToString(tags []Tag) []rune {
	var ret []rune
	if len(tags) == 0 {
		return ret
	}
	ret = append(ret, []rune("@")...)
	for i, tag := range tags {
		if i != 0 {
			ret = append(ret, []rune(";")...)
		}
		if tag.clientPrefix {
			ret = append(ret, []rune("+")...)
		}
		if tag.vendor != nil {
			ret = append(ret, tag.vendor...)
			ret = append(ret, []rune("/")...)
		}
		ret = append(ret, tag.key...)
		if tag.escapedValue != nil {
			ret = append(ret, []rune("=")...)
			ret = append(ret, tag.escapedValue...)
		}
	}
	ret = append(ret, []rune(" ")...)
	return ret
}

func SourceToString(source Source) []rune {
	var ret []rune
	if source.sourceName == nil {
		return ret
	}
	ret = append(ret, []rune(":")...)
	ret = append(ret, source.sourceName...)
	if source.user != nil {
		ret = append(ret, []rune("!")...)
		ret = append(ret, source.user...)
	}
	if source.host != nil {
		ret = append(ret, []rune("@")...)
		ret = append(ret, source.host...)
	}
	ret = append(ret, []rune(" ")...)
	return ret
}

func ParametersToString(parameters [][]rune) []rune {
	var ret []rune
	for i, param := range parameters {
		if i != 0 {
			ret = append(ret, []rune(" ")...)
		}
		if i == len(parameters)-1 {
			ret = append(ret, []rune(":")...)
		}
		ret = append(ret, param...)
	}
	return ret
}

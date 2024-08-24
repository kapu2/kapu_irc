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

func FindRunePos(line []rune, toFind rune) (int, error) {
	pos := -1
	for i, r := range line {
		if toFind == r {
			pos = i
			break
		}
	}
	if pos != -1 {
		return pos, nil
	} else {
		return pos, fmt.Errorf("failed to find rune position, toFind: %v line: %v", toFind, line)
	}
}

func ParseIRCMessage(line []rune) (IRCMessage, error) {
	msg := IRCMessage{}
	var tags []rune
	var source []rune
	var commandAndParameters []rune

	if len(line) > 0 && line[0] == '@' {
		pos, err := FindRunePos(line, ' ')
		if err == nil {
			tags = line[0:pos]
			// pos + 1, so we eat away space
			line = line[pos+1:]
		} else {
			return msg, fmt.Errorf("error: malformed IRCMessage: %s", string(line))
		}
	}
	if len(line) > 0 && line[0] == ':' {
		pos, err := FindRunePos(line, ' ')
		if err == nil {
			source = line[0:pos]
			// pos + 1, so we eat away space
			line = line[pos+1:]
		} else {
			return msg, fmt.Errorf("error: malformed IRCMessage: %s", string(line))
		}
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
		parsedSource := ParseSource(source)
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
			if i != len(parameters)-1 {
				finalParam = parameters[i+1:]
				// removing the : and final parameter
				parameters = parameters[0:i]
				addFinalParam = true
				break
			} else {
				finalParam = append(finalParam, []rune("")...)
				// removing the : and final parameter
				parameters = parameters[0:i]
				addFinalParam = true
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

func ParseSource(source []rune) Source {
	var ret Source
	sourceNameAndUserHost := Split(source, '!')
	ret.sourceName = sourceNameAndUserHost[0]
	if len(sourceNameAndUserHost) == 2 {
		userAndHost := Split(sourceNameAndUserHost[1], '@')
		ret.user = userAndHost[0]
		if len(userAndHost) == 2 {
			ret.host = userAndHost[1]
		}
	}
	return ret
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
		return ret, fmt.Errorf("error: tag only consisting of @")
	}
	tagList := Split(tags, ';')

	for _, v := range tagList {
		ret = append(ret, ParseTag(v))
	}
	return ret, nil
}

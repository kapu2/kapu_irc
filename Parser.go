package main

import (
	"fmt"
	"strings"
)

type Tag struct {
	clientPrefix bool   // optional
	vendor       string // optional, not sure what it does
	key          string
	escapedValue string // optional

}
type Source struct {
	sourceName string // nickname or sourcename, only nickname may have user or host fields
	user       string // optional
	host       string // optional
}

type IRCMessage struct {
	tags       []Tag
	source     Source
	command    string
	parameters []string
}

func (irc *IRCMessage) AddParameter(param string) {
	irc.parameters = append(irc.parameters, param)
}

func ParseIRCMessage(line string) (IRCMessage, error) {
	msg := IRCMessage{}
	var tags string
	var source string
	var commandAndParameters string
	var err error
	if len(line) > 0 && strings.IndexRune(line, '@') == 0 {
		pos := strings.IndexRune(line, ' ')
		tags = line[0:pos]

		line = line[pos:]
		// eat away space
		line = RemoveFirstRuneFromString(line)

	}
	if len(line) > 0 && strings.IndexRune(line, ':') == 0 {
		pos := strings.IndexRune(line, ' ')
		source = line[0:pos]
		// eat away space
		line = line[pos:]
		line = RemoveFirstRuneFromString(line)

	}
	if len(line) > 0 {
		commandAndParameters = line
	}
	if len(tags) > 0 {
		var parsedTags []Tag
		parsedTags, err = ParseTags(tags)
		if err != nil {
			fmt.Print(err)
		}
		msg.tags = parsedTags
	}
	if len(source) > 0 {
		var parsedSource Source
		parsedSource, err = ParseSource(source)
		if err != nil {
			fmt.Print(err)
		}
		msg.source = parsedSource
	}

	splitCommandAndParameters := strings.Split(commandAndParameters, " ")
	msg.command = splitCommandAndParameters[0]

	var parameters string
	for i, v := range splitCommandAndParameters {
		if i != 0 {
			parameters += v + " "
		}
	}
	parameters = parameters[0 : len(parameters)-1] // remove trailing space
	msg.parameters = ParseParameters([]rune(parameters))
	return msg, err
}

func ParseParameters(parameters []rune) []string {
	prevWhiteSpace := true
	addFinalParam := false
	var finalParam string
	for i, v := range parameters {
		if v == ' ' {
			prevWhiteSpace = true
		} else if v == ':' && prevWhiteSpace {
			if i != len(parameters)-1 && i != 0 {
				finalParam = string(parameters[i+1:])
				// removing the : and final parameter
				parameters = parameters[0:i]
				addFinalParam = true
				break
			} else if i != len(parameters)-1 {
				finalParam = string(parameters[i+1:])
				// we only have final parameter, so clear parameters
				parameters = []rune("")
				addFinalParam = true
				break
			} else {
				finalParam = string("")
				// removing final lonely :
				parameters = parameters[0:i]
				addFinalParam = false
				break
			}
		} else {
			prevWhiteSpace = false
		}
	}

	var parsedParams []string
	if len(parameters) != 0 {
		// trimspace, because we may have one extra space at the end
		parsedParams = strings.Split(strings.TrimSpace(string(parameters)), " ")
	}
	if addFinalParam {
		parsedParams = append(parsedParams, finalParam)
	}
	return parsedParams
}

func ParseSource(source string) (Source, error) {
	var ret Source
	if len(source) > 1 && source[0] == ':' {
		// eat away :
		source = source[1:]
	} else {
		return ret, fmt.Errorf("error: source only consisting of : and nothing else")
	}
	sourceNameUserAndHost := strings.Split(source, "@")
	if len(sourceNameUserAndHost) > 1 {
		ret.host = sourceNameUserAndHost[1]
	}
	sourceNameAndUser := strings.Split(sourceNameUserAndHost[0], "!")
	ret.sourceName = sourceNameAndUser[0]
	if len(sourceNameAndUser) == 2 {
		ret.user = sourceNameAndUser[1]
	}
	return ret, nil
}

func ParseTag(tag string) Tag {
	ret := Tag{}

	if len(tag) > 0 && tag[0] == '+' {
		ret.clientPrefix = true
		tag = tag[1:]
	}
	vendorKeyAndVal := strings.Split(tag, "=")
	vendorKey := vendorKeyAndVal[0]
	if len(vendorKeyAndVal) == 2 {
		ret.escapedValue = vendorKeyAndVal[1]
	} else if strings.ContainsRune(tag, '=') {
		ret.escapedValue = string("")
	}
	vendorAndKey := strings.Split(vendorKey, "/")
	if len(vendorAndKey) == 2 {
		ret.vendor = vendorAndKey[0]
		ret.key = vendorAndKey[1]
	} else {
		ret.key = vendorAndKey[0]
	}
	return ret
}
func ParseTags(tags string) ([]Tag, error) {
	var ret []Tag
	if len(tags) > 1 && tags[0] == '@' {
		// eat away @
		tags = RemoveFirstRuneFromString(tags)
	} else {
		return ret, fmt.Errorf("error: tag only consisting of @ and nothing else")
	}
	tagList := strings.Split(tags, ";")

	for _, v := range tagList {
		ret = append(ret, ParseTag(v))
	}
	return ret, nil
}

func IRCMessageToString(msg IRCMessage) string {
	// in separate function for easy testing
	ret := IRCMessageToStringWithoutNewline(msg)
	ret += "\r\n"
	return ret
}

func IRCMessageToStringWithoutNewline(msg IRCMessage) string {
	var ret string
	tags := TagsToString(msg.tags)
	source := SourceToString(msg.source)
	command := msg.command
	params := ParametersToString(msg.parameters)
	ret += tags + source + command + " " + params

	return ret
}

func TagsToString(tags []Tag) string {
	var ret string
	if len(tags) == 0 {
		return ret
	}
	ret += "@"
	for i, tag := range tags {
		if i != 0 {
			ret += ";"
		}
		if tag.clientPrefix {
			ret += "+"
		}
		if tag.vendor != "" {
			ret += tag.vendor + "/"
		}
		ret += tag.key
		// Technically value can be empty and still have =, but we will not send = without a value
		if tag.escapedValue != "" {
			ret += "=" + tag.escapedValue
		}
	}
	ret += " "
	return ret
}

func SourceToString(source Source) string {
	var ret string
	if source.sourceName == "" {
		return ret
	}
	ret += ":" + source.sourceName
	if source.user != "" {
		ret += "!" + source.user
	}
	if source.host != "" {
		ret += "@" + source.host
	}
	ret += " "
	return ret
}

func ParametersToString(parameters []string) string {
	var ret string
	for i, param := range parameters {
		if i != 0 {
			ret += " "
		}
		if i == len(parameters)-1 {
			ret += ":"
		}
		ret += param
	}
	return ret
}

package main

import (
	"testing"
)

func SpliceIsSame[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TwoDSpliceIsSame[T comparable](a [][]T, b [][]T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !SpliceIsSame(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestParseInput(t *testing.T) {
	v := NewTerminalView()
	m := NewStateKeeper()
	c := NewController(v, m)
	str := "/j #test-channel"
	want := "JOIN #test-channel\r\n"
	got := c.ParseInput(str)
	if !SpliceIsSame(got, []byte(want)) {
		t.Fatalf("expected: %s got: %s", want, string(got))
	}
}

func TestGetConnectionInfo(t *testing.T) {

	view := NewTerminalView()
	ipAndPort, nick := view.GetConnectionInfo()
	if ipAndPort == "" {
		t.Fatalf("error, ipAndPort empty")
	}

	if nick == "" {
		t.Fatalf("error, nick empty")
	}
}

func TestParseIRCMessageTags(t *testing.T) {
	msg := []rune("@url=;netsplit=tur,ty :irc.example.com CAP LS * :multi-prefix extended-join sasl")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := []Tag{{key: []rune("url")},
		{key: []rune("netsplit"), escapedValue: []rune("tur,ty")}}

	if len(parsedMsg.tags) != len(want) {
		t.Fatalf("error, unexpected amount of tags, want: %d got: %d", len(want), len(parsedMsg.tags))
	}
	for i, v := range parsedMsg.tags {
		if want[i].clientPrefix != v.clientPrefix {
			t.Fatalf("error, clientprefix differs, want: %t got: %t", want[i].clientPrefix, v.clientPrefix)
		}
		if !SpliceIsSame(want[i].escapedValue, parsedMsg.tags[i].escapedValue) {
			t.Fatalf("error, escapedValue differs, want: %s got: %s", string(want[i].escapedValue), string(v.escapedValue))
		}
		if !SpliceIsSame(want[i].key, parsedMsg.tags[i].key) {
			t.Fatalf("error, key differs, want: %s got: %s", string(want[i].key), string(v.key))
		}
		if !SpliceIsSame(want[i].vendor, parsedMsg.tags[i].vendor) {
			t.Fatalf("error, vendor differs, want: %s got: %s", string(want[i].vendor), string(v.vendor))
		}
	}
	got := IRCMessageToStringWithoutNewline(parsedMsg)
	if !SpliceIsSame(got, msg) {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}
}

func TestParseIRCMessageSource(t *testing.T) {
	msg := []rune(":nickname!user@host tauhkaa :parameters")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := Source{sourceName: []rune("nickname"), user: []rune("user"), host: []rune("host")}

	if !SpliceIsSame(want.sourceName, parsedMsg.source.sourceName) {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if !SpliceIsSame(want.user, parsedMsg.source.user) {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if !SpliceIsSame(want.host, parsedMsg.source.host) {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got := IRCMessageToStringWithoutNewline(parsedMsg)
	if !SpliceIsSame(got, msg) {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = []rune(":nickname@host tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: []rune("nickname"), user: nil, host: []rune("host")}
	if !SpliceIsSame(want.sourceName, parsedMsg.source.sourceName) {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if !SpliceIsSame(want.user, parsedMsg.source.user) {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if !SpliceIsSame(want.host, parsedMsg.source.host) {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if !SpliceIsSame(got, msg) {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = []rune(":nickname!user tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: []rune("nickname"), user: []rune("user"), host: nil}
	if !SpliceIsSame(want.sourceName, parsedMsg.source.sourceName) {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if !SpliceIsSame(want.user, parsedMsg.source.user) {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if !SpliceIsSame(want.host, parsedMsg.source.host) {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if !SpliceIsSame(got, msg) {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = []rune(":nickname tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: []rune("nickname"), user: nil, host: nil}
	if !SpliceIsSame(want.sourceName, parsedMsg.source.sourceName) {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if !SpliceIsSame(want.user, parsedMsg.source.user) {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if !SpliceIsSame(want.host, parsedMsg.source.host) {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if !SpliceIsSame(got, msg) {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}
}
func TestParseIRCMessageCommand(t *testing.T) {
	msg := []rune(":irc.example.com CAP * LIST :")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := []rune("CAP")
	if !SpliceIsSame(parsedMsg.command, want) {
		t.Fatalf("error, want: %s got: %s", string(want), string(parsedMsg.command))
	}
}

func TestParseIRCMessageParameters(t *testing.T) {
	msg := []rune("CAP * LS :multi-prefix sasl")
	parsedMsg, _ := ParseIRCMessage(msg)
	want2D := [][]rune{[]rune("*"), []rune("LS"), []rune("multi-prefix sasl")}
	if !TwoDSpliceIsSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = []rune("CAP REQ :sasl message-tags foo")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = [][]rune{[]rune("REQ"), []rune("sasl message-tags foo")}
	if !TwoDSpliceIsSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = []rune(":dan!d@localhost PRIVMSG #chan :Hey!")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = [][]rune{[]rune("#chan"), []rune("Hey!")}
	if !TwoDSpliceIsSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = []rune(":dan!d@localhost PRIVMSG #chan Hey!")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = [][]rune{[]rune("#chan"), []rune("Hey!")}
	if !TwoDSpliceIsSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = []rune(":dan!d@localhost PRIVMSG #chan ::-)")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = [][]rune{[]rune("#chan"), []rune(":-)")}
	if !TwoDSpliceIsSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}
}

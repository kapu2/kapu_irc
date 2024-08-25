package main

import (
	"testing"
)

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
	msg := string("@url;netsplit=tur,ty :irc.example.com CAP LS * :multi-prefix extended-join sasl")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := []Tag{{key: string("url")},
		{key: string("netsplit"), escapedValue: string("tur,ty")}}

	if len(parsedMsg.tags) != len(want) {
		t.Fatalf("error, unexpected amount of tags, want: %d got: %d", len(want), len(parsedMsg.tags))
	}
	for i, v := range parsedMsg.tags {
		if want[i].clientPrefix != v.clientPrefix {
			t.Fatalf("error, clientprefix differs, want: %t got: %t", want[i].clientPrefix, v.clientPrefix)
		}
		if want[i].escapedValue != parsedMsg.tags[i].escapedValue {
			t.Fatalf("error, escapedValue differs, want: %s got: %s", string(want[i].escapedValue), string(v.escapedValue))
		}
		if want[i].key != parsedMsg.tags[i].key {
			t.Fatalf("error, key differs, want: %s got: %s", string(want[i].key), string(v.key))
		}
		if want[i].vendor != parsedMsg.tags[i].vendor {
			t.Fatalf("error, vendor differs, want: %s got: %s", string(want[i].vendor), string(v.vendor))
		}
	}
	got := IRCMessageToStringWithoutNewline(parsedMsg)
	if got != msg {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}
}

func TestParseIRCMessageSource(t *testing.T) {
	msg := string(":nickname!user@host tauhkaa :parameters")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := Source{sourceName: string("nickname"), user: string("user"), host: string("host")}

	if want.sourceName != parsedMsg.source.sourceName {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if want.user != parsedMsg.source.user {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if want.host != parsedMsg.source.host {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got := IRCMessageToStringWithoutNewline(parsedMsg)
	if got != msg {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = string(":nickname@host tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: string("nickname"), user: "", host: string("host")}
	if want.sourceName != parsedMsg.source.sourceName {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if want.user != parsedMsg.source.user {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if want.host != parsedMsg.source.host {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if got != msg {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = string(":nickname!user tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: string("nickname"), user: string("user"), host: ""}
	if want.sourceName != parsedMsg.source.sourceName {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if want.user != parsedMsg.source.user {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if want.host != parsedMsg.source.host {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if got != msg {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}

	msg = string(":nickname tauhkaa :parameters")
	parsedMsg, _ = ParseIRCMessage(msg)
	want = Source{sourceName: string("nickname"), user: "", host: ""}
	if want.sourceName != parsedMsg.source.sourceName {
		t.Fatalf("error, sourceName differs, want: %s got: %s", string(want.sourceName), string(parsedMsg.source.sourceName))
	}
	if want.user != parsedMsg.source.user {
		t.Fatalf("error, user differs, want: %s got: %s", string(want.user), string(parsedMsg.source.user))
	}
	if want.host != parsedMsg.source.host {
		t.Fatalf("error, host differs, want: %s got: %s", string(want.host), string(parsedMsg.source.host))
	}

	got = IRCMessageToStringWithoutNewline(parsedMsg)
	if got != msg {
		t.Fatalf("error, parse and unparse failed, want: %s got: %s", string(msg), string(got))
	}
}
func TestParseIRCMessageCommand(t *testing.T) {
	msg := string(":irc.example.com CAP * LIST :")
	parsedMsg, _ := ParseIRCMessage(msg)
	want := string("CAP")
	if parsedMsg.command != want {
		t.Fatalf("error, want: %s got: %s", string(want), string(parsedMsg.command))
	}
}

func TestParseIRCMessageParameters(t *testing.T) {
	msg := string("CAP * LS :multi-prefix sasl")
	parsedMsg, _ := ParseIRCMessage(msg)
	want2D := []string{string("*"), string("LS"), string("multi-prefix sasl")}
	if !TwoDStringAreSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = string("CAP REQ :sasl message-tags foo")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = []string{string("REQ"), string("sasl message-tags foo")}
	if !TwoDStringAreSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = string(":dan!d@localhost PRIVMSG #chan :Hey!")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = []string{string("#chan"), string("Hey!")}
	if !TwoDStringAreSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = string(":dan!d@localhost PRIVMSG #chan Hey!")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = []string{string("#chan"), string("Hey!")}
	if !TwoDStringAreSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}

	msg = string(":dan!d@localhost PRIVMSG #chan ::-)")
	parsedMsg, _ = ParseIRCMessage(msg)
	want2D = []string{string("#chan"), string(":-)")}
	if !TwoDStringAreSame(parsedMsg.parameters, want2D) {
		t.Fatalf("error with parameters")
	}
}

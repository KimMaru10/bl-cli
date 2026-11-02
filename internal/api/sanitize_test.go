package api

import (
	"net/http"
	"strings"
	"testing"
)

func TestStripControlRemovesTerminalAndBidiControls(t *testing.T) {
	in := "a\x1b[2J\x1b]0;t\x07b\r‮c​d e\nf\tg 👨‍👩"
	want := "a[2J]0;tbcde\nf\tg 👨‍👩"
	if got := StripControl(in); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestDecodeSanitizesNestedStrings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"issueKey":"P-1","summary":"s\u001b[2J","status":{"name":"st\u001b]0;x\u0007"},` +
			`"attachments":[{"id":1,"name":"inv‮fdp.exe"}],` +
			`"customFields":[{"name":"f","value":"v\u001b[31m"}]}`))
	})
	is, err := c.GetIssue("P-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{is.Summary, is.Status.Name, is.Attachments[0].Name, is.CustomFields[0].DisplayValue()} {
		if strings.ContainsAny(s, "\x1b\x07‮") {
			t.Errorf("control char left in %q", s)
		}
	}
}

func TestAPIErrorMessageIsSanitized(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"message":"bad\u001b]52;c;x\u0007","code":7}]}`))
	})
	_, err := c.GetIssue("P-1")
	if err == nil || strings.ContainsAny(err.Error(), "\x1b\x07") {
		t.Fatalf("got %v", err)
	}
}

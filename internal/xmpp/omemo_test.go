package xmpp

import "testing"

func TestIsOMEMOMessage(t *testing.T) {
	fallback1 := "I sent you an OMEMO encrypted message but your client doesn’t seem to support that."
	if !IsOMEMOMessage(fallback1, false, false) {
		t.Fatalf("expected fallback1 to be recognized as OMEMO")
	}

	fallback2 := "I sent you an OMEMO encrypted message but your client doesn’t seem to support that. Find more information on https://conversations.im/omemo"
	if !IsOMEMOMessage(fallback2, false, false) {
		t.Fatalf("expected fallback2 to be recognized as OMEMO")
	}

	if !IsOMEMOMessage("", true, false) {
		t.Fatalf("expected hasEncrypted to be recognized as OMEMO")
	}

	if !IsOMEMOMessage("", false, true) {
		t.Fatalf("expected hasEME to be recognized as OMEMO")
	}

	normal := "Hello world"
	if IsOMEMOMessage(normal, false, false) {
		t.Fatalf("expected normal message not to be recognized as OMEMO")
	}
}

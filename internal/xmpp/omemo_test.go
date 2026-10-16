package xmpp

import (
	"strings"
	"testing"
)

func TestGetOrCreateDeviceID(t *testing.T) {
	id1 := GetOrCreateDeviceID()
	if id1 == 0 {
		t.Fatal("expected non-zero device ID")
	}
	id2 := GetOrCreateDeviceID()
	if id1 != id2 {
		t.Fatalf("expected persistent device ID %d, got %d", id1, id2)
	}
}

func TestInitOMEMO(t *testing.T) {
	devID := uint32(987654)
	om, err := InitOMEMO(devID)
	if err != nil {
		t.Fatalf("InitOMEMO error: %v", err)
	}
	if om.DeviceID != devID {
		t.Fatalf("expected device ID %d, got %d", devID, om.DeviceID)
	}
	if len(om.Fingerprint) == 0 {
		t.Fatal("expected non-empty identity fingerprint")
	}
	if len(om.Bundle.PreKeys) != 50 {
		t.Fatalf("expected 50 prekeys, got %d", len(om.Bundle.PreKeys))
	}
}

func TestOMEMOInfo(t *testing.T) {
	c := &Client{}
	info := c.OMEMOInfo()
	if !strings.Contains(info, "Not initialized") {
		t.Fatalf("expected not initialized message, got: %s", info)
	}

	om, err := InitOMEMO(123456)
	if err != nil {
		t.Fatalf("InitOMEMO error: %v", err)
	}
	c.omemo = om
	info2 := c.OMEMOInfo()
	if !strings.Contains(info2, "123456") {
		t.Fatalf("expected device ID in info, got: %s", info2)
	}
	if !strings.Contains(info2, om.Fingerprint) {
		t.Fatalf("expected fingerprint in info, got: %s", info2)
	}
}

func TestHasKeyForDevice(t *testing.T) {
	var enc rawEncrypted
	enc.Header.SID = 42
	enc.Header.Keys = []rawKey{{RID: 100}, {RID: 200}}
	var kj rawKeysJID
	kj.Keys = []rawKey{{RID: 300}}
	enc.Header.KeysJID = []rawKeysJID{kj}

	if !hasKeyForDevice(&enc, 100) {
		t.Fatal("expected key for device 100")
	}
	if !hasKeyForDevice(&enc, 300) {
		t.Fatal("expected key for device 300")
	}
	if hasKeyForDevice(&enc, 999) {
		t.Fatal("did not expect key for device 999")
	}
	if hasKeyForDevice(nil, 100) {
		t.Fatal("expected false for nil encrypted element")
	}
}

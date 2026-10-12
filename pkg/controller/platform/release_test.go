package platform

import (
	"testing"
	"time"
)

func TestNewestOfficial(t *testing.T) {
	got, ok := newestOfficial([]ghRelease{
		{TagName: "v3.4.0-rc.2", Prerelease: true},
		{TagName: "v3.9.0"},
		{TagName: "v3.10.0"},
		{TagName: "v3.2.0-r2"},
		{TagName: "v3.2.0"},
		{TagName: "nightly", Draft: true},
	})
	if !ok || got != "v3.10.0" {
		t.Fatalf("newest %q ok=%v", got, ok)
	}
}

func TestResolveSticky(t *testing.T) {
	now := time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)
	tag, write, after, err := resolveCoreVersion("", "", "", now, func() (string, error) {
		t.Fatal("lookup")
		return "", nil
	})
	if err != nil || tag != "v3.2.0" || !write || after != 0 {
		t.Fatalf("empty: %s write=%v after=%s err=%v", tag, write, after, err)
	}
	tag, write, after, err = resolveCoreVersion("v3.4.0-rc.2", "v3.4.0-rc.2", "", now, nil)
	if err != nil || tag != "v3.4.0-rc.2" || write || after != 0 {
		t.Fatalf("rc: %s write=%v after=%s err=%v", tag, write, after, err)
	}
	if _, _, _, err = resolveCoreVersion("latest", "", "", now, nil); err == nil {
		t.Fatal("latest should be rejected")
	}
}

func TestResolveAuto(t *testing.T) {
	now := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	checked := now.Add(-time.Hour).Format(time.RFC3339)
	tag, write, after, err := resolveCoreVersion(autoVersion, "v3.2.0", checked, now, func() (string, error) {
		t.Fatal("lookup")
		return "", nil
	})
	if err != nil || tag != "v3.2.0" || write || after <= 0 || after > 23*time.Hour {
		t.Fatalf("fresh: %s write=%v after=%s err=%v", tag, write, after, err)
	}
	calls := 0
	tag, write, _, err = resolveCoreVersion(autoVersion, "v3.2.0", now.Add(-25*time.Hour).Format(time.RFC3339), now, func() (string, error) {
		calls++
		return "v3.3.0", nil
	})
	if err != nil || calls != 1 || tag != "v3.3.0" || !write {
		t.Fatalf("due: %s write=%v calls=%d err=%v", tag, write, calls, err)
	}
	tag, write, after, err = resolveCoreVersion(autoVersion, "v3.2.0", "", now, func() (string, error) {
		return "", errLookup
	})
	if err != nil || tag != "v3.2.0" || write || after != releaseRetry {
		t.Fatalf("keep: %s write=%v after=%s err=%v", tag, write, after, err)
	}
}

var errLookup = errString("lookup failed")

type errString string

func (e errString) Error() string { return string(e) }

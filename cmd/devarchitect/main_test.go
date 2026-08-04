package main

import "testing"

func TestRun_Version(t *testing.T) {
	if code := run([]string{"version"}); code != 0 {
		t.Errorf("run(version) = %d, want 0", code)
	}
}

func TestRun_Help(t *testing.T) {
	if code := run([]string{"help"}); code != 0 {
		t.Errorf("run(help) = %d, want 0", code)
	}
}

func TestRun_NoArgs(t *testing.T) {
	if code := run(nil); code != 1 {
		t.Errorf("run(nil) = %d, want 1", code)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	if code := run([]string{"bogus"}); code != 1 {
		t.Errorf("run(bogus) = %d, want 1", code)
	}
}

func TestRun_AnalyzeValidPath(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/sample-go-repo"}); code != 0 {
		t.Errorf("run(analyze, sample-go-repo) = %d, want 0", code)
	}
}

func TestRun_AnalyzeMissingPathArg(t *testing.T) {
	if code := run([]string{"analyze"}); code != 1 {
		t.Errorf("run(analyze) = %d, want 1", code)
	}
}

func TestRun_AnalyzeNonexistentPath(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/does-not-exist"}); code != 1 {
		t.Errorf("run(analyze, does-not-exist) = %d, want 1", code)
	}
}

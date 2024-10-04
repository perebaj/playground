package tmp_test

import (
	"testing"

	tmp "github.com/perebaj/playground/golang/coverage"
)

func TestTmpPublic(t *testing.T) {
	if got := tmp.TmpPublic("a"); got != 1 {
		t.Errorf("TmpPublic() = %v, want 1", got)
	}

	if got := tmp.TmpPublic("b"); got != 2 {
		t.Errorf("TmpPublic() = %v, want 2", got)
	}

	if got := tmp.TmpPublic("c"); got != 0 {
		t.Errorf("TmpPublic() = %v, want 0", got)
	}
}

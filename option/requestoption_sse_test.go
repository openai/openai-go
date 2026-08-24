package option

import "testing"

func TestWithSSEMaxEventBytesRejectsNegativeLimit(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("WithSSEMaxEventBytes did not panic for a negative limit")
		}
	}()
	_ = WithSSEMaxEventBytes(-1)
}

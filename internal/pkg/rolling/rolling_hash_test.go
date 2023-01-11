package rolling

import "testing"

func TestRollingHash(t *testing.T) {
	tests := []struct {
		data   []byte
		result uint64
	}{
		{[]byte{1}, 1},
		{[]byte{1, 2}, 5},
		{[]byte{255}, 255},
		{[]byte{255, 254}, 763},
	}

	for _, test := range tests {
		result := Hash(test.data)
		if result != test.result {
			t.Errorf("unexpected result for data %v: got %d, want %d", test.data, result, test.result)
		}
	}
}

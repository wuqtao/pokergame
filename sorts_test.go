package pokergame

import (
	"testing"
)

func TestBubbleSortIntMin2Max(t *testing.T) {
	ints := []int{1,3,5,4,6,2}
	BubbleSortIntMin2Max(ints)
	if ints[0] != 1 ||
		ints[1] != 2 ||
		ints[2] != 3 ||
		ints[3] != 4 ||
		ints[4] != 5 ||
		ints[5] != 6{
			t.Error("TestBubbleSortIntMin2Max err")
	}

}

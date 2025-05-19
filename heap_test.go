package heap

import (
	"cmp"
	"slices"
	"strconv"
	"testing"
)

var unsortedIntTests = [][]int{
	{1, 5, 3, 4, 2},
	{},
	{0},
	{-10, 10, 25, -0},
}

var unsortedStringTests = [][]string{
	{"banana", "apple", "cherry", "date", "fig"},
	{},
	{"orange"},
	{"grape", "kiwi", "pear", "melon", "pear"},
}

func testNew[T cmp.Ordered](testCase []T) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		testCopy := make([]T, len(testCase))
		copy(testCopy, testCase)
		testCase = testCopy

		heap := New(testCase, func(a, b T) bool { return a < b })
		out := make([]T, len(testCase))
		for i := range out {
			out[i] = heap.Pop()
		}
		if !slices.IsSorted(out) {
			t.Errorf("expected sorted output, got %v", out)
		}
		if heap.Len() != 0 {
			t.Errorf("expected heap to be empty, got %d", heap.Len())
		}
	}
}

func testPush[T cmp.Ordered](testCase []T) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		heap := New([]T{}, func(a, b T) bool { return a < b })
		for _, v := range testCase {
			heap.Push(v)
		}
		out := make([]T, len(testCase))
		for i := range out {
			out[i] = heap.Pop()
		}
		if !slices.IsSorted(out) {
			t.Errorf("expected sorted output, got %v", out)
		}
		if heap.Len() != 0 {
			t.Errorf("expected heap to be empty, got %d", heap.Len())
		}
	}
}

func TestHeapNilIntNew(t *testing.T) {
	t.Parallel()

	heap := New(nil, func(a, b int) bool { return a < b })
	if heap.Len() != 0 {
		t.Errorf("expected new nil int heap to be empty, got %d", heap.Len())
	}

	heap.Push(5)
	heap.Push(3)
	heap.Push(7)
	if heap.Len() != 3 {
		t.Errorf("expected heap to have length 3, got %d", heap.Len())
	}
	if heap.Pop() != 3 {
		t.Errorf("expected heap to pop 3, got %v", heap.Pop())
	}
}

func TestHeapInts(t *testing.T) {
	t.Parallel()

	for i, testCase := range unsortedIntTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			t.Run("New", testNew(testCase))
			t.Run("Push", testPush(testCase))
		})
	}
}

func TestHeapStrings(t *testing.T) {
	t.Parallel()

	for i, testCase := range unsortedStringTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			t.Run("New", testNew(testCase))
			t.Run("Push", testPush(testCase))
		})
	}
}

func FuzzHeapInts(f *testing.F) {
	seeds := [][]int{
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
		{4, 3, 5, 1, 2},
		{0, -0, -459, 1234567890, -1234567890},
	}

	for _, seed := range seeds {
		f.Add(seed[0], seed[1], seed[2], seed[3], seed[4])
	}

	f.Fuzz(func(t *testing.T, one, two, three, four, five int) {
		t.Parallel()

		testCase := []int{one, two, three, four, five}
		t.Run("New", testNew(testCase))
		t.Run("Push", testPush(testCase))
	})
}

func FuzzHeapStrings(f *testing.F) {
	seeds := [][]string{
		{"apple", "banana", "cherry", "date", "elderberry"},
		{"elderberry", "date", "cherry", "banana", "apple"},
		{"cherry", "banana", "elderberry", "apple", "date"},
		{"apple", "apple", "banana", "cherry", "date"},
	}

	for _, seed := range seeds {
		f.Add(seed[0], seed[1], seed[2], seed[3], seed[4])
	}
	f.Fuzz(func(t *testing.T, one, two, three, four, five string) {
		t.Parallel()

		testCase := []string{one, two, three, four, five}
		t.Run("New", testNew(testCase))
		t.Run("Push", testPush(testCase))
	})
}

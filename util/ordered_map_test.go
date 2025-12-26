package util

import (
	"fmt"
	"testing"
)

func TestOrderedMapContains(t *testing.T) {
	orderedMap := NewOrderedMap[string, int]()
	orderedMap.Insert("a", 1)
	orderedMap.Insert("d", 2)
	orderedMap.Insert("c", 3)

	t.Logf("Resulting map: %+v", orderedMap)

	if !orderedMap.Contains("a") {
		t.Fatalf("%+v does not contain key %+v", orderedMap, "a")
	}

	if !orderedMap.Contains("d") {
		t.Fatalf("%+v does not contain key %+v", orderedMap, "d")
	}

	if !orderedMap.Contains("c") {
		t.Fatalf("%+v does not contain key %+v", orderedMap, "c")
	}
}

func TestOrderedMapValue(t *testing.T) {
	orderedMap := NewOrderedMap[string, int]()
	orderedMap.Insert("a", 1)
	orderedMap.Insert("d", 2)
	orderedMap.Insert("c", 3)

	t.Logf("Resulting map: %+v", orderedMap)

	va, err := orderedMap.Value("a")
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if va != 1 {
		t.Fatalf("%+v != %+v for key %+v", 1, va, "a")
	}

	vd, err := orderedMap.Value("d")
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if vd != 2 {
		t.Fatalf("%+v != %+v for key %+v", 2, vd, "d")
	}

	_, err = orderedMap.Value("m")
	if err == nil {
		t.Fatalf("Expected error for key %+v, got none", "m")
	}

	t.Log(err)
}

func TestOrderedIterationOrder(t *testing.T) {
	orderedMap := NewOrderedMap[string, int]()
	orderedMap.Insert("a", 1)
	orderedMap.Insert("d", 2)
	orderedMap.Insert("c", 3)

	t.Logf("Resulting map: %+v", orderedMap)

	// TODO should do assertions for this stuff.
	results := make([]string, 0)
	for k, v := range orderedMap.All() {
		f := fmt.Sprintf("%+v, %+v", k, v)
		results = append(results, f)
	}

	if len(results) != 3 {
		t.Fatalf("Expected to collect 3 entries")
	}

	for _, r := range results {
		t.Log(r)
	}
}

func TestOrderedIteratorEmpty(t *testing.T) {
	orderedMap := NewOrderedMap[string, int]()

	for k, v := range orderedMap.All() {
		t.Logf("%+v %+v", k, v)
		t.Fatalf("Never called")
	}
}

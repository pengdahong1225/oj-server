package utils

import "testing"

func TestSplitStringWithX(t *testing.T) {
	result := SplitStringWithX("#数组#双指针#哈希表", "#")
	t.Log(result)
}

func TestSpliceStringWithX(t *testing.T) {
	src := []string{"数组", "双指针", "哈希表"}
	result := SpliceStringWithX(src, "#")
	t.Log(result)
}

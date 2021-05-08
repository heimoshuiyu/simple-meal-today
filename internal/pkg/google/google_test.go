package google

import "testing"

func TestURL(t *testing.T) {
	google := NewGoogle()
	t.Log(google.GetSearchURL("初音未来"))
}

package gbutils_test

import (
	"testing"

	"io/ioutil"

	"github.com/gophersgang/gbdep/pkg/gbutils"
)

func TestComputeMD5Content(t *testing.T) {
	a := "some content"
	res := gbutils.ComputeMD5Content(a)
	if res != "9893532233caff98cd083a116b013c0b" {
		t.Fatalf("Expected %v, but was %v", "...", res)
	}
}

func TestComputeMD5(t *testing.T) {
	a := "some content"
	f, err := ioutil.TempFile("/tmp", "gbutils")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(a)
	res, err := gbutils.ComputeMD5(f.Name())
	if res != "9893532233caff98cd083a116b013c0b" {
		t.Fatalf("Expected %v, but was %v", "...", res)
	}
}

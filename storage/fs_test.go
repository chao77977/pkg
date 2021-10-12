package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFSType(t *testing.T) {
	path, err := ioutil.TempDir(os.TempDir(), "_pkg")
	defer os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}

	fsInfo, err := GetFSInfo(path)
	if err != nil {
		t.Fatal(err)
	}

	if fsInfo.Type() == "UNKNOWN" {
		t.Error("Unexpected FSType", fsInfo.Type())
	}
}

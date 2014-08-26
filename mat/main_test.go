package mat

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

const (
	fixturePath = "fixtures"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func createTempFile() string {
	file, _ := ioutil.TempFile("", "fixture")
	file.Close()

	return file.Name()
}

func TestOpen(t *testing.T) {
	path := path.Join(fixturePath, "data.mat")

	file, err := Open(path, "r")

	assertEqual(nil, err, t)
	assertEqual(path, file.path, t)

	file.Close()
}

func TestPutMatrix(t *testing.T) {
	path := createTempFile()
	defer os.Remove(path)

	name, rows, cols := "a", uint32(2), uint32(3)
	data := []float64{1, 4, 2, 5, 3, 6}

	file, _ := Open(path, "w7.3")
	err := file.PutMatrix(name, rows, cols, data)

	assertEqual(nil, err, t)

	file.Close()
}

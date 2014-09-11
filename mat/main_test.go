package mat

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-math/support/assert"
)

const (
	fixturePath = "fixtures"
)

func createTempFile() string {
	file, _ := ioutil.TempFile("", "fixture")
	file.Close()

	return file.Name()
}

func TestOpen(t *testing.T) {
	path := path.Join(fixturePath, "data.mat")

	file, err := Open(path, "r")

	assert.Success(err, t)

	file.Close()
}

func TestPut(t *testing.T) {
	path := createTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	name, value := "a", struct{
		A int8
		B uint8
		C int16
		D uint16
		E int32
		F uint32
		G int64
		H uint64
		I float32
		J float64
		Z []float64
	}{
		-1, 2, -3, 4, -5, 6, -7, 8, -9, 10,
		[]float64{1, 0, 1, 0, 1, 0},
	}

	assert.Success(file.Put(name, value), t)
}

func TestPutMatrix(t *testing.T) {
	path := createTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	name, rows, cols := "a", uint32(2), uint32(3)
	data := []float64{1, 2, 3, 4, 5, 6}

	assert.Success(file.PutMatrix(name, data, rows, cols), t)
}

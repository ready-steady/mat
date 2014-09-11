package mat

import (
	"fmt"
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

	objects := []interface{}{
		int8(-1),
		[]int8{-1, -1, -1},

		uint8(2),
		[]uint8{2, 2, 2},

		int16(-3),
		[]int16{-3, -3, -3},

		uint16(4),
		[]uint16{4, 4, 4},

		int32(-5),
		[]int32{-5, -5, -5},

		uint32(6),
		[]uint32{6, 6, 6},

		int64(-7),
		[]int64{-7, -7, -7},

		uint64(8),
		[]uint64{8, 8, 8},

		float32(9),
		[]float32{9, 9, 9},

		float64(10),
		[]float64{10, 10, 10},

		struct {
			A []float64
			B []float64
		}{
			A: []float64{1, 2, 3},
			B: []float64{4, 5, 6},
		},
	}

	for i, o := range objects {
		assert.Success(file.Put(fmt.Sprintf("%c", 'A'+i), o), t)
	}
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

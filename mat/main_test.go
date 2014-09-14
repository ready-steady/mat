package mat

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-math/support/assert"
	"github.com/go-math/support/fixture"
)

func TestOpen(t *testing.T) {
	path := findFixture("data.mat")

	file, err := Open(path, "r")

	assert.Success(err, t)

	file.Close()
}

func TestPut(t *testing.T) {
	path := fixture.MakeTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	for i, o := range fixtureObjects {
		err := file.Put(fmt.Sprintf("%c", 'A'+i), o)

		assert.Success(err, t)
	}
}

func TestPutMatrix(t *testing.T) {
	path := fixture.MakeTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	name, rows, cols := "a", uint32(2), uint32(3)
	data := []float64{1, 2, 3, 4, 5, 6}

	assert.Success(file.PutMatrix(name, data, rows, cols), t)
}

func TestGet(t *testing.T) {
	path := findFixture("data.mat")

	file, _ := Open(path, "r")
	defer file.Close()

	for i, o := range fixtureObjects {
		ptr := makeEmptyLike(o)
		err := file.Get(fmt.Sprintf("%c", 'A'+i), ptr)

		assert.Success(err, t)
		assert.Equal(reflect.Indirect(reflect.ValueOf(ptr)).Interface(), o, t)
	}
}

package mat

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ready-steady/support/assert"
	"github.com/ready-steady/support/fixture"
)

func TestPut(t *testing.T) {
	path := fixture.MakeTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	for i, o := range fixtureObjects {
		assert.Success(file.Put(fmt.Sprintf("%c", 'A'+i), o), t)
	}
}

func TestPutMatrix(t *testing.T) {
	path := fixture.MakeTempFile()
	defer os.Remove(path)

	file, _ := Open(path, "w7.3")
	defer file.Close()

	name, rows, cols := "a", uint(2), uint(3)
	data := []float64{1, 2, 3, 4, 5, 6}

	assert.Success(file.PutMatrix(name, data, rows, cols), t)
}

func TestGet(t *testing.T) {
	path := findFixture("data.mat")

	file, _ := Open(path, "r")
	defer file.Close()

	for i, o := range fixtureObjects {
		v := reflect.New(reflect.TypeOf(o))
		p := v.Interface()
		assert.Success(file.Get(fmt.Sprintf("%c", 'A'+i), p), t)
		assert.Equal(reflect.Indirect(v).Interface(), o, t)
	}
}

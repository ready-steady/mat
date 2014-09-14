package mat

import (
	"path"
)

const (
	fixturePath = "fixtures"
)

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}

type fixtureStruct struct {
	A []float64
	B []float64
}

var fixtureObjects = []interface{}{
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

	fixtureStruct{[]float64{1, 2, 3}, []float64{4, 5, 6}},
}

func makeEmptyLike(object interface{}) interface{} {
	switch object.(type) {
	case int8:
		var v int8
		return &v
	case []int8:
		v := []int8{}
		return &v
	case uint8:
		var v uint8
		return &v
	case []uint8:
		v := []uint8{}
		return &v
	case int16:
		var v int16
		return &v
	case []int16:
		v := []int16{}
		return &v
	case uint16:
		var v uint16
		return &v
	case []uint16:
		v := []uint16{}
		return &v
	case int32:
		var v int32
		return &v
	case []int32:
		v := []int32{}
		return &v
	case uint32:
		var v uint32
		return &v
	case []uint32:
		v := []uint32{}
		return &v
	case int64:
		var v int64
		return &v
	case []int64:
		v := []int64{}
		return &v
	case uint64:
		var v uint64
		return &v
	case []uint64:
		v := []uint64{}
		return &v
	case float32:
		var v float32
		return &v
	case []float32:
		v := []float32{}
		return &v
	case float64:
		var v float64
		return &v
	case []float64:
		v := []float64{}
		return &v
	case fixtureStruct:
		var v fixtureStruct
		return &v
	default:
		return nil
	}
}

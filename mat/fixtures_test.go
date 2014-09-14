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

	struct {
		A []float64
		B []float64
	}{
		[]float64{1, 2, 3},
		[]float64{4, 5, 6},
	},
}

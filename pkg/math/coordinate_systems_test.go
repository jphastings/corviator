package math_test

import (
	"testing"

	"github.com/jphastings/corviator/pkg/math"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestCoordinate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coordinate Suite")
}

var _ = Describe("LLACoords.ECEF()", func() {
	accuracy := 1 // Meters

	DescribeTable("known coordinate transforms",
		func(input math.LLACoords, expected math.ECEFCoords) {
			actual := input.ECEF(math.WGS84)
			Expect(actual.X).To(BeNumerically("~", expected.X, accuracy))
			Expect(actual.Y).To(BeNumerically("~", expected.Y, accuracy))
			Expect(actual.Z).To(BeNumerically("~", expected.Z, accuracy))
		},
		Entry("surface X",
			math.LLACoords{Φ: 0, Λ: 0, A: 0},
			math.ECEFCoords{X: 6378137, Y: 0, Z: 0}),
		Entry("surface Y",
			math.LLACoords{Φ: 0, Λ: 90, A: 0},
			math.ECEFCoords{X: 0, Y: 6378137, Z: 0}),
		Entry("surface Z",
			math.LLACoords{Φ: 90, Λ: 0, A: 0},
			math.ECEFCoords{X: 0, Y: 0, Z: 6356752}),
		Entry("surface -X",
			math.LLACoords{Φ: 0, Λ: 180, A: 0},
			math.ECEFCoords{X: -6378137, Y: 0, Z: 0}),
		Entry("surface -Y",
			math.LLACoords{Φ: 0, Λ: -90, A: 0},
			math.ECEFCoords{X: 0, Y: -6378137, Z: 0}),
		Entry("surface -Z",
			math.LLACoords{Φ: -90, Λ: 0, A: 0},
			math.ECEFCoords{X: 0, Y: 0, Z: -6356752}),
		Entry("Greenwich observatory",
			math.LLACoords{Φ: 51.4769, Λ: 0.0005, A: 48},
			math.ECEFCoords{X: 3980689, Y: 35, Z: 4966800}),
		Entry("Vernadsky Station Bar",
			math.LLACoords{Φ: -65.245724, Λ: -64.257668, A: 4},
			math.ECEFCoords{X: 1163168, Y: -2412321, Z: -5769239}),
		Entry("Marist Brothers Primary School in Suva",
			math.LLACoords{Φ: -18.140535, Λ: 178.428644, A: 33},
			math.ECEFCoords{X: -6060835, Y: 166262, Z: -1973182}),
	)
})
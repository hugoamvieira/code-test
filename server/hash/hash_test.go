package hash

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashNew(t *testing.T) {
	Convey("Given a valid string", t, func() {
		s := "hello"
		Convey("it should return its hash representation", func() {
			hash := New(s)
			So(hash, ShouldEqual, "62c0215")

			Convey("and that should be consistent", func() {
				hashSecond := New(s)
				So(hashSecond, ShouldEqual, hash)
			})
		})
	})
}

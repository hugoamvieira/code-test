package data

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("Calling data.New()", t, func() {
		websiteURL := "https://website.com"
		sessionID := "session"

		d := New(websiteURL, sessionID)

		Convey("should return a 'bare' data reference object with the same data", func() {
			So(d, ShouldNotBeNil)
			So(d.WebsiteURL, ShouldEqual, websiteURL)
			So(d.SessionID, ShouldEqual, sessionID)
			So(d.CopyAndPaste, ShouldBeEmpty)
			So(d.FormCompletionTime, ShouldBeZeroValue)
			So(d.ResizeFrom.Height, ShouldEqual, "")
			So(d.ResizeFrom.Width, ShouldEqual, "")
			So(d.ResizeTo.Height, ShouldEqual, "")
			So(d.ResizeTo.Width, ShouldEqual, "")
		})

		Convey("should store the object in the datastore", func() {
			stored, ok, err := Ds.Get(websiteURL, sessionID)
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
			So(stored, ShouldEqual, d)
		})
	})
}

package api

import (
	"testing"

	"github.com/hugoamvieira/code-test/server/data"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidForCopyPasteEvent(t *testing.T) {
	Convey("For an existing session", t, func() {
		websiteURL := "https://www.website1.com"
		session := "validSession1"

		d := data.New(websiteURL, session)

		Convey("given a valid copy-paste event", func() {
			cpe := &copyPasteEvent{
				WebsiteURL: d.WebsiteURL,
				SessionID:  d.SessionID,
				InputID:    "cardNumber",
			}

			Convey("it should be valid", func() {
				valid, err := cpe.Valid()
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})
	})

	Convey("Given an event for a session that doesn't exist", t, func() {
		cpe := &copyPasteEvent{
			WebsiteURL: "https://www.website2.com",
			SessionID:  "noSession2",
			InputID:    "cardNumber",
		}

		Convey("it should not be valid", func() {
			valid, err := cpe.Valid()
			So(err, ShouldBeNil)
			So(valid, ShouldBeFalse)
		})
	})
}

func TestValidForResizePageEvent(t *testing.T) {
	Convey("For an existing session", t, func() {
		websiteURL := "https://www.website3.com"
		session := "validSession3"

		d := data.New(websiteURL, session)

		Convey("given a valid resize page event", func() {
			rpe := &resizePageEvent{
				WebsiteURL: d.WebsiteURL,
				SessionID:  d.SessionID,
				ResizeFrom: data.Dimension{
					Width:  "100",
					Height: "200",
				},
				ResizeTo: data.Dimension{
					Width:  "101",
					Height: "201",
				},
			}

			Convey("it should be valid", func() {
				valid, err := rpe.Valid()
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})

		Convey("given an invalid resize page event", func() {
			rpe := &resizePageEvent{
				WebsiteURL: d.WebsiteURL,
				SessionID:  d.SessionID,
				ResizeFrom: data.Dimension{
					Width:  "",
					Height: "",
				},
				ResizeTo: data.Dimension{
					Width:  "",
					Height: "",
				},
			}

			Convey("it should not be valid", func() {
				valid, err := rpe.Valid()
				So(err, ShouldBeNil)
				So(valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an event for a session that doesn't exist", t, func() {
		rpe := &resizePageEvent{
			WebsiteURL: "https://www.website4.com",
			SessionID:  "noSession4",
			ResizeFrom: data.Dimension{
				Width:  "100",
				Height: "200",
			},
			ResizeTo: data.Dimension{
				Width:  "101",
				Height: "201",
			},
		}

		Convey("it should not be valid", func() {
			valid, err := rpe.Valid()
			So(err, ShouldBeNil)
			So(valid, ShouldBeFalse)
		})
	})
}

func TestValidForTimeTakenEvent(t *testing.T) {
	Convey("For an existing session", t, func() {
		websiteURL := "https://www.website5.com"
		session := "validSession5"

		d := data.New(websiteURL, session)

		Convey("given a valid time taken event", func() {
			tte := &timeTakenEvent{
				WebsiteURL: d.WebsiteURL,
				SessionID:  d.SessionID,
				TimeTaken:  10,
			}

			Convey("it should be valid", func() {
				valid, err := tte.Valid()
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})

		Convey("given an invalid time taken event", func() {
			tte := &timeTakenEvent{
				WebsiteURL: d.WebsiteURL,
				SessionID:  d.SessionID,
				TimeTaken:  -1,
			}

			Convey("it should not be valid", func() {
				valid, err := tte.Valid()
				So(err, ShouldBeNil)
				So(valid, ShouldBeFalse)
			})
		})
	})

	Convey("Given an event for a session that doesn't exist", t, func() {
		rpe := &timeTakenEvent{
			WebsiteURL: "https://www.website6.com",
			SessionID:  "noSession6",
			TimeTaken:  10,
		}

		Convey("it should not be valid", func() {
			valid, err := rpe.Valid()
			So(err, ShouldBeNil)
			So(valid, ShouldBeFalse)
		})
	})
}

func TestValidForNewSessionRequestEvent(t *testing.T) {
	Convey("Given a valid new session request event", t, func() {
		nsr := &newSessionRequest{
			WebsiteURL: "https://www.website7.com",
		}

		Convey("it should be valid", func() {
			valid := nsr.Valid()
			So(valid, ShouldBeTrue)
		})
	})
}

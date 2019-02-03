package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAndSet(t *testing.T) {
	Convey("Given a valid session in the map,", t, func() {
		sg := &sessionGen{
			sessions: make(map[string]bool),
		}

		sessionID := "session"
		sg.Set(sessionID)

		Convey("it should successfully get it", func() {
			ok := sg.Get(sessionID)
			So(ok, ShouldBeTrue)
		})
	})

	Convey("Given a sessionGen element with no sessions,", t, func() {
		sg := &sessionGen{
			sessions: make(map[string]bool),
		}

		Convey("it should fail to get something from the map", func() {
			ok := sg.Get("shouldnt-exist")
			So(ok, ShouldBeFalse)
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("Given a valid session,", t, func() {
		sg := &sessionGen{
			sessions: make(map[string]bool),
		}

		sessionID := "session"
		sg.Set(sessionID)

		Convey("it should have been successfully set in the first place", func() {
			ok := sg.Get(sessionID)
			So(ok, ShouldBeTrue)

			Convey("and Delete should remove it", func() {
				sg.Delete(sessionID)
				ok = sg.Get(sessionID)
				So(ok, ShouldBeFalse)
			})
		})
	})
}

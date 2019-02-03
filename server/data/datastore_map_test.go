package data

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDatastoreMapGet(t *testing.T) {
	Convey("Given an existing element in the map", t, func() {
		dm := &DatastoreMap{
			m: make(map[string]*Data),
		}

		d := &Data{
			WebsiteURL: "https://validwebsite.com",
			SessionID:  "validSessionForValidWebsite",
		}

		dm.m[getStoreKey(d.WebsiteURL, d.SessionID)] = d

		Convey("it should successfully obtain it", func() {
			obtained, exists, err := dm.Get(d.WebsiteURL, d.SessionID)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(obtained, ShouldEqual, d)
		})
	})

	Convey("Given an empty map", t, func() {
		dm := &DatastoreMap{
			m: make(map[string]*Data),
		}

		Convey("it shouldn't obtain anything", func() {
			obtained, exists, err := dm.Get("doesntexist", "alsodoesntexist")
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
			So(obtained, ShouldBeNil)
		})
	})

	Convey("Given a map with an element in it", t, func() {
		dm := &DatastoreMap{
			m: make(map[string]*Data),
		}

		d := &Data{
			WebsiteURL: "https://validwebsite.com",
			SessionID:  "validSessionForValidWebsite",
		}

		dm.m[getStoreKey(d.WebsiteURL, d.SessionID)] = d

		Convey("when trying to obtain another element, it should return nothing", func() {
			obtained, exists, err := dm.Get("thisdoesntexist", "nah")
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
			So(obtained, ShouldBeNil)
		})
	})
}

func TestDatastoreMapMutate(t *testing.T) {
	// TODO
}

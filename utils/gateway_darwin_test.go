package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDefaultGateway(t *testing.T) {
	Convey("default gateway should not be empty if network is ready", t, func() {
		So(GetDefaultGateway(), ShouldNotEqual, "")
	})
}

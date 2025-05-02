package enpoint_controller

import (
	"quicky-go/enpoint_controller/check"
)

var EnpointList = make([]interface{}, 0)

func GetEnpointList() []interface{} {
	return EnpointList
}
func init() {
	EnpointList = append(EnpointList, check.Hello)
}

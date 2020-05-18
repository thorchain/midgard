package usecase

import . "gopkg.in/check.v1"

var _ = Suite(&EventHandlerSuite{})

type EventHandlerSuite struct {
	dummyStore *StoreDummy
}

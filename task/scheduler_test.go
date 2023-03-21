package task

import (
	"testing"

	types "github.com/more-than-code/deploybot-service-builder/deploybot-types"
)

func TestHandleEvent(t *testing.T) {
	s := NewScheduler()
	s.PushEvent(types.Event{Key: "build", Value: "geoy-webapp"})

	e := s.PullEvent()
	if e.Key != "build" {
		t.Fail()
	}
}

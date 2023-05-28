package unnecessary

import (
	"fmt"
)

// type BehaviorCallback func() error
type BehaviorAjaxCallback func(target *AjaxTarget) error

type Behavior struct {
	Id            string
	page          *Component
	component     *Component
	ajaxEventCode string
	ajaxCallback  BehaviorAjaxCallback
	// callback      BehaviorCallback
}

func NewAjaxBehavior(event string, callback BehaviorAjaxCallback) *Behavior {
	behavior := Behavior{
		ajaxEventCode: event,
		ajaxCallback:  callback,
	}
	return &behavior
}

func (behavior *Behavior) GetScript(lazy bool) *string {
	if behavior.ajaxCallback != nil {
		// function UnnecessaryAddEventListener(lazy, behaviorId, eventCode, elementId)
		script := fmt.Sprintf("UnnecessaryAddEventListener(%t, \"%s\", \"%s\", \"%s\");",
			lazy, behavior.Id, behavior.ajaxEventCode, behavior.component.GetMarkupId())
		return &script
	} else {
		return nil
	}
}

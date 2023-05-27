package unnecessary

import (
	"fmt"
)

// type BejaviorCallback func() error
type BejaviorAjaxCallback func(target *AjaxTarget) error

type Bejavior struct {
	Id            string
	page          *Component
	component     *Component
	ajaxEventCode string
	ajaxCallback  BejaviorAjaxCallback
	// callback      BejaviorCallback
}

func NewAjaxBejavior(event string, callback BejaviorAjaxCallback) *Bejavior {
	bejavior := Bejavior{
		ajaxEventCode: event,
		ajaxCallback:  callback,
	}
	return &bejavior
}

func (bejavior *Bejavior) GetScript(lazy bool) *string {
	if bejavior.ajaxCallback != nil {
		// function UnnecessaryAddEventListener(lazy, bejaviorId, eventCode, elementId)
		script := fmt.Sprintf("UnnecessaryAddEventListener(%t, \"%s\", \"%s\", \"%s\");",
			lazy, bejavior.Id, bejavior.ajaxEventCode, bejavior.component.GetMarkupId())
		return &script
	} else {
		return nil
	}
}

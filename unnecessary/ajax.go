package unnecessary

import (
	"encoding/json"
	"fmt"
)

type AjaxRequest struct {
	BejaviorId string `json:"bejaviorId"`
	ElementId  string `json:"elementId"`
	Event      string `json:"event"`
}

type AjaxResponse struct {
	Checksum string             `json:"checksum"`
	Items    []AjaxResponseItem `json:"items"`
}

type AjaxResponseItem struct {
	MarkupId string  `json:"markupId"`
	Head     *string `json:"head"`
	Body     *string `json:"body"`
	Script   *string `json:"script"`
}

type AjaxTarget struct {
	page         *Component
	ajaxRequest  *AjaxRequest
	ajaxResponse *AjaxResponse
	body         []*Component
	script       []string
}

func NewAjaxTarget(page *Component) *AjaxTarget {
	target := AjaxTarget{
		page: page,
	}
	return &target
}

func (target *AjaxTarget) Add(component *Component) {
	target.body = append(target.body, component)
}

func (target *AjaxTarget) AddScript(script string) {
	target.script = append(target.script, script)
}

func (target *AjaxTarget) Unmarshal(body []byte) error {
	var ajaxRequest AjaxRequest
	if err := json.Unmarshal(body, &ajaxRequest); err != nil {
		return err
	}
	target.ajaxRequest = &ajaxRequest
	return nil
}

func (target *AjaxTarget) Marshal() ([]byte, error) {
	target.ajaxResponse = &AjaxResponse{}
	for _, component := range target.body {
		body, err := RenderComponent(component)
		if err != nil {
			return nil, err
		}

		script := collectAjaxScript(component)
		script = fmt.Sprintf("(function(){%s;})();", script)

		target.ajaxResponse.Items = append(target.ajaxResponse.Items, AjaxResponseItem{
			MarkupId: component.GetMarkupId(),
			Body:     &body,
			Script:   &script,
		})
	}
	data, err := json.Marshal(target.ajaxResponse)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (target *AjaxTarget) Process(requestData []byte) ([]byte, error) {
	if err := target.Unmarshal(requestData); err != nil {
		return nil, err
	}
	bejavior, ok := target.page.bejaviors[target.ajaxRequest.BejaviorId]
	if !ok {
		return nil, fmt.Errorf("bejavior %s not found", target.ajaxRequest.BejaviorId)
	}
	// holder := FindAjaxHandler(target.page, target.ajaxRequest.Handler)
	bejavior.ajaxCallback(target)
	responseData, err := target.Marshal()
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func ProcessAjaxRequest(page *Component, requestData []byte) ([]byte, error) {
	target := NewAjaxTarget(page)
	return target.Process(requestData)
}

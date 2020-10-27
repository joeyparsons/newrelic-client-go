package alerts

type AlertViolationEntity struct {
	Product string `json:"product"`
	Type    string `json:"type"`
	GroupId int    `json:"group_id"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
}

type AlertViolation struct {
	ID            int                  `json:"id"`
	Label         string               `json:"label"`
	Duration      int                  `json:"duration"`
	PolicyName    string               `json:"policy_name"`
	ConditionName string               `json:"condition_name"`
	Priority      string               `json:"priority"`
	OpenedAt      int                  `json:"opened_at"`
	ClosedAt      int                  `json:"closed_at"`
	Entity        AlertViolationEntity `json:"entity"`
}

type ListAlertViolationsParams struct {
	Page int `url:"page,omitempty"`
}

// ListAlertEvents is used to retrieve New Relic alert events
func (a *Alerts) ListAlertViolations(params *ListAlertViolationsParams) ([]*AlertViolation, error) {
	alertViolations := []*AlertViolation{}
	nextURL := a.config.Region().RestURL("alerts_violations.json")

	for nextURL != "" {
		response := alertViolationsResponse{}
		resp, err := a.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		alertViolations = append(alertViolations, response.AlertViolations...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertViolations, nil
}

type alertViolationsResponse struct {
	AlertViolations []*AlertViolation `json:"violations,omitempty"`
}

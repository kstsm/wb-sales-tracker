package dto

type ItemResponse struct {
	ID        string  `json:"id,omitempty"`
	Type      string  `json:"type,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Date      string  `json:"date,omitempty"`
	Category  string  `json:"category,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type ItemsListResponse struct {
	Items []ItemResponse `json:"items"`
	Total int            `json:"total"`
}

type AnalyticsResponse struct {
	From         string             `json:"from,omitempty"`
	To           string             `json:"to,omitempty"`
	Sum          float64            `json:"sum"`
	Avg          *float64           `json:"avg,omitempty"`
	Count        int                `json:"count"`
	Median       *float64           `json:"median,omitempty"`
	Percentile90 *float64           `json:"percentile_90,omitempty"`
	Grouped      []GroupedAnalytics `json:"grouped,omitempty"`
}

type GroupedAnalytics struct {
	Group        string   `json:"group"`
	Sum          *float64 `json:"sum,omitempty"`
	Avg          *float64 `json:"avg,omitempty"`
	Count        int      `json:"count"`
	Median       *float64 `json:"median,omitempty"`
	Percentile90 *float64 `json:"percentile90,omitempty"`
}

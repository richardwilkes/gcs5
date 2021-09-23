package settings

// Attribute holds the definition of an attribute.
type Attribute struct {
	ID                  string       `json:"id"`
	Type                string       `json:"type"`
	Name                string       `json:"name"`
	FullName            string       `json:"full_name"`
	AttributeBase       string       `json:"attribute_base"`
	CostPerPoint        int          `json:"cost_per_point"`
	CostAdjPercentPerSm int          `json:"cost_adj_percent_per_sm"`
	Thresholds          []*Threshold `json:"thresholds,omitempty"`
}

// Threshold holds a point within an attribute pool where changes in state occur.
type Threshold struct {
	State       string   `json:"state"`
	Explanation string   `json:"explanation"`
	Multiplier  int      `json:"multiplier"`
	Divisor     int      `json:"divisor"`
	Addition    int      `json:"addition"`
	Ops         []string `json:"ops"`
}

// FactoryAttributes returns the attribute factory settings.
func FactoryAttributes() []*Attribute {
	// TODO: Fill
	return nil
}

package settings

// HitLocations holds a set of hit locations.
type HitLocations struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Roll      string      `json:"roll"`
	Locations []*Location `json:"locations"`
}

// Location holds a single hit location.
type Location struct {
	ID          string       `json:"id"`
	ChoiceName  string       `json:"choice_name"`
	TableName   string       `json:"table_name"`
	Slots       int          `json:"slots"`
	HitPenalty  int          `json:"hit_penalty"`
	DRBonus     int          `json:"dr_bonus"`
	Description string       `json:"description"`
	Calc        LocationCalc `json:"calc"`
}

// LocationCalc holds values GCS calculates for a Location, but that we want to be present in any json output so that
// other uses of the data don't have to replicate the code to calculate it.
type LocationCalc struct {
	RollRange string `json:"roll_range"`
}

// FactoryHitLocations returns the hit location factory settings.
func FactoryHitLocations() *HitLocations {
	// TODO: Fill
	return &HitLocations{}
}

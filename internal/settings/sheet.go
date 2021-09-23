package settings

// Sheet settings.
type Sheet struct {
	DefaultLengthUnits         string        `json:"default_length_units"`
	DefaultWeightUnits         string        `json:"default_weight_units"`
	UserDescriptionDisplay     string        `json:"user_description_display"`
	ModifiersDisplay           string        `json:"modifiers_display"`
	NotesDisplay               string        `json:"notes_display"`
	SkillLevelAdjDisplay       string        `json:"skill_level_adj_display"`
	UseMultiplicativeModifiers bool          `json:"use_multiplicative_modifiers"`
	UseModifyingDicePlusAdds   bool          `json:"use_modifying_dice_plus_adds"`
	DamageProgression          string        `json:"damage_progression"`
	UseSimpleMetricConversions bool          `json:"use_simple_metric_conversions"`
	ShowCollegeInSheetSpells   bool          `json:"show_college_in_sheet_spells"`
	ShowDifficulty             bool          `json:"show_difficulty"`
	ShowAdvantageModifierAdj   bool          `json:"show_advantage_modifier_adj"`
	ShowEquipmentModifierAdj   bool          `json:"show_equipment_modifier_adj"`
	ShowSpellAdj               bool          `json:"show_spell_adj"`
	UseTitleInFooter           bool          `json:"use_title_in_footer"`
	Page                       Page          `json:"page"`
	BlockLayout                []string      `json:"block_layout"`
	Attributes                 []*Attribute  `json:"attributes"`
	HitLocations               *HitLocations `json:"hit_locations"`
}

// Page settings.
type Page struct {
	PaperSize    string `json:"paper_size"`
	TopMargin    string `json:"top_margin"`
	LeftMargin   string `json:"left_margin"`
	BottomMargin string `json:"bottom_margin"`
	RightMargin  string `json:"right_margin"`
	Orientation  string `json:"orientation"`
}

// NewSheet returns new sheet settings.
func NewSheet() *Sheet {
	return &Sheet{
		DefaultLengthUnits:         "ft_in",     // TODO: Use type
		DefaultWeightUnits:         "lb",        // TODO: Use type
		UserDescriptionDisplay:     "tooltip",   // TODO: Use type
		ModifiersDisplay:           "inline",    // TODO: Use type
		NotesDisplay:               "inline",    // TODO: Use type
		SkillLevelAdjDisplay:       "tooltip",   // TODO: Use type
		DamageProgression:          "basic_set", // TODO: Use type
		UseSimpleMetricConversions: true,
		ShowSpellAdj:               true,
		Page: Page{
			PaperSize:    "na-letter", // TODO: Use type
			TopMargin:    "0.25 in",   // TODO: Use type
			LeftMargin:   "0.25 in",   // TODO: Use type
			BottomMargin: "0.25 in",   // TODO: Use type
			RightMargin:  "0.25 in",   // TODO: Use type
			Orientation:  "portrait",  // TODO: Use type
		},
		BlockLayout:  FactoryBlockLayout(),
		Attributes:   FactoryAttributes(),
		HitLocations: FactoryHitLocations(),
	}
}

// FactoryBlockLayout returns the block layout factory settings.
func FactoryBlockLayout() []string {
	return []string{
		// TODO: Use constants
		"reactions conditional_modifiers",
		"melee",
		"ranged",
		"advantages skills",
		"spells",
		"equipment",
		"other_equipment",
		"notes",
	}
}

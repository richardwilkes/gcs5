/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package gurps

import (
	"fmt"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible SelfControlRollAdj values.
const (
	NoCRAdj SelfControlRollAdj = iota
	ActionPenalty
	ReactionPenalty
	FrightCheckPenalty
	FrightCheckBonus
	MinorCostOfLivingIncrease
	MajorCostOfLivingIncrease
)

type selfControlRollAdjData struct {
	Key         string
	String      string
	Description func(cr SelfControlRoll) string
	Adjustment  func(cr SelfControlRoll) int
	Features    func(cr SelfControlRoll) []*Feature
}

// SelfControlRollAdj holds an adjustment for a self-control roll.
type SelfControlRollAdj uint8

var selfControlRollAdjValues = []*selfControlRollAdjData{
	{
		Key:         "none",
		String:      i18n.Text("None"),
		Description: func(_ SelfControlRoll) string { return "" },
		Adjustment:  func(_ SelfControlRoll) int { return 0 },
		Features:    func(_ SelfControlRoll) []*Feature { return nil },
	},
	{
		Key:    "action_penalty",
		String: i18n.Text("Includes an Action Penalty for Failure"),
		Description: func(cr SelfControlRoll) string {
			if cr == NoneRequired {
				return ""
			}
			return fmt.Sprintf(i18n.Text("%d Action Penalty"), int(cr)-int(NoneRequired))
		},
		Adjustment: func(cr SelfControlRoll) int { return int(cr) - int(NoneRequired) },
		Features:   func(_ SelfControlRoll) []*Feature { return nil },
	},
	{
		Key:    "reaction_penalty",
		String: i18n.Text("Includes a Reaction Penalty for Failure"),
		Description: func(cr SelfControlRoll) string {
			if cr == NoneRequired {
				return ""
			}
			return fmt.Sprintf(i18n.Text("%d Reaction Penalty"), int(cr)-int(NoneRequired))
		},
		Adjustment: func(cr SelfControlRoll) int { return int(cr) - int(NoneRequired) },
		Features:   func(_ SelfControlRoll) []*Feature { return nil },
	},
	{
		Key:    "fright_check_penalty",
		String: i18n.Text("Includes Fright Check Penalty"),
		Description: func(cr SelfControlRoll) string {
			if cr == NoneRequired {
				return ""
			}
			return fmt.Sprintf(i18n.Text("%d Reaction Penalty"), int(cr)-int(NoneRequired))
		},
		Adjustment: func(cr SelfControlRoll) int { return int(cr) - int(NoneRequired) },
		Features:   func(_ SelfControlRoll) []*Feature { return nil },
	},
}

/*
    FRIGHT_CHECK_PENALTY {
        @Override
        public String toString() {
            return I18n.text("Includes Fright Check Penalty");
        }

        @Override
        public String getDescription(SelfControlRoll cr) {
            if (cr == SelfControlRoll.NONE_REQUIRED) {
                return "";
            }
            return MessageFormat.format(I18n.text("{0} Fright Check Penalty"), Numbers.formatWithForcedSign(getAdjustment(cr)));
        }

        @Override
        public int getAdjustment(SelfControlRoll cr) {
            return cr.ordinal() - SelfControlRoll.NONE_REQUIRED.ordinal();
        }
    },
    FRIGHT_CHECK_BONUS {
        @Override
        public String toString() {
            return I18n.text("Includes Fright Check Bonus");
        }

        @Override
        public String getDescription(SelfControlRoll cr) {
            if (cr == SelfControlRoll.NONE_REQUIRED) {
                return "";
            }
            return MessageFormat.format(I18n.text("{0} Fright Check Bonus"), Numbers.formatWithForcedSign(getAdjustment(cr)));
        }

        @Override
        public int getAdjustment(SelfControlRoll cr) {
            return SelfControlRoll.NONE_REQUIRED.ordinal() - cr.ordinal();
        }
    },
    MINOR_COST_OF_LIVING_INCREASE {
        @Override
        public String toString() {
            return I18n.text("Includes a Minor Cost of Living Increase");
        }

        @Override
        public String getDescription(SelfControlRoll cr) {
            if (cr == SelfControlRoll.NONE_REQUIRED) {
                return "";
            }
            return MessageFormat.format(I18n.text("{0}% Cost of Living Increase"), Numbers.formatWithForcedSign(getAdjustment(cr)));
        }

        @Override
        public int getAdjustment(SelfControlRoll cr) {
            return 5 * (SelfControlRoll.NONE_REQUIRED.ordinal() - cr.ordinal());
        }
    },
    MAJOR_COST_OF_LIVING_INCREASE {
        @Override
        public String toString() {
            return I18n.text("Includes a Major Cost of Living Increase and Merchant Skill Penalty");
        }

        @Override
        public String getDescription(SelfControlRoll cr) {
            if (cr == SelfControlRoll.NONE_REQUIRED) {
                return "";
            }
            return MessageFormat.format(I18n.text("{0}% Cost of Living Increase"), Numbers.formatWithForcedSign(getAdjustment(cr)));
        }

        @Override
        public int getAdjustment(SelfControlRoll cr) {
            return switch (cr) {
                case CR6 -> 80;
                case CR9 -> 40;
                case CR12 -> 20;
                case CR15 -> 10;
                default -> 0;
            };
        }

        @Override
        public List<Bonus> getBonuses(SelfControlRoll cr) {
            List<Bonus>    list     = new ArrayList<>();
            SkillBonus     bonus    = new SkillBonus();
            StringCriteria criteria = bonus.getNameCriteria();
            criteria.setType(StringCompareType.IS);
            criteria.setQualifier("Merchant");
            criteria = bonus.getSpecializationCriteria();
            criteria.setType(StringCompareType.ANY);
            LeveledAmount amount = bonus.getAmount();
            amount.setDecimal(false);
            amount.setPerLevel(false);
            amount.setAmount(cr.ordinal() - SelfControlRoll.NONE_REQUIRED.ordinal());
            list.add(bonus);
            return list;
        }
    };

    public abstract String getDescription(SelfControlRoll cr);

    public abstract int getAdjustment(SelfControlRoll cr);

    public List<Bonus> getBonuses(SelfControlRoll cr) {
        return Collections.emptyList();
    }
}
*/

// EnsureValid returns the first SelfControlRollAdj if this SelfControlRollAdj is not a known value.
func (s SelfControlRollAdj) EnsureValid() SelfControlRollAdj {
	if int(s) < len(selfControlRollAdjValues) {
		return s
	}
	return NoCRAdj
}

// String implements fmt.Stringer.
func (s SelfControlRollAdj) String() string {
	return selfControlRollAdjValues[s.EnsureValid()].String
}

// Description returns a formatted description.
func (s SelfControlRollAdj) Description(cr SelfControlRoll) string {
	if cr == NoneRequired {
		return ""
	}
	return selfControlRollAdjValues[s.EnsureValid()].Description(cr)
}

// Adjustment returns the adjustment to make.
func (s SelfControlRollAdj) Adjustment(cr SelfControlRoll) int {
	return selfControlRollAdjValues[s.EnsureValid()].Adjustment(cr)
}

// Features returns the set of features to apply.
func (s SelfControlRollAdj) Features(cr SelfControlRoll) []*Feature {
	return selfControlRollAdjValues[s.EnsureValid()].Features(cr)
}

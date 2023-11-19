package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
	Level int
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked toon of desired level of the race/class you want, ideally without any talents to mess up base stats.
//  Basic stats are as-shown (str/agi/stm/int/spirit)

// Base Spell Crit is calculated by
//   1. Take as-shown value (troll shaman have 3.5%)
//   2. Calculate the bonus from int (for troll shaman that would be 104/78.1=1.331% crit)
//   3. Subtract as-shown from int bouns (3.5-1.331=2.169)
//   4. 2.169*22.08 (rating per crit percent) = 47.89 crit rating.

// Base mana can be looked up here: https://wowwiki-archive.fandom.com/wiki/Base_mana

// These are also scattered in various dbc/casc files,
// `octbasempbyclass.txt`, `combatratings.txt`, `chancetospellcritbase.txt`, etc.

var RaceOffsets = map[proto.Race]stats.Stats{
	proto.Race_RaceUnknown: stats.Stats{},
	proto.Race_RaceHuman:   stats.Stats{},
	proto.Race_RaceOrc: {
		stats.Agility:   -3,
		stats.Strength:  3,
		stats.Intellect: -3,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   1,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   4,
		stats.Strength:  -4,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   0,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -4,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   2,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   0,
	},
	proto.Race_RaceBloodElf: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceDraenei: {
		stats.Agility:   -3,
		stats.Strength:  1,
		stats.Intellect: 0,
		stats.Spirit:    2,
		stats.Stamina:   0,
	},
}

// TODO: Classic base stats
var ClassBaseStats = map[proto.Class]map[int]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassPaladin: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassHunter: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassRogue: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassPriest: {
		25: {
			stats.Health:    222,
			stats.Mana:      217,
			stats.Agility:   26,
			stats.Strength:  22,
			stats.Intellect: 53,
			stats.Spirit:    55,
			stats.Stamina:   44,
		},
		40: {
			stats.Health:    457,
			stats.Mana:      631,
			stats.Agility:   0,
			stats.Strength:  26,
			stats.Intellect: 78,
			stats.Spirit:    81,
			stats.Stamina:   39,
		},
		50: {
			stats.Health:    792,
			stats.Mana:      886,
			stats.Agility:   35,
			stats.Strength:  29,
			stats.Intellect: 98,
			stats.Spirit:    102,
			stats.Stamina:   45,
		},
		60: {
			stats.Health:    1217,
			stats.Mana:      1096,
			stats.Agility:   40,
			stats.Strength:  32,
			stats.Intellect: 120,
			stats.Spirit:    125,
			stats.Stamina:   52,
		},
	},
	proto.Class_ClassShaman: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassMage: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassWarlock: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
	proto.Class_ClassDruid: {
		25: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		40: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
		60: {
			stats.Health:      0,
			stats.Mana:        0,
			stats.Agility:     0,
			stats.Strength:    0,
			stats.Intellect:   0,
			stats.Spirit:      0,
			stats.Stamina:     0,
			stats.AttackPower: 0,
		},
	},
}

// Retrieves base stats, with race offsets, and crit rating adjustments per level
func getBaseStatsCombo(r proto.Race, c proto.Class, level int) stats.Stats {
	if level == 0 {
		level = 60
	}

	starting := ClassBaseStats[c][level]

	return starting.Add(RaceOffsets[r]).Add(ExtraClassBaseStats[c][level])
}

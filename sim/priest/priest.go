package priest

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{15, 16, 16}

const (
	SpellCode_PriestNone int32 = iota
	SpellCode_PriestFlashHeal
	SpellCode_PriestHeal
	SpellCode_PriestGreaterHeal
)

type Priest struct {
	core.Character
	Talents *proto.PriestTalents

	Latency float64

	// Auras
	InnerFocusAura     *core.Aura
	ShadowWeavingAuras core.AuraArray
	WeakenedSouls      core.AuraArray

	// Base Damage Spells
	DevouringPlague []*core.Spell
	HolyFire        []*core.Spell
	MindBlast       []*core.Spell
	MindFlay        [][]*core.Spell
	ShadowWordPain  []*core.Spell
	Smite           []*core.Spell

	// Base Healing Spells
	FlashHeal       []*core.Spell
	GreaterHeal     []*core.Spell
	PowerWordShield []*core.Spell
	PrayerOfHealing []*core.Spell
	Renew           []*core.Spell

	// Other Base Spells
	InnerFocus *core.Spell

	// Runes
	CircleOfHealing   *core.Spell
	MindFlayModifier  float64       // For Twisted Faith
	MindBlastModifier float64       // For Twisted Faith
	MindSear          []*core.Spell // 1 entry for each tick
	Penance           *core.Spell
	PenanceHeal       *core.Spell
	PrayerOfMending   *core.Spell
	EmpoweredRenew    *core.Spell
	ShadowWordDeath   *core.Spell
	VoidPlague        *core.Spell

	ProcPrayerOfMending core.ApplySpellResults

	DpInitMultiplier float64
}

func (priest *Priest) GetCharacter() *core.Character {
	return &priest.Character
}

func (priest *Priest) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ShadowProtection = true
	raidBuffs.DivineSpirit = true

	raidBuffs.PowerWordFortitude = max(raidBuffs.PowerWordFortitude, core.MakeTristateValue(
		true,
		priest.Talents.ImprovedPowerWordFortitude == 2))
}

func (priest *Priest) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (priest *Priest) Initialize() {
	priest.registerMindBlast()
	priest.registerMindFlay()
	priest.registerShadowWordPainSpell()
	priest.registerDevouringPlagueSpell()
	priest.RegisterSmiteSpell()
	priest.registerHolyFire()

	priest.registerPowerInfusionCD()
}

func (priest *Priest) RegisterHealingSpells() {
	// priest.registerFlashHealSpell()
	// priest.registerGreaterHealSpell()
	// priest.registerPowerWordShieldSpell()
	// priest.registerPrayerOfHealingSpell()
	// priest.registerRenewSpell()
}

func (priest *Priest) AddShadowWeavingStack(sim *core.Simulation, target *core.Unit) {
	if priest.ShadowWeavingAuras == nil {
		return
	}

	if sim.RollWithLabel(0, 1, "ShadowWeaving") < (0.2 * float64(priest.Talents.ShadowWeaving)) {
		priest.ShadowWeavingAuras.Get(target).Activate(sim)
		priest.ShadowWeavingAuras.Get(target).AddStack(sim)
	}
}

func (priest *Priest) Reset(_ *core.Simulation) {
	priest.MindFlayModifier = 1
	priest.MindBlastModifier = 1
}

func New(char *core.Character, talents string) *Priest {
	priest := &Priest{
		Character: *char,
		Talents:   &proto.PriestTalents{},
	}
	core.FillTalentsProto(priest.Talents.ProtoReflect(), talents, TalentTreeSizes)

	priest.EnableManaBar()

	priest.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[priest.Class][int(priest.Level)]*core.SpellCritRatingPerCritChance)

	// Set mana regen to 12.5 + Spirit/4 each 2s tick
	priest.SpiritManaRegenPerSecond = func() float64 {
		return 6.25 + priest.GetStat(stats.Spirit)/8
	}

	return priest
}

func (priest *Priest) HasRune(rune proto.PriestRune) bool {
	return priest.HasRuneById(int32(rune))
}

// Agent is a generic way to access underlying priest on any of the agents.
type PriestAgent interface {
	GetPriest() *Priest
}

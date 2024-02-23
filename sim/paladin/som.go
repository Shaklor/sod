package paladin

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

// Seal of Martyrdom is a spell consisting of:
// - A judgement that deals 70% weapon damage that is not normalised.
// - A guaranteed on-hit proc that deals 40% weapon damage that is normalised.

func (paladin *Paladin) registerSealOfMartyrdomSpellAndAura() {

	if !paladin.HasRune(proto.PaladinRune_RuneChestSealofMartyrdom) {
		return
	}

	onJudgementProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 407803}, // Judgement of Righteousness.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplier: 0.7,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	onSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 407799},
		SpellSchool:      core.SpellSchoolHoly,
		ProcMask:         core.ProcMaskEmpty, // This needs figured out properly
		Flags:            core.SpellFlagMeleeMetrics,
		RequiredLevel:    1,
		DamageMultiplier: 0.4,
		ThreatMultiplier: 1.0,
		CritMultiplier:   paladin.MeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.BonusWeaponDamage()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	// TODO: SoM also procs on a variety of weapon proc effects like Talwar's shadow bolt.
	// There is a brief icd preventing SoM and weapon procs chaining back and forth, they
	// can proc each other exactly once.
	// We need the proc rates figured out to see if they differ between weapons and what variables
	// the proc rates may be a function of.
	// icd := core.Cooldown{
	// 	Timer:    paladin.NewTimer(),
	// 	Duration: time.Millisecond * 5,
	// }

	auraActionID := core.ActionID{SpellID: 407798}
	paladin.SealOfMartyrdomAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Martyrdom",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				onSwingProc.Cast(sim, result.Target)
			}
			if spell.ProcMask.Matches(core.ProcMaskProc) {
				// if icd.IsReady(sim) {
				// 	icd.Use(sim)
				// }
				onSwingProc.Cast(sim, result.Target)
			}
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, result.Target)
			}
		},
	})
	// Necessary because of the mix of % base mana cost and flat reduction on the libram
	manaCost := paladin.BaseMana * 0.04
	if paladin.Ranged().ID == LibramOfBenediction {
		manaCost -= 10
	}
	aura := paladin.SealOfMartyrdomAura
	paladin.SealOfMartyrdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1.0 - (float64(paladin.Talents.Benediction) * 0.03),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
			paladin.CurrentSealExpiration = sim.CurrentTime + SealDuration
		},
	})
}

package core

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type ResourceKey struct {
	ActionID ActionID
	Type     proto.ResourceType
}

type DistributionMetrics struct {
	// Values for the current iteration. These are cleared after each iteration.
	Total float64

	// Aggregate values. These are updated after each iteration.
	aggregator
	max     float64
	min     float64
	maxSeed int64
	minSeed int64
	hist    map[int32]int32 // rounded DPS to count
	sample  []float64
}

func (distMetrics *DistributionMetrics) reset() {
	distMetrics.Total = 0
}

// This should be called when a Sim iteration is complete.
func (distMetrics *DistributionMetrics) doneIteration(sim *Simulation) {
	dps := distMetrics.Total / sim.Duration.Seconds()
	distMetrics.add(dps)

	if sim.Options.SaveAllValues {
		if cap(distMetrics.sample) < int(sim.Options.Iterations) {
			distMetrics.sample = make([]float64, 0, sim.Options.Iterations)
		}
		distMetrics.sample = append(distMetrics.sample, dps)
	}

	if dps > distMetrics.max {
		distMetrics.max = dps
		distMetrics.maxSeed = sim.rand.GetSeed()
	}
	if dps <= distMetrics.min || distMetrics.min < 0 {
		distMetrics.min = dps
		distMetrics.minSeed = sim.rand.GetSeed()
	}

	dpsRounded := int32(math.Round(dps/10) * 10)
	distMetrics.hist[dpsRounded]++
}

func (distMetrics *DistributionMetrics) ToProto() *proto.DistributionMetrics {
	mean, stdev := distMetrics.meanAndStdDev()

	return &proto.DistributionMetrics{
		Avg:       mean,
		Stdev:     stdev,
		Max:       distMetrics.max,
		Min:       distMetrics.min,
		MaxSeed:   distMetrics.maxSeed,
		MinSeed:   distMetrics.minSeed,
		Hist:      distMetrics.hist,
		AllValues: distMetrics.sample,
	}
}

func NewDistributionMetrics() DistributionMetrics {
	return DistributionMetrics{
		hist: make(map[int32]int32),
		min:  -1,
	}
}

type UnitMetrics struct {
	dps    DistributionMetrics
	dpasp  DistributionMetrics
	threat DistributionMetrics
	dtps   DistributionMetrics
	tmi    DistributionMetrics
	hps    DistributionMetrics
	tto    DistributionMetrics

	tmiList   []tmiListItem
	isTanking bool
	tmiBin    int32

	CharacterIterationMetrics

	// Aggregate values. These are updated after each iteration.
	numItersDead int32
	oomTimeSum   float64
	actions      map[ActionID]*ActionMetrics
	resources    []*ResourceMetrics
}

// Metrics for the current iteration, for 1 agent. Keep this as a separate
// struct, so it's easy to clear.
type CharacterIterationMetrics struct {
	Died    bool // Whether this unit died in the current iteration.
	WentOOM bool // Whether the agent has hit OOM at least once in this iteration.

	ManaSpent  float64
	ManaGained float64

	OOMTime time.Duration // time spent not casting and waiting for regen.

	FirstOOMTimestamp time.Duration // Timestamp at which unit first went OOM.
}

type ActionMetrics struct {
	IsMelee bool // True if melee action, false if spell action.

	// Metrics for this action, for each possible target.
	Targets []TargetedActionMetrics
}

type tmiListItem struct {
	Timestamp      time.Duration
	WeightedDamage float64
}

func (actionMetrics *ActionMetrics) ToProto(actionID ActionID) *proto.ActionMetrics {
	targetMetrics := make([]*proto.TargetedActionMetrics, 0, len(actionMetrics.Targets))
	for _, tam := range actionMetrics.Targets {
		targetMetrics = append(targetMetrics, tam.ToProto())
	}

	return &proto.ActionMetrics{
		Id:      actionID.ToProto(),
		IsMelee: actionMetrics.IsMelee,
		Targets: targetMetrics,
	}
}

// Metric totals for a spell against a specific target, for the current iteration.
type SpellMetrics struct {
	Casts   int32
	Misses  int32
	Hits    int32
	Crits   int32
	Crushes int32
	Dodges  int32
	Glances int32
	Parries int32
	Blocks  int32

	// Partial or full resists aren't tracked, at the moment, cp. applyResistances()

	TotalDamage    float64 // Damage done by all casts of this spell.
	TotalThreat    float64 // Threat generated by all casts of this spell.
	TotalHealing   float64 // Healing done by all casts of this spell.
	TotalShielding float64 // Shielding done by all casts of this spell.
	TotalCastTime  time.Duration
}

type TargetedActionMetrics struct {
	UnitIndex int32

	Casts   int32
	Hits    int32
	Crits   int32
	Misses  int32
	Dodges  int32
	Parries int32
	Blocks  int32
	Glances int32

	Damage    float64
	Threat    float64
	Healing   float64
	Shielding float64
	CastTime  time.Duration
}

func (tam *TargetedActionMetrics) ToProto() *proto.TargetedActionMetrics {
	return &proto.TargetedActionMetrics{
		UnitIndex: tam.UnitIndex,

		Casts:      tam.Casts,
		Hits:       tam.Hits,
		Crits:      tam.Crits,
		Misses:     tam.Misses,
		Dodges:     tam.Dodges,
		Parries:    tam.Parries,
		Blocks:     tam.Blocks,
		Glances:    tam.Glances,
		Damage:     tam.Damage,
		Threat:     tam.Threat,
		Healing:    tam.Healing,
		Shielding:  tam.Shielding,
		CastTimeMs: float64(tam.CastTime.Milliseconds()),
	}
}

func NewUnitMetrics() UnitMetrics {
	return UnitMetrics{
		dps:     NewDistributionMetrics(),
		dpasp:   NewDistributionMetrics(),
		threat:  NewDistributionMetrics(),
		dtps:    NewDistributionMetrics(),
		tmi:     NewDistributionMetrics(),
		hps:     NewDistributionMetrics(),
		tto:     NewDistributionMetrics(),
		actions: make(map[ActionID]*ActionMetrics),
	}
}

type ResourceMetrics struct {
	ActionID ActionID
	Type     proto.ResourceType

	Events     int32
	Gain       float64
	ActualGain float64

	EventsFromPreviousIterations     int32
	ActualGainFromPreviousIterations float64
}

func (resourceMetrics *ResourceMetrics) ToProto() *proto.ResourceMetrics {
	return &proto.ResourceMetrics{
		Id:   resourceMetrics.ActionID.ToProto(),
		Type: resourceMetrics.Type,

		Events:     resourceMetrics.Events,
		Gain:       resourceMetrics.Gain,
		ActualGain: resourceMetrics.ActualGain,
	}
}

func (resourceMetrics *ResourceMetrics) reset() {
	resourceMetrics.EventsFromPreviousIterations = resourceMetrics.Events
	resourceMetrics.ActualGainFromPreviousIterations = resourceMetrics.ActualGain
}
func (resourceMetrics *ResourceMetrics) EventsForCurrentIteration() int32 {
	return resourceMetrics.Events - resourceMetrics.EventsFromPreviousIterations
}
func (resourceMetrics *ResourceMetrics) ActualGainForCurrentIteration() float64 {
	return resourceMetrics.ActualGain - resourceMetrics.ActualGainFromPreviousIterations
}

func (resourceMetrics *ResourceMetrics) AddEvent(gain float64, actualGain float64) {
	resourceMetrics.Events++
	resourceMetrics.Gain += gain
	resourceMetrics.ActualGain += actualGain
}

func (unitMetrics *UnitMetrics) NewResourceMetrics(actionID ActionID, resourceType proto.ResourceType) *ResourceMetrics {
	newMetrics := &ResourceMetrics{
		ActionID: actionID,
		Type:     resourceType,
	}
	unitMetrics.resources = append(unitMetrics.resources, newMetrics)
	return newMetrics
}

// Convenience helpers for NewResourceMetrics.
func (unit *Unit) NewHealthMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeHealth)
}
func (unit *Unit) NewManaMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeMana)
}
func (unit *Unit) NewRageMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeRage)
}
func (unit *Unit) NewEnergyMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeEnergy)
}
func (unit *Unit) NewComboPointMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeComboPoints)
}
func (unit *Unit) NewFocusMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeFocus)
}

// Adds the results of a spell to the character metrics.
func (unitMetrics *UnitMetrics) addSpellMetrics(spell *Spell, actionID ActionID, spellMetrics []SpellMetrics) {
	actionMetrics, ok := unitMetrics.actions[actionID]

	if !ok {
		actionMetrics = &ActionMetrics{IsMelee: spell.Flags.Matches(SpellFlagMeleeMetrics)}
		unitMetrics.actions[actionID] = actionMetrics
	}

	if len(actionMetrics.Targets) == 0 {
		actionMetrics.Targets = make([]TargetedActionMetrics, len(spellMetrics))
		for i := range actionMetrics.Targets {
			tam := &actionMetrics.Targets[i]
			tam.UnitIndex = spell.Unit.AttackTables[i].Defender.UnitIndex
		}
	}

	for i, spellTargetMetrics := range spellMetrics {
		tam := &actionMetrics.Targets[i]
		tam.Casts += spellTargetMetrics.Casts
		tam.Misses += spellTargetMetrics.Misses
		tam.Hits += spellTargetMetrics.Hits
		tam.Crits += spellTargetMetrics.Crits
		tam.Dodges += spellTargetMetrics.Dodges
		tam.Parries += spellTargetMetrics.Parries
		tam.Blocks += spellTargetMetrics.Blocks
		tam.Glances += spellTargetMetrics.Glances
		tam.Damage += spellTargetMetrics.TotalDamage
		tam.Threat += spellTargetMetrics.TotalThreat
		tam.Healing += spellTargetMetrics.TotalHealing
		tam.Shielding += spellTargetMetrics.TotalShielding
		tam.CastTime += spellTargetMetrics.TotalCastTime

		target := spell.Unit.AttackTables[i].Defender
		target.Metrics.dtps.Total += spellTargetMetrics.TotalDamage

		if spell.Unit.IsOpponent(target) {
			unitMetrics.dps.Total += spellTargetMetrics.TotalDamage
			unitMetrics.threat.Total += spellTargetMetrics.TotalThreat
		} else {
			unitMetrics.hps.Total += spellTargetMetrics.TotalHealing + spellTargetMetrics.TotalShielding
		}
	}
}

// This should be called at the end of each iteration, to include metrics from Pets in
// those of their owner.
// Assumes that doneIteration() has already been called on the pet metrics.
func (unitMetrics *UnitMetrics) AddFinalPetMetrics(petMetrics *UnitMetrics) {
	unitMetrics.dps.Total += petMetrics.dps.Total
}

func (unitMetrics *UnitMetrics) AddOOMTime(sim *Simulation, dur time.Duration) {
	unitMetrics.CharacterIterationMetrics.OOMTime += dur
	unitMetrics.MarkOOM(sim)
}
func (unitMetrics *UnitMetrics) MarkOOM(sim *Simulation) {
	if !unitMetrics.WentOOM {
		unitMetrics.WentOOM = true
		unitMetrics.FirstOOMTimestamp = sim.CurrentTime
	}
}

func (unitMetrics *UnitMetrics) UpdateDpasp(dpspSeconds float64) {
	// We store the total of seconds * spell power due to how DistributionMetrics work internally.
	unitMetrics.dpasp.Total += dpspSeconds
}

func (unitMetrics *UnitMetrics) reset() {
	unitMetrics.dps.reset()
	unitMetrics.dpasp.reset()
	unitMetrics.threat.reset()
	unitMetrics.dtps.reset()
	unitMetrics.tmi.reset()
	unitMetrics.tmiList = nil
	unitMetrics.hps.reset()
	unitMetrics.tto.reset()
	unitMetrics.CharacterIterationMetrics = CharacterIterationMetrics{}

	for _, resourceMetrics := range unitMetrics.resources {
		resourceMetrics.reset()
	}
}

// This should be called when a Sim iteration is complete.
func (unitMetrics *UnitMetrics) doneIteration(unit *Unit, sim *Simulation) {
	if unit.HasManaBar() {
		encounterDurationSeconds := sim.Duration.Seconds()
		timeToOOM := unitMetrics.FirstOOMTimestamp
		if !unitMetrics.WentOOM {
			// If we didn't actually go OOM in this iteration, infer TTO based on remaining mana.
			manaSpentPerSecond := (unitMetrics.ManaSpent - unitMetrics.ManaGained) / encounterDurationSeconds
			remainingTTO := DurationFromSeconds(unit.CurrentMana() / manaSpentPerSecond)
			timeToOOM = DurationFromSeconds(encounterDurationSeconds) + remainingTTO
			timeToOOM = min(timeToOOM, time.Minute*60)
		}

		if timeToOOM < 0 {
			// This happens sometimes when staying at 100% mana throughout the sim duration.
			timeToOOM = time.Minute * 60
		}

		unitMetrics.tto.Total = timeToOOM.Seconds()
		// Hack because of the way DistributionMetrics does its calculations.
		unitMetrics.tto.Total *= encounterDurationSeconds
	}

	if unitMetrics.isTanking {
		unitMetrics.tmi.Total = unitMetrics.calculateTMI(unit, sim)

		// Hack because of the way DistributionMetrics does its calculations.
		unitMetrics.tmi.Total *= sim.Duration.Seconds()
	}

	unitMetrics.dps.doneIteration(sim)
	unitMetrics.dpasp.doneIteration(sim)
	unitMetrics.threat.doneIteration(sim)
	unitMetrics.dtps.doneIteration(sim)
	unitMetrics.tmi.doneIteration(sim)
	unitMetrics.hps.doneIteration(sim)
	unitMetrics.tto.doneIteration(sim)

	unitMetrics.oomTimeSum += unitMetrics.OOMTime.Seconds()
	if unitMetrics.Died {
		unitMetrics.numItersDead++
	}
}

func (unitMetrics *UnitMetrics) calculateTMI(unit *Unit, sim *Simulation) float64 {
	if unit.Metrics.tmiList == nil || unitMetrics.tmiBin == 0 {
		return 0
	}

	bin := int(unitMetrics.tmiBin) // Seconds width for bin, default = 6
	firstEvent := 0                // Marks event at start of current bin
	ev := 0                        // Marks event at end of current bin
	lastEvent := len(unit.Metrics.tmiList)
	var buckets []float64

	// Traverse event array via marching time bins
	for tStep := 0; float64(tStep) < sim.Duration.Seconds()-float64(bin); tStep++ {

		// Increment event counter until we exceed the bin start
		for ; firstEvent < lastEvent && unit.Metrics.tmiList[firstEvent].Timestamp.Seconds() < float64(tStep); firstEvent++ {
		}

		// Increment event counter until we exceed the bin end
		for ; ev < lastEvent && unit.Metrics.tmiList[ev].Timestamp.Seconds() < float64(tStep+bin); ev++ {
		}

		if ev-firstEvent > 0 {
			sum := float64(0)

			// Add up everything in the bin
			for j := firstEvent; j < ev; j++ {
				sum += unit.Metrics.tmiList[j].WeightedDamage
			}

			//if sim.Log != nil {
			//	unit.Log(sim, "Bucket from %ds to %ds with events %d to %d totaled %f", tStep, tStep+bin, firstEvent, ev-1, sum)
			//}
			buckets = append(buckets, sum)
		} else { // an entire window with zero damage midfight still needs to be included
			if firstEvent < lastEvent {
				buckets = append(buckets, 0)
			}
		}

	}

	if buckets == nil {
		return 0
	}

	sum := float64(0)

	for i := 0; i < len(buckets); i++ {
		sum += math.Pow(math.E, buckets[i]*float64(10))
	}

	/* DEBUG LOGS
	if sim.Log != nil {
		raw_avg := float64(0)
		for i:=0; i < len(buckets); i++ {
			raw_avg += buckets[i]
		}
		raw_avg = raw_avg / float64(len(buckets))
		unit.Log(sim, "Sum of %d buckets was %f and raw mean bucket was %f", len(buckets), sum, raw_avg)
		unit.Log(sim, "TMI should be reported as %f", float64(10000) * math.Log(float64(1)/float64(len(buckets)) * sum))
	}
	*/

	return float64(10) * math.Log(float64(1)/float64(len(buckets))*sum)

	// 100000 / factor * ln ( Sum( p(window) * e ^ (factor * dmg(window) / hp ) ) )
	// factor = 10, multiplier of 100000 equivalent to 100% HP
	// Rescale to 100 = 100%

}

func (unitMetrics *UnitMetrics) ToProto() *proto.UnitMetrics {
	n := float64(unitMetrics.dps.n)
	protoMetrics := &proto.UnitMetrics{
		Dps:           unitMetrics.dps.ToProto(),
		Dpasp:         unitMetrics.dpasp.ToProto(),
		Threat:        unitMetrics.threat.ToProto(),
		Dtps:          unitMetrics.dtps.ToProto(),
		Tmi:           unitMetrics.tmi.ToProto(),
		Hps:           unitMetrics.hps.ToProto(),
		Tto:           unitMetrics.tto.ToProto(),
		SecondsOomAvg: unitMetrics.oomTimeSum / n,
		ChanceOfDeath: float64(unitMetrics.numItersDead) / n,
	}

	protoMetrics.Actions = make([]*proto.ActionMetrics, 0, len(unitMetrics.actions))
	for actionID, action := range unitMetrics.actions {
		protoMetrics.Actions = append(protoMetrics.Actions, action.ToProto(actionID))
	}

	protoMetrics.Resources = make([]*proto.ResourceMetrics, 0, len(unitMetrics.resources))
	for _, resource := range unitMetrics.resources {
		if resource.Events > 0 {
			protoMetrics.Resources = append(protoMetrics.Resources, resource.ToProto())
		}
	}

	return protoMetrics
}

type AuraMetrics struct {
	ID ActionID

	// Metrics for the current iteration.
	Uptime time.Duration
	Procs  int32

	// Aggregate values. These are updated after each iteration.
	aggregator
	procsSum int32
}

func (auraMetrics *AuraMetrics) reset() {
	auraMetrics.Uptime = 0
	auraMetrics.Procs = 0
}

// This should be called when a Sim iteration is complete.
func (auraMetrics *AuraMetrics) doneIteration() {
	auraMetrics.add(auraMetrics.Uptime.Seconds())
	auraMetrics.procsSum += auraMetrics.Procs
}

func (auraMetrics *AuraMetrics) ToProto() *proto.AuraMetrics {
	mean, stdev := auraMetrics.meanAndStdDev()

	return &proto.AuraMetrics{
		Id: auraMetrics.ID.ToProto(),

		UptimeSecondsAvg:   mean,
		UptimeSecondsStdev: stdev,
		ProcsAvg:           float64(auraMetrics.procsSum) / float64(auraMetrics.n),
	}
}

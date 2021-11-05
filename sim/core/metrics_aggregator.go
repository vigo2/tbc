package core

import (
	"math"
	"time"
)

// Returns a unique number based on an ActionID.
// This works by making item IDs negative to avoid collisions, and assumes
// there are no collisions with OtherID.
// Actual key values dont matter, just need something unique and fast to compute.
type ActionKey struct {
	ActionID int32
	Tag      int32
}
func NewActionKey(actionID ActionID, tag int32) ActionKey {
	return ActionKey{
		ActionID: int32(actionID.OtherID) + actionID.SpellID - actionID.ItemID,
		Tag: tag,
	}
}

type MetricsAggregator struct {
	// Duration of each iteration, in seconds
	encounterDuration float64

	startTime time.Time
	numIterations  int

	// Metrics for each Agent, for the current iteration
	agentIterations []AgentIterationMetrics

	// Aggregate values for each agent over all iterations
	agentAggregates []AgentAggregateMetrics
}

type AgentIterationMetrics struct {
	TotalDamage float64
	ManaSpent   float64
	DamageAtOOM float64
	OOMAt       time.Duration
}

type AgentAggregateMetrics struct {
	DpsSum        float64
	DpsSumSquared float64
	DpsMax        float64
	DpsHist       map[int32]int32 // rounded DPS to count

	NumOom      int
	OomAtSum    float64
	DpsAtOomSum float64

	Actions map[ActionKey]ActionMetric
}

type SimResult struct {
	ExecutionDurationMs int64
	Logs                string

	Agents []AgentResult
}

type AgentResult struct {
	DpsAvg   float64
	DpsStDev float64
	DpsMax   float64
	DpsHist  map[int32]int32 // rounded DPS to count

	NumOom      int
	OomAtAvg    float64
	DpsAtOomAvg float64

	Actions []ActionMetric
}

type ActionMetric struct {
	ActionID ActionID

	Tag int32

	Casts  int32
	Hits   int32
	Crits  int32
	Misses int32
	// Resists []int32   // Count of Resists

	Damage float64
}

func NewMetricsAggregator(numAgents int, encounterDuration float64) *MetricsAggregator {
	aggregator := &MetricsAggregator{
		encounterDuration: encounterDuration,
		startTime: time.Now(),
	}

	for i := 0; i < numAgents; i++ {
		aggregator.agentIterations = append(aggregator.agentIterations, AgentIterationMetrics{})
		aggregator.agentAggregates = append(aggregator.agentAggregates, AgentAggregateMetrics{})

		aggregator.agentAggregates[i].Actions = make(map[ActionKey]ActionMetric)
		aggregator.agentAggregates[i].DpsHist = make(map[int32]int32)
	}

	return aggregator
}

// Adds the results of an action to the aggregated metrics.
func (aggregator *MetricsAggregator) AddCastAction(cast DirectCastAction, castResults []DirectCastDamageResult) {
	actionID := cast.GetActionID()
	tag := cast.GetTag()

	actionKey := NewActionKey(actionID, tag)

	agentID := cast.GetCharacter().ID

	iterationMetrics := &aggregator.agentIterations[agentID]
	if !cast.castInput.IgnoreManaCost {
		iterationMetrics.ManaSpent += cast.castInput.ManaCost
	}

	aggregateMetrics := &aggregator.agentAggregates[agentID]
	actionMetrics, ok := aggregateMetrics.Actions[actionKey]

	if !ok {
		actionMetrics.ActionID = actionID
		actionMetrics.Tag = tag
	}

	actionMetrics.Casts++
	for _, result := range castResults {
		if result.Hit {
			actionMetrics.Hits++
		} else {
			actionMetrics.Misses++
		}

		if result.Crit {
			actionMetrics.Crits++
		}

		actionMetrics.Damage += result.Damage
		iterationMetrics.TotalDamage += result.Damage
	}

	aggregateMetrics.Actions[actionKey] = actionMetrics
}

func (aggregator *MetricsAggregator) MarkOOM(character *Character, oomAtTime time.Duration) {
	agentID := character.ID

	if aggregator.agentIterations[agentID].OOMAt == 0 {
		aggregator.agentIterations[agentID].DamageAtOOM = aggregator.agentIterations[agentID].TotalDamage
		aggregator.agentIterations[agentID].OOMAt = oomAtTime
	}
}

// This should be called when a Sim iteration is complete.
func (aggregator *MetricsAggregator) doneIteration() {
	aggregator.numIterations++

	// Loop for each agent
	for i, iterationMetrics := range aggregator.agentIterations {
		aggregateMetrics := &aggregator.agentAggregates[i]

		dps := iterationMetrics.TotalDamage / aggregator.encounterDuration
		// log.Printf("total: %0.1f, dur: %0.1f, dps: %0.1f", metrics.TotalDamage, aggregator.encounterDuration, dps)

		aggregateMetrics.DpsSum += dps
		aggregateMetrics.DpsSumSquared += dps * dps
		aggregateMetrics.DpsMax = MaxFloat(aggregateMetrics.DpsMax, dps)

		dpsRounded := int32(math.Round(dps/10) * 10)
		aggregateMetrics.DpsHist[dpsRounded]++

		if iterationMetrics.OOMAt > 0 {
			aggregateMetrics.NumOom++
			aggregateMetrics.OomAtSum += float64(iterationMetrics.OOMAt)
			aggregateMetrics.DpsAtOomSum += float64(iterationMetrics.DamageAtOOM) / float64(iterationMetrics.OOMAt.Seconds())
		}

		// Clear the iteration metrics
		aggregator.agentIterations[i] = AgentIterationMetrics{}
	}
}

func (aggregator *MetricsAggregator) getResult() SimResult {
	result := SimResult{}
	result.ExecutionDurationMs = time.Since(aggregator.startTime).Milliseconds()

	numIterations := float64(aggregator.numIterations)
	numAgents := len(aggregator.agentAggregates)

	result.Agents = make([]AgentResult, numAgents)
	for i := 0; i < numAgents; i++ {
		agentAggregate := &aggregator.agentAggregates[i]
		agentResult := &result.Agents[i]

		agentResult.DpsAvg = agentAggregate.DpsSum / numIterations
		agentResult.DpsStDev = math.Sqrt((agentAggregate.DpsSumSquared / numIterations) - (agentResult.DpsAvg * agentResult.DpsAvg))
		agentResult.DpsMax = agentAggregate.DpsMax
		agentResult.DpsHist = agentAggregate.DpsHist

		agentResult.NumOom = agentAggregate.NumOom
		if agentResult.NumOom > 0 {
			agentResult.OomAtAvg = agentAggregate.OomAtSum / float64(agentAggregate.NumOom)
			agentResult.DpsAtOomAvg = agentAggregate.DpsAtOomSum / float64(agentAggregate.NumOom)
		}

		for _, action := range agentAggregate.Actions {
			agentResult.Actions = append(agentResult.Actions, action)
		}
	}

	return result
}
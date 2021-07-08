package types

type ComponentState string

const (
	ComponentStateApprove             ComponentState = "operational"
	ComponentStateDegradedPerformance ComponentState = "degraded_performance"
)

func (ComponentState) Values() []ComponentState {
	return []ComponentState{
		"operational",
		"degraded_performance",
	}
}

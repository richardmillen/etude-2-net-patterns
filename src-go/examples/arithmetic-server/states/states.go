package states

var (
	NumState = &fsm.State{
		Name: "number",
		Events: []*fsm.Event{
			{
				Input: msgs.Num,
				MoveTo: []*fsm.State{
					OpState,
					CalcState,
				},
			},
		},
	}
	OpState = &fsm.State{
		Name: "operator",
		Events: []*fsm.Event{
			{
				Input:  msgs.Op,
				MoveTo: []*fsm.State{NumState},
			},
		},
	}
	CalcState = &fsm.State{
		Name: "calculate",
		Events: []*fsm.Event{
			{
				Input: msgs.Num,
			},
		},
	}
)

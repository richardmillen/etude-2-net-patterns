package msgs

const opersExpr = `[\+\-/\*]`

// Operator represents a hypothetical message type from client to server.
// it embeds an fsm.Regex type that would provide regular expression matching.
// the operator type comes pre-baked with 'opersExpr' that matches basic arithmetic
// operators.
type Operator struct {
	fsm.Regex
}

// Oper is a convenience method that performs the arithmetical operation
// on operands 'a' and 'b'.
func (o Operator) Oper(a float32, b float32) (result float32) {
	// TODO: perform relevant operation:
	result = a + b
	return
}

var (
	// Num is a numeric operand.
	Num = &fsm.Float{
		Hint: "numeric operand",
	}
	Op     = &Operator{}
	Equals = &fsm.String{
		Hint:  "equals sign",
		Match: "=",
	}
	Any = &fsm.Any{}
)

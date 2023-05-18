package plonky2_verifier

import (
	"fmt"
	. "gnark-plonky2-verifier/field"
	"regexp"
	"strconv"
)

var mulExtensionGateRegex = regexp.MustCompile("MulExtensionGate { num_ops: (?P<numOps>[0-9]+) }")

func deserializeMulExtensionGate(parameters map[string]string) gate {
	// Has the format "MulExtensionGate { num_ops: 13 }"
	numOps, hasNumOps := parameters["numOps"]
	if !hasNumOps {
		panic("Missing field num_ops in MulExtensionGate")
	}

	numOpsInt, err := strconv.Atoi(numOps)
	if err != nil {
		panic("Invalid num_ops field in MulExtensionGate")
	}

	return NewMultiplicationExtensionGate(uint64(numOpsInt))
}

type MultiplicationExtensionGate struct {
	numOps uint64
}

func NewMultiplicationExtensionGate(numOps uint64) *MultiplicationExtensionGate {
	return &MultiplicationExtensionGate{
		numOps: numOps,
	}
}

func (g *MultiplicationExtensionGate) Id() string {
	return fmt.Sprintf("ArithmeticExtensionGate { num_ops: %d }", g.numOps)
}

func (g *MultiplicationExtensionGate) wiresIthMultiplicand0(i uint64) Range {
	return Range{3 * D * i, 3*D*i + D}
}

func (g *MultiplicationExtensionGate) wiresIthMultiplicand1(i uint64) Range {
	return Range{3*D*i + D, 3*D*i + 2*D}
}

func (g *MultiplicationExtensionGate) wiresIthOutput(i uint64) Range {
	return Range{3*D*i + 2*D, 3*D*i + 3*D}
}

func (g *MultiplicationExtensionGate) EvalUnfiltered(p *PlonkChip, vars EvaluationVars) []QuadraticExtension {
	const0 := vars.localConstants[0]

	constraints := []QuadraticExtension{}
	for i := uint64(0); i < g.numOps; i++ {
		multiplicand0 := vars.GetLocalExtAlgebra(g.wiresIthMultiplicand0(i))
		multiplicand1 := vars.GetLocalExtAlgebra(g.wiresIthMultiplicand1(i))
		output := vars.GetLocalExtAlgebra(g.wiresIthOutput(i))

		mul := p.qeAPI.MulExtensionAlgebra(multiplicand0, multiplicand1)
		computed_output := p.qeAPI.ScalarMulExtensionAlgebra(const0, mul)

		diff := p.qeAPI.SubExtensionAlgebra(output, computed_output)
		for j := 0; j < D; j++ {
			constraints = append(constraints, diff[j])
		}
	}

	return constraints
}

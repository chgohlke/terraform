package terraform

import (
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/dag"
)

// DestroyPlanGraphBuilder implements GraphBuilder and is responsible for
// planning a pure-destroy.
//
// Planning a pure destroy operation is simple because we can ignore most
// ordering configuration and simply reverse the state.
type DestroyPlanGraphBuilder struct {
	// Module is the root module for the graph to build.
	Module *module.Tree

	// State is the current state
	State *State

	// Targets are resources to target
	Targets []string
}

// See GraphBuilder
func (b *DestroyPlanGraphBuilder) Build(path []string) (*Graph, error) {
	return (&BasicGraphBuilder{
		Steps:    b.Steps(),
		Validate: true,
	}).Build(path)
}

// See GraphBuilder
func (b *DestroyPlanGraphBuilder) Steps() []GraphTransformer {
	concreteResource := func(a *NodeAbstractResource) dag.Vertex {
		return &NodePlanDestroyableResource{
			NodeAbstractResource: a,
		}
	}

	steps := []GraphTransformer{
		// Creates all the nodes represented in the state.
		&StateTransformer{
			Concrete: concreteResource,
			State:    b.State,
		},

		// Target
		&TargetsTransformer{Targets: b.Targets},

		// Attach the configuration to any resources
		&AttachResourceConfigTransformer{Module: b.Module},

		// Single root
		&RootTransformer{},
	}

	return steps
}

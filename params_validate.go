// Copyright (c) triagegeist authors: Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov.
// Licensed under the EUPL.

package triagegeist

import "github.com/olaflaitinen/triagegeist/validate"

// ValidateParamsExternal runs the validate package's Params check on p.
// Use when you want the same validation logic as validate.Params without
// constructing validate.ParamsLike manually.
func ValidateParamsExternal(p Params) bool {
	pl := validate.ParamsLike{
		VitalWeights:   p.VitalWeights,
		MaxResources:   p.MaxResources,
		ResourceWeight: p.ResourceWeight,
		T1:             p.T1, T2: p.T2, T3: p.T3, T4: p.T4,
	}
	return validate.ParamsValid(pl)
}

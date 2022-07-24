package common

import spectypes "github.com/lavanet/lava/x/spec/types"

func CreateMockSpec() spectypes.Spec {
	specName := "mockSpec"
	spec := spectypes.Spec{}
	spec.Name = specName
	spec.Index = specName
	spec.Enabled = true
	apiInterface := spectypes.ApiInterface{Interface: "mockInt", Type: "get"}
	spec.Apis = append(spec.Apis, spectypes.ServiceApi{Name: specName + "API", ComputeUnits: 100, Enabled: true, ApiInterfaces: []spectypes.ApiInterface{apiInterface}})

	return spec
}
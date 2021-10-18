package knowledge

type ClusterType struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ClusterVersion struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ClusterComponent struct {
	ComponentType string `json:"componentType"`
	ComponentName string `json:"componentName"`
}

type ClusterTypeSpec struct {
	ClusterType  ClusterType          `json:"clusterType"`
	VersionSpecs []ClusterVersionSpec `json:"versionSpecs"`
}

func (s *ClusterTypeSpec) GetVersionSpec(versionCode string) (versionSpec *ClusterVersionSpec) {
	for i := range s.VersionSpecs {
		if s.VersionSpecs[i].ClusterVersion.Code == versionCode {
			return &s.VersionSpecs[i]
		}
	}
	return nil
}

type ClusterVersionSpec struct {
	ClusterVersion ClusterVersion         `json:"clusterVersion"`
	ComponentSpecs []ClusterComponentSpec `json:"componentSpecs"`
}

func (s *ClusterVersionSpec) GetComponentSpec(componentType string) (componentSpec *ClusterComponentSpec) {
	for i := range s.ComponentSpecs {
		if s.ComponentSpecs[i].ClusterComponent.ComponentType == componentType {
			return &s.ComponentSpecs[i]
		}
	}
	return nil
}

type ComponentPortConstraint struct {
	Start int `json:"portRangeStart"`
	End   int `json:"portRangeEnd"`
	Count int `json:"portCount"`
}

type ClusterComponentSpec struct {
	ClusterComponent    ClusterComponent        `json:"clusterComponent"`
	ComponentConstraint ComponentConstraint     `json:"componentConstraint"`
	PortConstraint      ComponentPortConstraint `json:"compentPortConstraint"`
}

type ComponentConstraint struct {
	ComponentRequired       bool     `json:"componentRequired"`
	SuggestedNodeQuantities []int    `json:"suggestedNodeQuantities"`
	AvailableSpecCodes      []string `json:"availableSpecCodes"`
	MinZoneQuantity         int      `json:"minZoneQuantity"`
}
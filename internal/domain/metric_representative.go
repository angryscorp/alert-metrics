package domain

type MetricRepresentative struct {
	Type  MetricType
	Name  string
	Value string
}

func (m MetricRepresentative) String() string {
	return string(m.Type) + "." + m.Name + ": " + m.Value
}

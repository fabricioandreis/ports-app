package domain

type Port struct {
	ID          string
	Code        string
	Name        string
	City        string
	Country     string
	Province    string
	Timezone    string
	Alias       []string
	Coordinates []float32
	Regions     []string
	Unlocs      []string
}

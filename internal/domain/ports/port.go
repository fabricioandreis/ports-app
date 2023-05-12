package ports

type Port struct {
	ID          string
	Code        string
	Name        string
	City        string
	Province    string
	Country     string
	Timezone    string
	Alias       []string
	Coordinates []float32
	Regions     []string
	Unlocs      []string
}

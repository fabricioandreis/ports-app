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
	Coordinates Coordinates
	Regions     []string
	Unlocs      []string
}

type Coordinates struct {
	Lat, Long float32
}

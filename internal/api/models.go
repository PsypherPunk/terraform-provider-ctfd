package api

// Tag -
type Tag struct {
	Value string `json:"value"`
}

// Challenge -
type Challenge struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Solves     uint   `json:"solves"`
	SolvedByMe bool   `json:"solved_by_me"`
	Category   string `json:"category"`
	Tags       []Tag  `json:"tags"`
	Template   string `json:"template"`
	Script     string `json:"script"`
}

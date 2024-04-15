package holiday

type Holiday struct {
	Schema string   `json:"$schema"`
	ID     string   `json:"$id"`
	Year   int      `json:"year"`
	Papers []string `json:"papers"`
	Days   []*Day   `json:"days"`
}
type Day struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	IsOffDay bool   `json:"isOffDay"`
}

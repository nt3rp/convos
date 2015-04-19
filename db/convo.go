package db

type Convo struct {
	Id        int      `json:"id"`
	Sender    int      `json:"sender"`
	Recipient int      `json:"recipient"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Status    string   `json:"status"`
	Children  []*Convo `json:"replies"`
}

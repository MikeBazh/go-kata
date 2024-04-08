package PetModel

type Pet struct {
	Category struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status    string   `json:"status"`
	Tags      []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

package opencars

type Transport struct {
	ID                  int    `json:"id"`
	RegistrationAddress string `json:"registration_address"`
	RegistrationCode    int    `json:"registration_code"`
	Registration        string `json:"registration"`
	Date                string `json:"date"`
	Model               string `json:"model"`
	Year                int    `json:"year"`
	Color               string `json:"color"`
	Kind                string `json:"kind"`
	Body                string `json:"body"`
	Fuel                string `json:"fuel"`
	Capacity            int    `json:"capacity"`
	Weight              int    `json:"own_weight"`
	Number              string `json:"number"`
}

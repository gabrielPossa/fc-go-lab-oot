package cep

type CEP struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
	Erro       string `json:"erro"`
}

type RequestCEP struct {
	Cep string `json:"cep"`
}

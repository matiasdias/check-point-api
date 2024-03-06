package models

// UpdatePassWord responsável pelo mapeamento dos campos de atualização da senha
type UpdatePassWord struct {
	CurrentPassWord string `json:"current_password,omitempty"`
	NewPassWord     string `json:"new_passWord,omitempty"`
}

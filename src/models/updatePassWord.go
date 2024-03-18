package models

// UpdatePassWord responsável pelo mapeamento dos campos de atualização da senha
type UpdatePassWord struct {
	NewPassWord string `json:"new_passWord,omitempty"`
}

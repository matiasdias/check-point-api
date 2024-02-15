package models

// UpdatePassWord responsavel por informar a senha atual e a nova senha
type UpdatePassWord struct {
	CurrentPassWord string `json:"senhaAtual,omitempty"`
	NewPassWord     string `json:"newPassWord,omitempty"`
}

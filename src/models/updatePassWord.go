package models

type UpdatePassWord struct {
	CurrentPassWord string `json:"senhaAtual,omitempty"`
	NewPassWord     string `json:"newPassWord,omitempty"`
}

package domain

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	ProviderID   string
	ProviderName string
}

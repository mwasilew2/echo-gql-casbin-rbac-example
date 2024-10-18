// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Account struct {
	Ulid string `json:"ulid"`
	Name string `json:"name"`
	// The account's ID
	ID uint `gorm:"primaryKey"`
}

type Mutation struct {
}

type Namespace struct {
	Ulid    string   `json:"ulid"`
	Name    string   `json:"name"`
	Account *Account `json:"account"`
}

type NewAccount struct {
	Name string `json:"name"`
}

type NewNamespace struct {
	Name string `json:"name"`
}

type NewStack struct {
	Name        string `json:"name"`
	NamespaceID string `json:"namespaceId"`
}

type Query struct {
}

type Stack struct {
	Ulid      string     `json:"ulid"`
	Name      string     `json:"name"`
	Namespace *Namespace `json:"namespace"`
	Account   *Account   `json:"account"`
}

type User struct {
	Ulid     string   `json:"ulid"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Account  *Account `json:"account"`
	// The account's ID
	AccountID uint `json:"-"`
	// The user's ID
	ID uint `gorm:"primaryKey"`
}

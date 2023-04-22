package entities

// Service interface would be here if not defined in proto

type Repository interface {
	CreateLink(original string, hashed string) string
	ReturnLink(hashed string) string
	CheckIfOriginalExists(original string) error
}

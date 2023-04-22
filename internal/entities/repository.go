package entities

import "context"

// Service interface would be here if not defined in proto

type Repository interface {
	CreateLink(ctx context.Context, hashed string, original string) error
	ReturnLink(ctx context.Context, hashed string) (string, error)
	CheckIfHashedExists(ctx context.Context, hashed string) error
}

package outbound

type TokenGenerator interface {
	Generate(string) (string, error)
}

type TokenValidator interface {
	Validate(string) (string, error)
}

type TokenGeneratorValidator interface {
	TokenGenerator
	TokenValidator
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

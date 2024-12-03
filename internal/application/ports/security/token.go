package security

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

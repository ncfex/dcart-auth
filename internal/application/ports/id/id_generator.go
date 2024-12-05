package id

type IDGenerator interface {
	GenerateFromData(data []byte) string
}

package domain

type Encrypter interface {
	Encrypt(data []byte) ([]byte, error)
}

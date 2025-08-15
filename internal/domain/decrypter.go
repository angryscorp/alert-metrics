package domain

type Decrypter interface {
	Decrypt(data []byte) ([]byte, error)
}

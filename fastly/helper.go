package fastly

import "github.com/fastly/compute-sdk-go/secretstore"

func GetSecretStoreKey(storeName string, keyName string) (string, error) {
	st, err := secretstore.Open(storeName)
	if err != nil {
		return "", err
	}

	s, err := st.Get(keyName)
	if err != nil {
		return "", err
	}

	v, err := s.Plaintext()
	if err != nil {
		return "", err
	}

	return string(v), err
}

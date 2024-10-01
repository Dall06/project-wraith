package alchemy

import (
	"encoding/json"
	"errors"
)

func StructIntoString(entity interface{}, secret string) (string, error) {
	data, err := json.Marshal(entity)
	if err != nil {
		return "", err
	}

	encrypted, err := Encrypt(string(data), secret)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func StringToStruct(encrypted string, secret *string, entity interface{}) error {
	if secret == nil {
		return errors.New("secret cannot be nil")
	}

	decrypted, err := Decrypt(encrypted, *secret)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(decrypted), entity)
	if err != nil {
		return err
	}

	return nil
}

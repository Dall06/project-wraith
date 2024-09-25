package alchemy

import (
	"reflect"
)

func Transmutation(entity interface{}, secret string) error {
	val := reflect.ValueOf(entity).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			encrypted, err := Encrypt(secret, field.String())
			if err != nil {
				return err
			}
			field.SetString(encrypted)
		}
	}
	return nil
}

func Revert(entity interface{}, secret string) error {
	val := reflect.ValueOf(entity).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			decrypted, err := Decrypt(field.String(), secret)
			if err != nil {
				return err
			}
			field.SetString(decrypted)
		}
	}
	return nil
}

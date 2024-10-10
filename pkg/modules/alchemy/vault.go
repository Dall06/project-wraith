package alchemy

import (
	"os"
)

type VaultKeeper interface {
	Secure(filePath string, replace bool) error
	ObliterateFile(filePath string) error
	Release(filePath string) (string, error)
}

type ArcaneVault struct {
	secret string
}

func (av *ArcaneVault) Secure(filePath string, replace bool) error {
	scroll, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	securedScroll, err := Encrypt(string(scroll), av.secret)
	if err != nil {
		return err
	}

	if replace {
		return os.WriteFile(filePath, []byte(securedScroll), 0644)
	}

	newFilePath := filePath + ".secured"
	return os.WriteFile(newFilePath, []byte(securedScroll), 0644)
}

func (av *ArcaneVault) Release(filePath string) (string, error) {
	securedScroll, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return Decrypt(string(securedScroll), av.secret)
}

func (av *ArcaneVault) ObliterateFile(filePath string) error {
	return os.Remove(filePath)
}

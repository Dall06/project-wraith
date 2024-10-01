package core

import (
	"errors"
	"project-wraith/pkg/modules/lics"
	"time"
)

func Activate(repo lics.LicenseRepository, licenseKey string) error {
	lic := &lics.License{
		LicenseKey: licenseKey,
	}

	result, err := repo.Get(*lic)
	if err != nil {
		return err
	}

	if !result.IsActive && result.ExpiryDate.After(time.Now()) {
		return errors.New("license is either inactive or expired")
	}

	return nil
}

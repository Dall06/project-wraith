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

	// Check if result is nil
	if result == nil {
		return errors.New("license not found")
	}

	if !result.IsActive || result.ExpiryDate.Before(time.Now()) {
		return errors.New("license is either inactive or expired")
	}

	return nil
}

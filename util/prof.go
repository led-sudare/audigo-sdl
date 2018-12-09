package util

import (
	"github.com/pkg/profile"
)

func GetProfile() *profile.Profile {
	return profile.Start(profile.ProfilePath(".")).(*profile.Profile)
}

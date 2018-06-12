package environment

import "fmt"

type Environment struct {
	Timezone string
	Locale   string
	Lang     string
}

func (e *Environment) GetTimezone() string {
	return fmt.Sprintf("%s", e.Timezone)
}

func (e *Environment) GetLang() string {
	return ""
}

func (e *Environment) GetLocale() string {
	return ""
}

func newEnvironment() *Environment {
	return &Environment{
		Timezone: "",
		Locale:   "en_US.UTF-8",
		Lang:     "en_US.UTF-8",
	}
}

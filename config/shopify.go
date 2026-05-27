package config

type Domain struct {
	// Console   string `mapstructure:"console" json:"console" yaml:"console"`
	Artist   string `mapstructure:"artist" json:"artist" yaml:"artist"`
	Referral string `mapstructure:"referral" json:"referral" yaml:"referral"`
	Admin    string `mapstructure:"admin" json:"admin" yaml:"admin"`
	View     string `mapstructure:"view" json:"view" yaml:"view"`
}

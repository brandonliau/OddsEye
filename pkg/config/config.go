package config

type Config interface {
	load(string) error
	validate() error
}

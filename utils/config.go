package utils

type Config struct {
	screenHeight int
	screenWidth  int
}

func NewConfig(sh, sw int) *Config {
	return &Config{
		screenHeight: sh,
		screenWidth:  sw,
	}
}

package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type IpClientConfig struct {
	Uri string `mapstructure:"URI"`
}

type DoConfig struct {
	Token string `mapstructure:"TOKEN"`
}

type Config struct {
	DoConfig DoConfig       `mapstructure:"DIGITAL_OCEAN"`
	IpClient IpClientConfig `mapstructure:"IP_CLIENT"`
	Interval int            `mapstructure:"INTERVAL"`
	Domains  []string       `mapstructure:"DOMAINS"`
}

func LoadConfig(cfgFile string) *Config {
	log.Default().Printf("Loading config from %s", cfgFile)
	config := &Config{}
	v := viper.New()
	v.AutomaticEnv()
	// Load config from file
	v.SetConfigFile(cfgFile)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	v.SetConfigName("")
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); (!ok) && err.Error() != "open .env: no such file or directory" {
			log.Fatal("Cannot read cofiguration")
		}
	}

	// Unmarshal config
	err = v.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
	}

	mergeEnv(config)

	return config

}

func mergeEnv(c *Config) {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		return
	}

	if value, ok := myEnv["DO_TOKEN"]; ok {
		c.DoConfig.Token = value
	}

	if value, ok := myEnv["IP_CLIENT_URI"]; ok {
		c.IpClient.Uri = value
	}

}

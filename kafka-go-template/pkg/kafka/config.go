package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// These represents a limited set of kafka configurations used by the Customers consumer/producer applications.
//
//	Additional properties (supported by rdkafka) can be added into the configuration with the properties fields
type Config struct {
	BootstrapServers string `mapstructure:"bootstrap.servers"`
	ClientId         string `mapstructure:"client.id"`

	SchemaRegistryUrl      string  `mapstructure:"schema.registry.url"`
	SchemaRegistryUser     *string `mapstructure:"schema.registry.user"`
	SchemaRegistryPassword *string `mapstructure:"schema.registry.password"`

	TopicName string `mapstructure:"topic.name"`

	// Additional broker properties
	Properties map[string]string `mapstructure:"properties"`
	// Consumer properties
	Consumer map[string]string `mapstructure:"consumer"`
	// Producer properties
	Producer map[string]string `mapstructure:"producer"`
}

const Prefix = "ASHKAL"

func InitConfig(configname string) (Config, error) {
	// Initialization
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigName(configname)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvPrefix(Prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	envVars := []string{"bootstrap.servers", "client.id", "schema.registry.url",
		"schema.registry.user", "schema.registry.password"}
	for _, envName := range envVars {
		v.BindEnv(envName)
	}

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config := Config{}
	err = v.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Printf("%+v\n", config)

	config.Properties = loadEnvProperties("KAFKA", config.Properties)
	config.Producer = loadEnvProperties("PRODUCER", config.Producer)
	config.Consumer = loadEnvProperties("CONSUMER", config.Consumer)

	return config, err
}

// loads environment properties matching into a map
func loadEnvProperties(group string, props map[string]string) map[string]string {
	var result [][]string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.Contains(pair[0], Prefix+"_"+group+"_") {
			result = append(result, []string{strings.Replace(pair[0], Prefix+"_"+group+"_", "", 1), pair[1]})
		}
	}
	if len(result) > 0 {
		if props == nil {
			props = make(map[string]string)
		}
		for _, v := range result {
			lcase := strings.ToLower(v[0])
			propertyName := strings.ReplaceAll(lcase, "_", ".")
			props[propertyName] = v[1]
		}
	}

	return props
}

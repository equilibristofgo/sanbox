package properties

import (
	"errors"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/common/properties/models"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type ApplicationProperties struct {
	Micro       models.Micro       `yaml:"micro"`
	HttpHandler models.HttpHandler `yaml:"http_handler"`
	Repository  models.Repository  `yaml:"infrastructure"`
	Auth        models.Auth        `yaml:"auth"`
}

func (p *ApplicationProperties) GetConfiguration(applicationProperties string) {
	//Define Default value
	p.Micro = p.Micro.DefaultTag()
	p.HttpHandler = p.HttpHandler.DefaultTag()
	p.Repository = p.Repository.DefaultTag()
	p.Auth = p.Auth.DefaultTag()

	//Read Yaml Properties
	v := setViperByEnvironment(applicationProperties)
	mapProperties := v.AllSettings()
	yamlString, _ := yaml.Marshal(mapProperties)
	yaml.Unmarshal(yamlString, &p)
}

func setViperByEnvironment(fileName string) *viper.Viper {
	var effectiveFileName string
	v := viper.New()
	if os.Getenv("ENVIRONMENT") != "" {
		substr := strings.Split(fileName, ".yaml")
		effectiveFileName = substr[0] + "_" + strings.ToLower(os.Getenv("ENVIRONMENT")) + ".yaml"
	} else {
		effectiveFileName = fileName
	}
	err := existsFile(effectiveFileName)
	if err != nil {
		panic("properties file " + effectiveFileName + " not found ")
	}
	v.SetConfigType("yaml")
	v.SetConfigFile(effectiveFileName)
	v.AutomaticEnv()
	envReplacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(envReplacer)
	v.ReadInConfig()

	return v
}

func existsFile(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return err
	}
	return errors.New("Unknown Error")
}

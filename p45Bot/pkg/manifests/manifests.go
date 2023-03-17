package manifests

import (
	"github.com/spf13/viper"
)

func ReadManifest(path string) {
	
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")
}

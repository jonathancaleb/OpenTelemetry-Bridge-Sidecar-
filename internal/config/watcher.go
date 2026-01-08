package config

import (
    "fmt"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

func WatchConfigFile(filePath string, onChange func()) error {
    if filePath == "" {
        return fmt.Errorf("empty config file path")
    }

    viper.SetConfigFile(filePath)

    if err := viper.ReadInConfig(); err != nil {
        return err
    }

    viper.OnConfigChange(func(e fsnotify.Event) {
        if onChange != nil {
            onChange()
        }
    })

    viper.WatchConfig()
    return nil
}
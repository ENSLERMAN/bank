/*
 Банковская система на обычном REST ( монолит ).
 */
package main

import (
	"flag"
	"log"
	"github.com/BurntSushi/toml"
	"github.com/ENSLERMAN/soft-eng/project/internal/app/apiserver"
)

var (
	configPath string
)

// в init путь до файла с конфигурацией.
func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
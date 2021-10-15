package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"dev11/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/apiserver.json", "Path to JSON config file")
}

func main() {
	flag.Parse()
	// Парсинг конфигурационного файла из json
	config := apiserver.NewConfig()
	// ReadFile читает файл, путь до которого передается в аргументе. Если путь не абсолютный, то отсчитываться он будет от текущей рабочей
	// директории
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal производит перенос значений параметров конфигурации в экземпляр структуры Config.
	// При определении структурного типа Config необходимо состыковать поля структуры и соответствующие поля в json-файле
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	// Инициализация сервера (создание экземпляра APIServer)
	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

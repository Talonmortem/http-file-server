package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

/*type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Storage struct {
		UploadDir   string `yaml:"upload_dir"`
		PublicDir   string `yaml:"public_dir"`
		MaxFileSize int    `yaml:"max_file_size"`
	} `yaml:"storage"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	JWT struct {
		SecretKey string `yaml:"secret_key"`
		ExpiresIn int    `yaml:"expires_in"`
	}
}*/

type Config struct {
	Server   Server
	Storage  Storage
	JWT      JWT
	Database Database
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Storage struct {
	UploadDir   string `yaml:"upload_dir"`
	PublicDir   string `yaml:"public_dir"`
	TemplateDir string `yaml:"template_dir"`
	WebDir      string `yaml:"web_dir"`
}

type JWT struct {
	SecretKey string `yaml:"secret_key"`
	ExpiresIn int    `yaml:"expires_in"`
}

type Database struct {
	Path string `yaml:"path"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(path string) (*Config, error) {
	// Читаем содержимое файла
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Ошибка чтения конфигурационного файла: %v", err)
		return nil, err
	}

	// Создаем структуру для хранения настроек
	var cfg Config

	// Декодируем YAML в структуру
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Ошибка разбора YAML: %v", err)
		return nil, err
	}

	log.Println("Config file loaded successfully", cfg)

	return &cfg, nil
}

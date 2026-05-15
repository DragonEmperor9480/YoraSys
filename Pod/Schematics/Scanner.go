package schematics

type AnamolyReader struct {
	ID          int      `yaml:"id"`
	Name        string   `yaml:"name"`
	Category    string   `yaml:"category"`
	Description string   `yaml:"description"`
	Paths       []string `yaml:"paths"`
}

type Registry struct {
	Schema   Schema          `yaml:"schema"`
	Platform string          `yaml:"platform"`
	Caches   []AnamolyReader `yaml:"caches"`
}

type Schema struct {
	Name    string  `yaml:"name"`
	Version float32 `yaml:"version"`
}


package logic

type Provider struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Env         map[string]string `json:"env"`
}

type Config struct {
	Providers    []Provider `json:"providers"`
	CurrentIndex int        `json:"currentIndex"`
}

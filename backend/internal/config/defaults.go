package config

func (c *Config) SetDefaults() {

	c.Server.Address = ":8080"
	c.MaxAPI.BaseURL = "https://platform-api.max.ru"
	c.App.JWTSecret = "your-secret-key-change-in-production"
	c.App.JWTTTL = 3600
	c.App.RefreshTTL = 604800
	c.App.WebSocketPath = "/ws"
	c.App.MaxSessionSize = 20
}

package config

func (c *Config) SetDefaults() {

	c.Server.Address = ":8080"
	// Telegram Bot API использует стандартный endpoint https://api.telegram.org
	c.App.JWTSecret = "your-secret-key-change-in-production"
	c.App.JWTTTL = 900        // 15 minutes for access token
	c.App.RefreshTTL = 604800 // 7 days for refresh token
	c.App.WebSocketPath = "/ws"
	c.App.MaxSessionSize = 20
}

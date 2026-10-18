package models



type Config struct {
	DBDriver			string		`mapstructure:"DB_DRIVER"`
	DBSource			string		`mapstructure:"DB_SOURCE"`
	ServerAddress		string		`mapstructure:"SERVER_ADDRESS"`
	SecretKey			string		`mapstructure:"SECRET_KEY"`
	AccessTokenMinutes	int			`mapstructure:"ACCESS_TOKEN_MINUTES"`
	CloudinaryUrl		string		`mapstructure:"CLOUDINARY_URL"`
}


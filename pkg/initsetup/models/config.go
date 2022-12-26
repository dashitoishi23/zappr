package initsetupmodels

type Config struct {
	DatabaseHost     string `json:"databaseHost"`
	DatabaseName     string `json:"databaseName"`
	DatabaseUser     string `json:"databaseUser"`
	DatabasePassword string `json:"databasePassword"`
	DatabaseSSLMode  string `json:"databaseSSLMode"`
	DatabasePort     int    `json:"databasPort"`
	JWTSigningKey    string `json:"jwtSigningKey"`
	RedisHost        string `json:"redisHost"`
	RedisPassword    string `json:"redisPassword"`
}
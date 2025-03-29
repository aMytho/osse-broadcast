package config

import "os"

type OsseConfig struct {
	HttpHost         string
	RedisHost        string
	OsseClientOrigin string
}

func GetOsseConfig() OsseConfig {
	httpHost := getEnvVar("OSSE_BROADCAST_HOST")
	redisHost := getEnvVar("OSSE_REDIS_HOST")
	osseClientOrigin := getEnvVar("OSSE_ALLOWED_ORIGIN")

	return OsseConfig{httpHost, redisHost, osseClientOrigin}
}

func getEnvVar(key string) string {
	result, varExists := os.LookupEnv(key)

	if !varExists {
		println("The environment variable " + key + " was not set. Please set this var in the osse config file.")
		println("Osse Broadcast is shutting down!")
		os.Exit(1)
	}

	return result
}

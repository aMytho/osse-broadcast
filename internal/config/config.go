package config

import "os"

type OsseConfig struct {
	HttpHost  string
	RedisHost string
}

func GetOsseConfig() OsseConfig {
	httpHost := getEnvVar("OSSE_BROADCAST_HOST")
	redisHost := getEnvVar("OSSE_REDIS_HOST")

	return OsseConfig{httpHost, redisHost}
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

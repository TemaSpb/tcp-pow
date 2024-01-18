package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost              string
	ServerPort              int64
	HashcashZeros           int
	HashcashChallengeLength int
	HashcashChallenge       string
	HashcashDuration        int64
	HashcashMaxIterations   int
}

func LoadConfig(confPath string) (*Config, error) {
	var myEnvs map[string]string

	myEnvs, err := godotenv.Read(confPath)
	if err != nil {
		return nil, err
	}

	serverPort, err := strconv.ParseInt(myEnvs["SERVER_PORT"], 10, 64)
	if err != nil {
		return nil, err
	}

	hashcashZeros, err := strconv.Atoi(myEnvs["HASHCASH_ZEROS"])
	if err != nil {
		return nil, err
	}

	hashcashChallengeLength, err := strconv.Atoi(myEnvs["HASHCASH_CHALLENGE_LENGTH"])
	if err != nil {
		return nil, err
	}

	hashcashDuration, err := strconv.ParseInt(myEnvs["HASHCASH_DURATION"], 10, 64)
	if err != nil {
		return nil, err
	}

	hashcashMaxIterations, err := strconv.Atoi(myEnvs["HASHCASH_MAX_ITERATIONS"])
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerHost:              myEnvs["SERVER_HOST"],
		ServerPort:              serverPort,
		HashcashZeros:           hashcashZeros,
		HashcashChallenge:       myEnvs["HASHCASH_CHALLENGE"],
		HashcashChallengeLength: hashcashChallengeLength,
		HashcashDuration:        hashcashDuration,
		HashcashMaxIterations:   hashcashMaxIterations,
	}, nil
}

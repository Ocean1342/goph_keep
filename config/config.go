package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL string        `env:"DATABASE_URI"`
	RunAddr     string        `env:"RUN_ADDRESS"`
	SecretKey   string        `env:"SECRET_KEY"`
	TokenTTL    time.Duration `env:"TOKEN_TTL"`
	TLSEnabled  bool          `env:"TLS_ENABLED"`
	CertFile    string        `env:"CERT_FILE"`
	KeyFile     string        `env:"KEY_FILE"`
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Errorf("Error loading .env file")
	}
	ttl, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		log.Errorf("could not set token ttl.err:%s", err)
	}
	if ttl == 0 {
		ttl = 24
	}
	dbURL, ok := os.LookupEnv("DATABASE_URI")
	if !ok {
		log.Fatalf("DATABASE_URI env variable not set")
	}
	runAddr, ok := os.LookupEnv("RUN_ADDRESS")
	if !ok {
		log.Fatalf("RUN_ADDRESS env variable not set")
	}
	secretKey, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatalf("SECRET_KEY env variable not set")
	}
	v, ok := os.LookupEnv("TLS_ENABLED")
	if !ok {
		log.Fatalf("TLS_ENABLED env variable not set")
	}
	tlsEnabled, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("TLS_ENABLED env variable not set")
	}
	tlsCert, ok := os.LookupEnv("CERT_FILE")
	if !ok && tlsEnabled || tlsEnabled && tlsCert == "" {
		log.Fatalf("CERT_FILE env variable not set when TLS is enabled")
	}
	tlsCert = ""
	tlsKey, ok := os.LookupEnv("KEY_FILE")
	if !ok && tlsEnabled || tlsEnabled && tlsKey == "" {
		log.Fatalf("KEY_FILE env variable not set when TLS is enabled")
	}
	tlsKey = ""
	return &Config{
		DatabaseURL: dbURL,
		RunAddr:     runAddr,
		SecretKey:   secretKey,
		TokenTTL:    time.Duration(ttl) * time.Hour,
		TLSEnabled:  tlsEnabled,
		CertFile:    tlsCert,
		KeyFile:     tlsKey,
	}
}

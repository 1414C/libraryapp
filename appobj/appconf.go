package appobj

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/1414C/lw"
)

// DBConfig type hold pg config info
type DBConfig struct {
	DBDialect           string `json:"db_dialect"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Usr                 string `json:"Usr"`
	Password            string `json:"Password"`
	Name                string `json:"name"`
	ORMLogActive        bool   `json:"ormLogActive"`
	ORMDebugTraceActive bool   `json:"ormDebugTraceActive"`
}

// LeadSetGetConfig holds group leadership KVS config
type LeadSetGetConfig struct {
	StandAlone StandAloneGKVSConfig `json:"local_standalone"`
	Redis      RedisGKVSConfig      `json:"redis"`
	Memcached  MemcachedGKVSConfig  `json:"memcached"`
	Sluggo     SluggoGKVSConfig     `json:"sluggo"`
}

// StandAloneGKVSConfig holds the configuration values used to establish
// the connection to the group-leadership KVS when the application is
// running with a single stand-alone application server.
type StandAloneGKVSConfig struct {
	Active          bool   `json:"active"`
	InternalAddress string `json:"internal_address"`
}

// RedisGKVSConfig holds the configuration values used to establish
// the connection to Redis when it is being used as the group-leadership
// KVS.
type RedisGKVSConfig struct {
	Active        bool   `json:"active"`
	MaxIdle       int    `json:"max_idle"`
	MaxActive     int    `json:"max_active"`
	RedisProtocol string `json:"redis_protocol"`
	RedisAddress  string `json:"redis_address"`
}

// MemcachedGKVSConfig holds the configuration values used to establish
// the connection to Memcached when it is being used as the group-
// leadership KVS.
type MemcachedGKVSConfig struct {
	Active             bool     `json:"active"`
	MemcachedAddresses []string `json:"memcached_addresses"`
}

// SluggoGKVSConfig holds the configuration values used to establish
// the connection to Sluggo when it is being used as the group-
// leadership KVS.
type SluggoGKVSConfig struct {
	Active        bool   `json:"active"`
	SluggoAddress string `json:"sluggo_address"`
}

// LogConfig holds logging config info
type LogConfig struct {
	Active        bool `json:"active"`
	CallLocation  bool `json:"callLocation"`
	ColorMsgTypes bool `json:"colorMsgTypes"`
	InfoMsgs      bool `json:"infoMsgs"`
	WarningMsgs   bool `json:"warningMsgs"`
	ErrorMsgs     bool `json:"errorMsgs"`
	DebugMsgs     bool `json:"debugMsgs"`
	TraceMsgs     bool `json:"traceMsgs"`
}

// ServiceActivation struct
type ServiceActivation struct {
	ServiceName   string `json:"service_name"`
	ServiceActive bool   `json:"service_active"`
}

// ConnectionInfo returns a DBConfig string
func (c DBConfig) ConnectionInfo() string {

	switch c.Dialect() {
	case "postgres":
		if c.Password == "" {
			return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Usr, c.Name)
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Usr, c.Password, c.Name)

	case "mssql":
		// "sqlserver://SA:my_passwd@my_mssql.server.com:1401?database=sqlx")
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", c.Usr, c.Password, c.Host, c.Port, c.Name)

	case "hdb":
		// "hdb://my_user:my_passwd@my_hdb.server.com:30047")
		return fmt.Sprintf("hdb://%s:%s@%s:%d", c.Usr, c.Password, c.Host, c.Port)

	case "sqlite":
		// "sqlite3", "testdb.sqlite"
		return fmt.Sprintf("%s", c.Name)

	case "mysql":
		// "my_user:my_passwd@tcp(my_mysql.server.com:3306)/sqlx?charset=utf8&parseTime=True&loc=Local")
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.Usr, c.Password, c.Host, c.Port, c.Name)

	default:
		panic(fmt.Errorf("dialect %s is not recognized", c.DBDialect))

	}
}

// Dialect returns the db type
func (c DBConfig) Dialect() string {
	return c.DBDialect
}

// DefaultDBConfig provides a default db config
// for use in development testing
func DefaultDBConfig() DBConfig {
	return DBConfig{
		DBDialect:           "postgres",
		Host:                "localhost",
		Port:                5432,
		Usr:                 "my_user",
		Password:            "my_passwd",
		Name:                "glrestgen",
		ORMLogActive:        false,
		ORMDebugTraceActive: false,
	}
}

// DefaultLogConfig provides a default logging
// configuration defaulted to active.
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Active:        true,
		CallLocation:  false,
		ColorMsgTypes: true,
		InfoMsgs:      true,
		WarningMsgs:   true,
		ErrorMsgs:     true,
		DebugMsgs:     false,
		TraceMsgs:     false,
	}
}

// DefaultLeadSetGetConfig returns a valid default configuration
// suitable for running a jiffy-application in stand-alone mode.
func DefaultLeadSetGetConfig() LeadSetGetConfig {
	return LeadSetGetConfig{
		StandAlone: StandAloneGKVSConfig{
			Active:          true,
			InternalAddress: "127.0.0.1:4444",
		},
		Redis: RedisGKVSConfig{
			Active:        false,
			MaxIdle:       80,
			MaxActive:     12000,
			RedisProtocol: "tcp",
			RedisAddress:  "127.0.0.1:6379",
		},
		Memcached: MemcachedGKVSConfig{
			Active:             false,
			MemcachedAddresses: []string{"127.0.0.1:11211"},
		},
		Sluggo: SluggoGKVSConfig{
			Active:        false,
			SluggoAddress: "127.0.0.1:7070",
		},
	}
}

// Config isProd := false
// likely don't need ther private-keys beyond what you use in Login() to
// generate the jwt-token in cases where the access token was not created
// by an IDP.
type Config struct {
	ExternalAddress     string              `json:"external_address"`
	InternalAddress     string              `json:"internal_address"`
	Env                 string              `json:"env"`
	PingCycle           uint                `json:"ping_cycle"`
	FailureThreshold    uint64              `json:"failure_threshold"`
	Pepper              string              `json:"pepper"`
	Database            DBConfig            `json:"database"`
	LeadSetGet          LeadSetGetConfig    `json:"group_leader_kvs"`
	Logging             LogConfig           `json:"logging"`
	CertFile            string              `json:"cert_file"`
	KeyFile             string              `json:"key_file"`
	RSA256PrivKeyFile   string              `json:"rsa256_priv_key_file"`
	RSA256PubKeyFile    string              `json:"rsa256_pub_key_file"`
	RSA384PrivKeyFile   string              `json:"rsa384_priv_key_file"`
	RSA384PubKeyFile    string              `json:"rsa384_pub_key_file"`
	RSA512PrivKeyFile   string              `json:"rsa512_priv_key_file"`
	RSA512PubKeyFile    string              `json:"rsa512_pub_key_file"`
	ECDSA256PrivKeyFile string              `json:"ecdsa256_priv_key_file"`
	ECDSA256PubKeyFile  string              `json:"ecdsa256_pub_key_file"`
	ECDSA384PrivKeyFile string              `json:"ecdsa384_priv_key_file"`
	ECDSA384PubKeyFile  string              `json:"ecdsa384_pub_key_file"`
	ECDSA521PrivKeyFile string              `json:"ecdsa521_priv_key_file"`
	ECDSA521PubKeyFile  string              `json:"ecdsa521_pub_key_file"`
	JWTSignMethod       string              `json:"jwt_sign_method"`
	JWTLifetime         uint                `json:"jwt_lifetime"`
	ServiceActivations  []ServiceActivation `json:"service_activations"`
}

// IsProd informs the app which environment it is running in
func (c Config) IsProd() bool {
	if c.Env == "prod" {
		return true
	}
	return false
}

// IsDev informs the app which environment it is running in
func (c Config) IsDev() bool {
	if c.Env == "dev" {
		return true
	}
	return false
}

// DefaultServiceActivations returns the app's default service activations
func DefaultServiceActivations() []ServiceActivation {
	s := ServiceActivation{}
	sa := []ServiceActivation{}
	s.ServiceName = "Library"
	s.ServiceActive = true
	sa = append(sa, s)

	s.ServiceName = "Book"
	s.ServiceActive = true
	sa = append(sa, s)

	return sa
}

// DefaultConfig returns the app's default config in a Config structure
func DefaultConfig() Config {
	return Config{
		ExternalAddress:     "127.0.0.1:3000",
		InternalAddress:     "127.0.0.1:4444",
		Env:                 "def",
		PingCycle:           1,
		FailureThreshold:    5,
		Pepper:              "secret-pepper-key",
		Database:            DefaultDBConfig(),
		LeadSetGet:          DefaultLeadSetGetConfig(),
		Logging:             DefaultLogConfig(),
		CertFile:            "", // https
		KeyFile:             "", // https
		RSA256PrivKeyFile:   "",
		RSA256PubKeyFile:    "",
		RSA384PrivKeyFile:   "",
		RSA384PubKeyFile:    "",
		RSA512PrivKeyFile:   "",
		RSA512PubKeyFile:    "",
		ECDSA256PrivKeyFile: "",
		ECDSA256PubKeyFile:  "",
		ECDSA384PrivKeyFile: "jwtkeys/ecdsa/ec384.priv.pem",
		ECDSA384PubKeyFile:  "jwtkeys/ecdsa/ec384.priv.pem",
		ECDSA521PrivKeyFile: "",
		ECDSA521PubKeyFile:  "",
		JWTSignMethod:       "ES384",
		JWTLifetime:         120,
		ServiceActivations:  DefaultServiceActivations(),
	}
}

// LoadConfig loads the config file, or falls back to the default
func LoadConfig(configReq RunMode) Config {

	var fName string
	switch configReq {
	case cDev:
		fName = ".dev.config.json"
	case cPrd:
		fName = ".prd.config.json"
	default:
		fName = "default_config" // not a file ;)
	}

	f, err := os.Open(fName) // cDef will always fail here - ok
	if err != nil {
		if configReq != cDef {
			panic(err)
		}
		lw.Console("using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	lw.Console("successfully loaded the config file...")
	// log.Println("config:", c)
	return c
}

package appobj

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/1414C/libraryapp/group/gmcom"
	"github.com/1414C/lw"
	"github.com/1414C/sluggo/wscl"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/garyburd/redigo/redis"
)

// LeadSetGet provides a sample implementation of the gmcom.GMLeaderSetterGetter interface in order
// to support persistence of the current group-leader information.

// StandAloneLeadSetGet is an implementation of the gmcom.GMLeaderSetterGetter interface
// in order to facilitate the running of a standalone jiffy application-server.
// In such a scenario, there is no need for a leader per se, but the codebase
// demands an implementation of the interface.
type StandAloneLeadSetGet struct {
	gmcom.GMLeaderSetterGetter
	LocalLeaderIPAddress string
}

// Cleanup is needed to satisfy interface gmcom.GMLeaderSetterGetter but is
// not really needed to tidyup an open pool/connection.
func (sa *StandAloneLeadSetGet) Cleanup() error {
	// nothing to do for stand-alone
	return nil
}

// GetDBLeader retrieves the current leader information from the persistence layer.
func (sa *StandAloneLeadSetGet) GetDBLeader() (*gmcom.GMLeader, error) {

	// access the database here to read the current leader
	l := &gmcom.GMLeader{
		LeaderID:        1,
		LeaderIPAddress: sa.LocalLeaderIPAddress,
	}
	return l, nil
}

// SetDBLeader stores the current leader information in the persistence layer.
func (sa *StandAloneLeadSetGet) SetDBLeader(l gmcom.GMLeader) error {

	// running as local/standalone so the leader info doesn't
	// really matter.  the stand-alone group leadership info
	// is set by the call it InitializeStandAloneLeadSetGet()
	// and is deemed to be immutable.
	return nil
}

// For testing with sluggo, go get -u github.com/1414C/sluggo
//
// Execute sluggo from the command-line as follows:
// go run main.go -a <ipaddress:port>
// For example:
// $ go run main.go -a 192.168.1.40:5050

// SluggoLeadSetGet is a struct implementing the gmcom.GMLeaderSetterGetter
// interface in order to support access to the sluggo KVS.
type SluggoLeadSetGet struct {
	gmcom.GMLeaderSetterGetter
	internalAddress string
}

// CreateSluggoLeadSetGet creates an instance of interface gmcom.GMLeaderSetterGetter
// targeting a running sluggo KVS process.
func (sg *SluggoLeadSetGet) CreateSluggoLeadSetGet(c SluggoGKVSConfig) gmcom.GMLeaderSetterGetter {

	return &SluggoLeadSetGet{
		internalAddress: c.SluggoAddress,
	}
}

// Cleanup is needed to satisfy interface gmcom.GMLeaderSetterGetter but is
// not really needed to tidyup an open pool/connection.
func (sg *SluggoLeadSetGet) Cleanup() error {
	// nothing to do for sluggo
	return nil
}

// GetDBLeader retrieves the current leader information from the persistence layer.
func (sg *SluggoLeadSetGet) GetDBLeader() (*gmcom.GMLeader, error) {

	// access the database here to read the current leader
	l := &gmcom.GMLeader{}
	// wscl.GetCacheEntry("LEADER", l, "192.168.112.192:7070")
	wscl.GetCacheEntry("LEADER", l, sg.internalAddress)
	return l, nil
}

// SetDBLeader stores the current leader information in the persistence layer.
func (sg *SluggoLeadSetGet) SetDBLeader(l gmcom.GMLeader) error {

	// access the database here to set a new current leader
	// wscl.AddUpdCacheEntry("LEADER", &l, "192.168.112.192:7070")
	wscl.AddUpdCacheEntry("LEADER", &l, sg.internalAddress)
	return nil
}

// RedisLeadSetGet is a struct implementing the gmcom.GMLeaderSetterGetter
// interface in order to support access to the redis KVS.
type RedisLeadSetGet struct {
	gmcom.GMLeaderSetterGetter
	pool *redis.Pool
	conn redis.Conn
}

// InitializeRedisLeadSetGet creates an instance of interface gmcom.GMLeaderSetterGetter
// targeting a running sluggo KVS process.
func (rg *RedisLeadSetGet) InitializeRedisLeadSetGet(c RedisGKVSConfig) error {

	rg.pool = rg.newPool(c)
	rg.conn = rg.pool.Get()
	return nil
}

// Cleanup - its hard to block in here, so we rely on an external call to tear
// the redis connection down on SIGINT/KILL etc.  fugly.
func (rg *RedisLeadSetGet) Cleanup() error {

	lw.Console("Redis KVS: RedisLeadSetGet running Cleanup()")
	err := rg.conn.Flush()
	if err != nil {
		return err
	}
	err = rg.conn.Close()
	if err != nil {
		return err
	}
	err = rg.pool.Close()
	if err != nil {
		return err
	}
	return nil
}

// newPool attempts to create a new connection pool and attach to
// redis using the config information.  failure to connect results
// in a panic.
func (rg *RedisLeadSetGet) newPool(c RedisGKVSConfig) *redis.Pool {

	return &redis.Pool{
		MaxIdle:   c.MaxIdle,
		MaxActive: c.MaxActive,
		Dial: func() (redis.Conn, error) {
			// c, err := redis.Dial("tcp", ":6379")
			c, err := redis.Dial(c.RedisProtocol, c.RedisAddress)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// GetDBLeader retrieves the current leader information from the persistence layer.
func (rg *RedisLeadSetGet) GetDBLeader() (*gmcom.GMLeader, error) {

	// access the database here to read the current leader
	l := gmcom.GMLeader{}
	// wscl.GetCacheEntry("LEADER", l, "192.168.112.192:7070")

	// read the gob-encoded leader information from redis
	result, err := rg.conn.Do("GET", "LEADER")
	if err != nil {
		return nil, err
	}

	if result != nil {
		raw := result.([]byte)
		decBuf := bytes.NewBuffer(raw)
		err = gob.NewDecoder(decBuf).Decode(&l)
		if err != nil {
			return nil, err
		}
		lw.Console("Redis KVS: RedisLeadSetGet.GetDBLeader got: %v", l)
		return &l, nil
	}
	return &l, nil
}

// SetDBLeader stores the current leader information in the persistence layer.
func (rg *RedisLeadSetGet) SetDBLeader(l gmcom.GMLeader) error {

	// access redis to set the new leader
	// wscl.AddUpdCacheEntry("LEADER", &l, "192.168.112.192:7070")
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(l)
	if err != nil {
		lw.Console("Redis KVS: SetDBLeader failed to gob-encode the new leader record: %s", err)
	}

	// write the gob-encoded leader information to redis
	_, err = rg.conn.Do("SET", "LEADER", encBuf)
	if err != nil {
		lw.Console("Redis KVS: SetDBLeader failed to update redis with the new leader record: %s", err)
	}
	return nil
}

// MemcachedLeadSetGet is a struct implementing the gmcom.GMLeaderSetterGetter
// interface in order to support access to the memcached KVS.
type MemcachedLeadSetGet struct {
	gmcom.GMLeaderSetterGetter
	client *memcache.Client
}

// InitializeMemcachedLeadSetGet creates an instance of interface gmcom.GMLeaderSetterGetter
// targeting a running memcached KVS system/cluster.
func (mc *MemcachedLeadSetGet) InitializeMemcachedLeadSetGet(c MemcachedGKVSConfig) error {

	// mc.client = memcache.New("192.168.112.50:11211")
	mc.client = memcache.New(c.MemcachedAddresses...)
	if mc.client == nil {
		lw.Console("Memcached KVS: InitializeMemcachedLeadSetGet failed initialize a memcached client.")
		return errors.New("memcached KVS: InitializeMemcachedLeadSetGet failed initialize a memcached client")
	}
	return nil
}

// Cleanup - its hard to block in here, so we rely on an external call to tear
// the memcached connection down on SIGINT/KILL etc.  fugly.
func (mc *MemcachedLeadSetGet) Cleanup() error {

	lw.Console("Memcached KVS: MemcachedLeadSetGet running Cleanup()")
	// there is no apparent cleanup function in the
	// memcached client - set the connection to nil?
	mc.client = nil
	return nil
}

// GetDBLeader retrieves the current leader information from the persistence layer.
func (mc *MemcachedLeadSetGet) GetDBLeader() (*gmcom.GMLeader, error) {

	// access the database here to read the current leader
	l := gmcom.GMLeader{}

	// read the gob-encoded leader information from memcached
	it, err := mc.client.Get("LEADER")
	if err != nil {
		return nil, err
	}

	if it.Value != nil {
		decBuf := bytes.NewBuffer(it.Value)
		err = gob.NewDecoder(decBuf).Decode(&l)
		if err != nil {
			return nil, err
		}
		lw.Console("Memcached KVS: MemcachedLeadSetGet.GetDBLeader got: %v", l)
		return &l, nil
	}
	return &l, nil
}

// SetDBLeader stores the current leader information in the persistence layer.
func (mc *MemcachedLeadSetGet) SetDBLeader(l gmcom.GMLeader) error {

	// access memcached to set the new leader
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(l)
	if err != nil {
		lw.Console("Memcached KVS: SetDBLeader failed to gob-encode the new leader record: %s", err)
	}

	// write the gob-encoded leader information to memcached
	// _, err = rg.conn.Do("SET", "LEADER", encBuf)
	err = mc.client.Set(&memcache.Item{Key: "LEADER", Value: encBuf.Bytes()})
	if err != nil {
		lw.Console("Memcached KVS: SetDBLeader failed to update the new leader record: %s", err)
	}
	return nil
}

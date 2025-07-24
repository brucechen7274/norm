package norm

import (
	"fmt"
	"github.com/haysons/norm/logger"
	"github.com/haysons/norm/resolver"
	"github.com/haysons/norm/statement"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"net"
	"strconv"
	"time"
)

// DB uses statement.Statement to construct nGQL statements,
// and then executes them through nebula.SessionPool.
// You can retrieve the results via methods such as Find, Exec, or Pluck.
// DB is concurrency-safe: multiple statements can be executed concurrently using a single DB instance.
//
// Note:
//   - nebula-graph officially provides a SessionPool, so there's no need to manage a custom connection pool.
//     In most cases, a single DB instance is sufficient for the application.
//   - statement.Statement is NOT concurrency-safe.
//     Do not build nGQL statements concurrently using the same Statement instance.
//   - Embedded fields in struct definitions are NOT supported.
//     Avoid using embedded fields when defining vertex/edge structs.
type DB struct {
	Statement   *statement.Statement
	conf        *Config
	sessionPool *nebula.SessionPool
	clone       int
}

// Open creates a new DB instance.
//
// It initializes configuration options, resolves timezone, sets logger,
// parses the server address, and creates the session pool.
// The returned DB instance is ready to execute nGQL statements.
func Open(conf *Config, opts ...ConfigOption) (*DB, error) {
	for _, o := range opts {
		o.apply(conf)
	}

	if conf.TimezoneName != "" {
		loc, err := time.LoadLocation(conf.TimezoneName)
		if err != nil {
			return nil, fmt.Errorf("norm: load timezone failed: %v", err)
		}
		conf.timezone = loc
	} else {
		conf.timezone = time.Local
	}
	resolver.SetTimezone(conf.timezone)

	if conf.logger == nil {
		conf.logger = logger.Default
	}

	hostAddr, err := parseServerAddr(conf.Addresses)
	if err != nil {
		return nil, err
	}
	poolConf, err := nebula.NewSessionPoolConf(conf.Username, conf.Password, hostAddr, conf.SpaceName, parseSessionOptions(conf)...)
	if err != nil {
		return nil, fmt.Errorf("norm: build session pool conf failed: %v", err)
	}
	pool, err := nebula.NewSessionPool(*poolConf, nebula.DefaultLogger{})
	if err != nil {
		return nil, fmt.Errorf("norm: create session pool failed: %v", err)
	}

	db := &DB{
		Statement:   statement.New(),
		conf:        conf,
		sessionPool: pool,
		clone:       1, // when clone is 1, the Statement object will be copied to ensure that the same singleton build statement does not affect each other.
	}
	return db, nil
}

func parseServerAddr(addrList []string) ([]nebula.HostAddress, error) {
	hostAddr := make([]nebula.HostAddress, 0, len(addrList))
	for _, addr := range addrList {
		host, portTmp, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("norm: parse server addr failed: %w", err)
		}
		port, err := strconv.Atoi(portTmp)
		if err != nil {
			return nil, fmt.Errorf("norm: convert server addr port failed: %w", err)
		}
		hostAddr = append(hostAddr, nebula.HostAddress{
			Host: host,
			Port: port,
		})
	}
	return hostAddr, nil
}

func parseSessionOptions(conf *Config) []nebula.SessionPoolConfOption {
	poolOptions := make([]nebula.SessionPoolConfOption, 0)
	if conf.MaxOpenConns > 0 {
		poolOptions = append(poolOptions, nebula.WithMaxSize(conf.MaxOpenConns))
	}
	if conf.MinOpenConns > 0 {
		poolOptions = append(poolOptions, nebula.WithMinSize(conf.MinOpenConns))
	}
	if conf.ConnTimeout > 0 {
		poolOptions = append(poolOptions, nebula.WithTimeOut(conf.ConnTimeout))
	}
	if conf.ConnMaxIdleTime > 0 {
		poolOptions = append(poolOptions, nebula.WithIdleTime(conf.ConnMaxIdleTime))
	}
	poolOptions = append(poolOptions, conf.nebulaSessionOpts...)
	return poolOptions
}

func (db *DB) getInstance() *DB {
	if db.clone > 0 {
		tx := &DB{conf: db.conf, sessionPool: db.sessionPool, clone: 0}
		tx.Statement = statement.New()
		return tx
	}
	return db
}

func (db *DB) Close() error {
	db.sessionPool.Close()
	return nil
}

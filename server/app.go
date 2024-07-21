package server

import (
	"fmt"
	"nrs16/cme/config"
	"nrs16/cme/repository"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	redis "github.com/redis/go-redis/v9"
)

type App struct {
	Router                 *mux.Router
	AuthenticationDatabase repository.AuthenticationDatabase
	ChatDatabase           repository.ChatsDatabase
	Redis                  *redis.Client
}

func InitialiseApp(conf config.Config) (*App, error) {

	app := new(App)

	// configure router
	router := mux.NewRouter()
	app.Router = router
	ConfigureRoutes(app)

	//connect to database

	cluster := gocql.NewCluster(conf.Database.Host)
	cluster.Keyspace = conf.Database.KeySpace
	cluster.Consistency = gocql.Quorum

	cqlSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	app.AuthenticationDatabase = repository.AuthenticationDatabase{Database: cqlSession}
	app.ChatDatabase = repository.ChatsDatabase{Database: cqlSession}
	app.Redis = StartRedisClient(conf.Redis)
	return app, nil
}

func StartRedisClient(r config.RedisConfiguration) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", r.Host, r.Port),
	})
}

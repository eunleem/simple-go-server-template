package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	redistore "github.com/boj/redistore"
	"github.com/go-redis/redis"

	"gopkg.in/mgo.v2"

	mgoutil "gitlab.com/eunleem/gopack/mgoutil-v1"
	redisutil "gitlab.com/eunleem/gopack/redisutil-v1"
)

type SharedData struct {
	MongoConnection      *mgo.Session
	RedisSessionStore    *redistore.RediStore
	RedisCacheConnection *redis.Client
}

var sharedData SharedData

func closeConnections() {
	if sharedData.MongoConnection != nil {
		mgoutil.CloseSession(sharedData.MongoConnection)
	}

	if sharedData.RedisSessionStore != nil {
		sharedData.RedisSessionStore.Close()
	}

	if sharedData.RedisCacheConnection != nil {
		sharedData.RedisCacheConnection.Close()
	}

	log.Print("All the connections have been closed!")
}

func openMongoDb() {
	// If MongoDB is being used
	// Open Mongo Connection
	if dbSession, err := mgoutil.OpenSession(conf.MongoDb.Host, conf.MongoDb.Username, conf.MongoDb.Password); err != nil {
		panic(err)
	} else {
		sharedData.MongoConnection = dbSession
	}
}

func openRedisSession() {
	if store, err := redistore.NewRediStore(10, "tcp", conf.Redis.Host, conf.Redis.Password, []byte("So.merAn.D-oM$Tr|ng")); err != nil {
		panic(err)
	} else {
		sharedData.RedisSessionStore = store
	}
}

func openRedisCache() {
	if redisClient, err := redisutil.New(conf.Redis.Host, conf.Redis.Password, 10); err != nil {
		panic(err)
	} else {
		sharedData.RedisCacheConnection = redisClient
	}
}

// config.go's init
func main() {
	log.Print("Server Launched")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if conf.MongoDb.Host != "" {
		openMongoDb()
	}

	if conf.Redis.Host != "" {
		openRedisSession()
		openRedisCache()
	}

	// Register Routes
	http.HandleFunc("/", indexHandler)

	// Run Server
	portStr := strconv.Itoa(conf.WebServer.Port)
	log.Print("Listening on port " + portStr)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Fatal("ListenAndServe on "+portStr, err)
	}

	defer closeConnections()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// filePath := filepath.Join(conf.WebServer.WebDir, "index.html")
	// log.Print(filePath)
	// http.ServeFile(w, r, filePath)
	http.ServeFile(w, r, filepath.Join(conf.WebServer.WebDir, "index.html"))
}

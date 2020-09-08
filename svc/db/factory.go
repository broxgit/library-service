package db

import "os"

//GetDb is a factory method that will retrieve the user specified Database service
func GetDb() (Database, error) {
	keySpace, found := os.LookupEnv("CASSANDRA_KEYSPACE")

	if found {
		return newCassandraDb(keySpace)
	} else {
		return newCacheDatabase()
	}
}

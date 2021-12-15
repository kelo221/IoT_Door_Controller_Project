package main

import (
	"context"
	_ "context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var db driver.Database
var userCol driver.Collection

func handleDatabase() {

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
		TLSConfig: &tls.Config{ /*...*/ },
	})
	if err != nil {
		fmt.Println(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "1234"),
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create a database
	db, err = client.Database(nil, "DOOR_DATA")
	if err != nil {
		fmt.Println(err, "creating new...")
		ctx := context.Background()
		db, err = client.CreateDatabase(ctx, "DOOR_DATA", nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Create a collection for users
	userCol, err = db.Collection(nil, "DOOR_LOGIN")
	if err != nil {
		fmt.Println(err, "creating new...")
		ctx := context.Background()
		options := &driver.CreateCollectionOptions{ /* ... */ }
		userCol, err = db.CreateCollection(ctx, "DOOR_LOGIN", options)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Create a collection for RFID
	userCol, err = db.Collection(nil, "DOOR_RFID")
	if err != nil {
		fmt.Println(err, "creating new...")
		ctx := context.Background()
		options := &driver.CreateCollectionOptions{ /* ... */ }
		userCol, err = db.CreateCollection(ctx, "DOOR_RFID", options)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Create a collection for history
	userCol, err = db.Collection(nil, "DOOR_HISTORY")
	if err != nil {
		fmt.Println(err, "creating new...")
		ctx := context.Background()
		options := &driver.CreateCollectionOptions{ /* ... */ }
		userCol, err = db.CreateCollection(ctx, "DOOR_HISTORY", options)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Create a collection for lock changes
	userCol, err = db.Collection(nil, "LOCK_HISTORY")
	if err != nil {
		fmt.Println(err, "creating new...")
		ctx := context.Background()
		options := &driver.CreateCollectionOptions{ /* ... */ }
		userCol, err = db.CreateCollection(ctx, "LOCK_HISTORY", options)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func createAccounts() {
	salt := []byte("salt")

	aqlNoReturn("UPSERT { username: 'User' } " +
		"INSERT { username: 'User', name: 'Marko', hash: '" + hex.EncodeToString(HashPassword([]byte("User"), salt)) + "', dateCreated: DATE_NOW() } " +
		"UPDATE {} IN DOOR_LOGIN")

	aqlNoReturn("UPSERT { username: 'User2' } " +
		"INSERT { username: 'User2', name: 'Marie', hash: '" + hex.EncodeToString(HashPassword([]byte("User2"), salt)) + "', dateCreated: DATE_NOW() } " +
		"UPDATE {} IN DOOR_LOGIN")

}

func createRDIF() {
	salt := []byte("salt")

	aqlNoReturn("UPSERT { RFID_OWNER: 'User' } " +
		"INSERT { RFID_OWNER: 'User', HASHED_RFID: '" + hex.EncodeToString(HashPassword([]byte("F7 20 05 3A"), salt)) + "', dateCreated: DATE_NOW() } " +
		"UPDATE {} IN DOOR_RFID")

}

func aqlNoReturn(query string) {

	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func(cursor driver.Cursor) {
		err3 := cursor.Close()
		if err3 != nil {
			fmt.Println(err3)
		}
	}(cursor)

}

func aqlToString(query string) string {

	var result string

	ctx := context.Background()
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		// handle error
	}
	defer func(cursor driver.Cursor) {
		err3 := cursor.Close()
		if err3 != nil {
			fmt.Println(err3)
		}
	}(cursor)
	for {
		_, err2 := cursor.ReadDocument(ctx, &result)
		if driver.IsNoMoreDocuments(err2) {
			break
		} else if err2 != nil {
			fmt.Println(err2)
		}
	}

	return result
}

type logData struct {
	Mode string `json:"mode,omitempty"`
	Name string `json:"name,omitempty"`
	Time int    `json:"time,omitempty"`
}

func aqlJSON(query string) []logData {

	var dataPayload []logData

	ctx := context.Background()
	//query = "FOR Speed IN IOT_DATA_SENSOR RETURN Speed"
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		// handle error
	}
	defer func(cursor driver.Cursor) {
		err3 := cursor.Close()
		if err3 != nil {
			fmt.Println(err3)
		}
	}(cursor)
	for {
		var doc logData
		_, err2 := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err2) {
			break
		} else if err2 != nil {
			fmt.Println(err2)
		}
		//fmt.Printf("Got doc with key '%s' from query\n", meta.Rev)
		//fmt.Println(doc)
		dataPayload = append(dataPayload, doc)
	}

	return dataPayload

}

package apis

import (
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/firebase"
	"github.com/keyurboss/Generic-Servers/go/lint-shortner/utility"
)

var firebaseDb = firebase.GetFirebaseFirestoreDatabase()

func init() {
	// firebase.GetFirebaseFirestoreDatabase()
	currentIdConfig, err := firebaseDb.GetPublicData("server_config", "current_uni_id")
	if err != nil {
		panic(err)
	}
	if id, ok := currentIdConfig["id"]; ok {
		if id, ok := id.(string); ok {
			utility.UniqueIdConst.SetUniqueId(id)
		}
	}
	if c, ok := currentIdConfig["changeNoOfDigits"]; ok {
		if c, ok := c.(int); ok {
			utility.UniqueIdConst.ChangeNoOfDigits = c
		}
	}
	if c, ok := currentIdConfig["increseDigitBy"]; ok {
		if c, ok := c.(int); ok {
			utility.UniqueIdConst.IncreseDigitBy = c
		}
	}
	serverConfig, err := firebaseDb.GetPublicData("server_config", "server_url")
	if err != nil {
		panic(err)
	}
	if url, ok := serverConfig["server_url"]; ok {
		if url, ok := url.(string); ok {
			utility.ServerUrl = url
		}
	}
	// println("Firebase Database Service Initialized")
}

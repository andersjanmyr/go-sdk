// to run the CPU profiling: go build -ldflags "-X main.ProfileMode=mem" main_profile_feature.go && ./main_profile_feature

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/optimizely/go-sdk/optimizely/client"
	"github.com/optimizely/go-sdk/optimizely/decision"
	"github.com/optimizely/go-sdk/optimizely/entities"
	"github.com/optimizely/go-sdk/optimizely/notification"

	"github.com/pkg/profile"
)

func stressTest() {
	/*
		For the test app, the biggest json file is used with 100 entities.
		DATAFILES_DIR has to be set to point to the path where 100_entities.json is located.
	*/

	var datafileDir = path.Join(os.Getenv("DATAFILES_DIR"), "100_entities.json")

	datafile, err := ioutil.ReadFile(datafileDir)
	if err != nil {
		log.Fatal(err)
	}

	optlyClient := &client.OptimizelyFactory{
		Datafile: datafile,
	}

	user := entities.UserContext{
		ID: "test_user_1",
		Attributes: map[string]interface{}{
			"attr_5": "testvalue",
		},
	}

	// Creates a default, canceleable context
	notificationCenter := notification.NewNotificationCenter()
	decisionService := decision.NewCompositeService(notificationCenter)

	clientOptions := client.Options{
		DecisionService: decisionService,
	}
	clientApp, err := optlyClient.ClientWithOptions(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	clientApp.IsFeatureEnabled("feature_5", user)
}

var ProfileMode = ""

const RUN_NUMBER = 50

func main() {

	if ProfileMode != "" {

		switch ProfileMode {
		case "mem":
			defer profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.MemProfileRate(1)).Stop()
		case "cpu":
			defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
		}

		for i := 0; i < RUN_NUMBER; i++ {
			stressTest()
		}
	}

}

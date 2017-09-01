package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tebeka/selenium"
)

func TestWithSelenium(t *testing.T) {
	const (
		seleniumPath = "vendor/selenium-server-standalone-3.5.3.jar"
		port         = 8080
	)

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.Output(os.Stderr),
	}

	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "safari"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

}
func TestOver20Connections(t *testing.T) {
	envvars = getEnv()
	db := dbmanager{url: envvars["postgres_url"]}
	db.OpenConnection()

	for j := 0; j < 100; j++ {
		out := make(chan int)
		index := j
		go func() {
			t.Logf("%d", index)
			db.db.Query("select * from Responses")
			db.GetRandomResponse(&Response{wave: 2}, false)
			out <- 0
			close(out)
		}()
	}

	success, err := db.CheckConnection()
	if success || err == nil {
		t.Fatalf("Postgres overflow connection did not ping incorrectly")
	}

	t.Logf("Postgres connection error: %s. Success?", err.Error())
}

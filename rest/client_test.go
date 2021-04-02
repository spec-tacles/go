package rest_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spec-tacles/go/rest"
	"github.com/spec-tacles/go/types"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestDo(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/@me", nil)
	if err != nil {
		t.Error(err)
		return
	}

	client := rest.NewClient(os.Getenv("DISCORD_TOKEN"), "8")
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(res)
	if res.StatusCode >= 400 {
		t.Fail()
	}
}

func TestGateway(t *testing.T) {
	var info types.GatewayBot
	client := rest.NewClient(os.Getenv("DISCORD_TOKEN"), "8")
	err := client.DoJSON(http.MethodGet, "/gateway/bot", nil, &info)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", info)
}

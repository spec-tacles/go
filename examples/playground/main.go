package main

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func main() {
	var a struct {
		IDs []snowflake.ID `json:"ids"`
		ID  snowflake.ID   `json:"id"`
	}
	if err := json.Unmarshal([]byte(`
	{
		"ids": ["123", "456", "789"],
		"id": "0123"
	}
	`), &a); err != nil {
		panic(err)
	}

	fmt.Println(a.IDs)
	fmt.Println(a.ID)
}

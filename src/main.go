package main

import (
	"token_price/src/actions"
	"token_price/src/config"
	"token_price/src/utils"
)

func main() {

	utils.LoadEnv()
	//connect to mongodb
	mongoURI := utils.GetValue("MONGO_URI", "")
	utils.ConnectMongoDB(mongoURI)

	//getPrice
	tokens := config.TOKENS
	for i := 0; i < len(tokens); i++ {
		actions.CheckPrice(tokens[i])
	}

}

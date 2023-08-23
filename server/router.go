package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

/*
/: ping (POST)
/newGame : new game (POST)
/newUser : new user (POST)
/login   : login as user (POST)
*/
func CreateRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/v1/", ping)
	router.GET("/api/v1/fen/:fen", fen)
	router.GET("/api/v1/test", testDB)
	router.POST("/api/v1/create", newGame)
	router.GET("/api/v1/game/:id", getGame)
	return router
}

func ping(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Info("Server pinged.")
	fmt.Fprintf(w, "pong")
}

func fen(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fenIDstring := ps.ByName("fen")

	fenID, err := strconv.Atoi(fenIDstring)
	if err != nil {
		log.WithField("fen", fenIDstring).Error("Unable to convert fen to int")
	}

	log.WithField("fenID", fenID).Info("received request")

	var game Game
	res := db.First(&game, fenID)

	if res.Error != nil {
		log.Errorln(res.Error)
		log.Info("Sending bad request response")
		http.Error(w, "Inavlid FEN Number", http.StatusBadRequest)
	} else {

		resp := make(map[string]string)
		resp["fen"] = game.Fen

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Error(err)
		}
		w.Write(jsonResp)
	}
}

func testDB(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Info("Test db.")

	games := []Game{}
	db.Find(&games)
	fmt.Println(games)
}

func getGame(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	gameIDstring := ps.ByName("id")

	gameID, err := strconv.Atoi(gameIDstring)
	if err != nil {
		log.WithField("id", gameIDstring).Error("Unable to convert gameID to int")
	}
	var game = &Game{ID: gameID}
	res := db.First(game)

	if res.Error != nil {
		log.Errorln(res.Error)
		log.Info("Sending bad request response")
		http.Error(w, "Inavlid GameID Number", http.StatusBadRequest)
	} else {
		type Response struct {
			ID      int
			Fen     string
			Done    bool
			PgnPath string
		}

		resp := Response{
			ID:      game.ID,
			Fen:     game.Fen,
			Done:    game.Done,
			PgnPath: game.PgnPath,
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Error(err)
		}
		w.Write(jsonResp)
	}

}

func newGame(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	game, err := insertNewGame()
	if err != nil {
		log.Errorf("Error creating game: %v", err)
	}

	type gameResponse struct {
		ID    int
		Token string
	}

	response := gameResponse{
		ID: game.ID, Token: game.WhiteToken,
	}
	//response := make(map[string]interface{})
	//response["id"] = game.ID
	//response["token"] = game.WhiteToken
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(jsonResp)
}
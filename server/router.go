package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/connoraubry/chessbot-go/engine"
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
	router.GET("/api/v1/game/:id", getGame)

	router.POST("/api/v1/create", newGame)
	router.POST("/api/v1/move", move)
	router.POST("/api/v1/join", joinGame)
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
func joinGame(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	log.Info("Join game requested")
	type joinReq struct {
		ID int
	}

	var request joinReq
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Errorf("Error parsing req body into joinReq")
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	var game = &Game{ID: request.ID}
	res := db.First(game)

	if res.Error != nil {
		log.Errorln(res.Error)
		log.Info("Sending bad request response")
		http.Error(w, "Inavlid GameID Number", http.StatusBadRequest)
		return
	}

	if game.BlackJoined {
		log.Error("Black player already joined game")
		http.Error(w, "Another player has already joined the game", http.StatusBadRequest)
		return
	}

	log.WithField("id", request.ID).Info("Joining game")

	type gameResponse struct {
		ID    int
		Token string
	}

	response := gameResponse{
		ID: game.ID, Token: game.BlackToken,
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}

	game.BlackJoined = true
	res = db.Save(&game)
	if res.Error != nil {
		log.Errorln(res.Error)
		http.Error(w, "Error updating database, try again.", http.StatusBadRequest)
		return
	}

	w.Write(jsonResp)
}

func move(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	type MoveReq struct {
		Move  string
		ID    int
		Token string
	}
	var m MoveReq
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		log.Errorf("Error parsing req body into MoveReq")
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	log.Info(m)
	var game = &Game{ID: m.ID}
	res := db.First(game)

	if res.Error != nil {
		log.Errorln(res.Error)
		log.Info("Sending bad request response")
		http.Error(w, "Inavlid GameID Number", http.StatusBadRequest)
		return
	}

	e := engine.NewEngine(engine.OptFenString(game.Fen))

	playerTurn := e.CurrentGamestate().Player
	if playerTurn == engine.WHITE {
		if m.Token == game.WhiteToken {
			log.Info("Provided token successfully matches white")
		} else {
			log.Error("Move is black, token is white")
			http.Error(w, "Not user's turn", http.StatusBadRequest)
			return
		}
	} else {
		if m.Token == game.BlackToken {
			log.Info("Provided token successfully matches black")
		} else {
			log.Error("Move is white, token is black")
			http.Error(w, "Not user's turn", http.StatusBadRequest)
			return
		}
	}

	log.Infof("Taking move %v", m.Move)

	moves := e.GetStringToMoveMap(e.GetValidMoves())

	move, ok := moves[m.Move]
	if !ok {
		log.Errorf("Move %v not found in current moves", m.Move)
		http.Error(w, "Error with move", http.StatusBadRequest)
		return
	}

	if !e.TakeMove(move) {
		log.Errorf("Error taking move %v %v", m.Move, move)
		http.Error(w, "Error with move", http.StatusBadRequest)
		return
	}
	newFen := e.ExportToFEN()
	log.WithField("fen", newFen).Info("Updating FEN")

	game.Fen = newFen
	res = db.Save(&game)
	if res.Error != nil {
		log.Errorln(res.Error)
		http.Error(w, "Error updating database, try again.", http.StatusBadRequest)
		return
	}

}

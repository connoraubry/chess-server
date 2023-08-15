package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	port string
}

func New() *Server {
	s := &Server{}
	db = NewDB()
	return s
}

func ping(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Info("Server pinged.")
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
func newGame(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pgnPath := "/test.pgn"
	game := &Game{Fen: fen, Done: false, PgnPath: pgnPath}

	result := db.Create(game)

	if result.Error != nil {
		log.Fatalln(result.Error)
	} else {
		log.WithFields(log.Fields{"id": game.ID}).Info("Created new game")
	}
	response := make(map[string]int)
	response["id"] = game.ID
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(jsonResp)
}

func (s *Server) Run() {
	log.Info("Booting up server")
	router := httprouter.New()
	router.GET("/api/v1/", ping)
	router.GET("/api/v1/fen/:fen", fen)
	router.GET("/api/v1/test", testDB)
	router.GET("/api/v1/create", newGame)
	//testCreate()
	log.WithFields(log.Fields{"port": "3030"}).Info("Running server.")
	log.Fatal(http.ListenAndServe(":3030", router))
}

//func testCreate() {
//	db.Create(&Game{Fen: "8/8/8/8", Done: true, PgnPath: "/lltest.pgn"})
//}

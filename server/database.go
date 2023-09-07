package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
	type User struct {
		gorm.Model
		ID    int
		Name  string
		Token string
	}
*/
type Game struct {
	gorm.Model
	ID   int
	Fen  string
	Done bool

	MoveList string

	WhiteToken string
	BlackToken string

	BlackJoined bool

	IsDev bool
}

func (g *Game) GetLastMove() string {
	var lastMove string
	if g.MoveList != "" {
		moves := strings.Fields(g.MoveList)
		return moves[len(moves)-1]
	}
	return lastMove
}
func (g *Game) AddMove(move string) {
	g.MoveList = fmt.Sprintf("%s %s", g.MoveList, move)
}

var db *gorm.DB

func NewDB() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		"server")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Game{})
	//db.AutoMigrate(&User{})
	return db
}

func genToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func insertNewGame(isDev bool) (*Game, error) {

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	game := &Game{
		Fen:         fen,
		Done:        false,
		MoveList:    "",
		WhiteToken:  genToken(20),
		BlackToken:  genToken(20),
		BlackJoined: false,
		IsDev:       isDev,
	}

	result := db.Create(game)

	if result.Error != nil {
		log.Error(result.Error)
	} else {
		log.WithFields(log.Fields{"id": game.ID, "isDev": isDev}).Info("Created new game")
	}

	return game, result.Error
}

/*
func InsertNewUser(username string) (string, error) {

	_, err := GetUserByUsername(username)
	if err == nil {
		return "", fmt.Errorf("Error: found existing user in database")
	}

	token := genToken(20)
	user := &User{Name: username, Token: token}
	result := db.Create(user)
	return token, result.Error
}

func GetUserByUsername(username string) (User, error) {
	var user User
	result := db.Where("name = ?", username).First(&user)
	return user, result.Error
}
*/

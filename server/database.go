package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Game struct {
	gorm.Model
	ID      int
	Fen     string
	Done    bool
	PgnPath string
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
	//db.AutoMigrate(&Game{})

	return db
}

func GetGameByID(db *sql.DB, id string) ([]Game, error) {
	var games []Game
	rows, err := db.Query("SELECT * FROM game WHERE id = ?", id)

	if err != nil {
		return nil, fmt.Errorf("db.GetGameByID(%v): %v", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.Fen, &g.Done, &g.PgnPath); err != nil {
			return nil, fmt.Errorf("getAllGames: %v", err)
		}
		games = append(games, g)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllGames: %v", err)
	}

	return games, nil
}

func getAllGames(db *sql.DB) ([]Game, error) {
	var games []Game

	rows, err := db.Query("SELECT * FROM game")

	if err != nil {
		return nil, fmt.Errorf("getAllGames: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.Fen, &g.Done, &g.PgnPath); err != nil {
			return nil, fmt.Errorf("getAllGames: %v", err)
		}
		games = append(games, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllGames: %v", err)
	}

	return games, nil

}

func insertNewGame(db *sql.DB) error {

	v := "INSERT INTO game (fen, done, pgnpath) (? ? ?)"
	_, err := db.Exec(v, "7/6/5/4", false, "/path/to/pgn")
	return err
}

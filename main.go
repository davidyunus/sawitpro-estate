package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/davidyunus/sawitpro-estate/src/common"
	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/davidyunus/sawitpro-estate/src/helper"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	estatehttp "github.com/davidyunus/sawitpro-estate/src/estate/delivery/http"
	estatesql "github.com/davidyunus/sawitpro-estate/src/estate/repository/sql"
	estateuc "github.com/davidyunus/sawitpro-estate/src/estate/usecase"
	palmtreelocation "github.com/davidyunus/sawitpro-estate/src/palm_tree/repository/sql"
)

var (
	dbConn               *sql.DB
	estateUsecase        domain.EstateUsecase
	estateRepo           domain.EstateRepository
	palmTreeLocationRepo domain.PalmTreeLocationRepository

	manager *helper.Manager
)

func initDB() (err error) {
	psqlInfo := fmt.Sprint("host=localhost port=5432 user=postgres " +
		"password=postgres dbname=postgres sslmode=disable")

	dbConn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	manager = helper.NewManager(dbConn, common.TransactionContextKey)

	return nil
}

func initRepo() error {
	estateRepo = estatesql.NewEstateRepositorySql(dbConn, manager)
	palmTreeLocationRepo = palmtreelocation.NewPalmTreeRepositorySql(dbConn, manager)

	return nil
}

func initUsecase() error {
	estateUsecase = estateuc.NewEstateUsecase(estateRepo, palmTreeLocationRepo)

	return nil
}

func initHTTP() error {
	e := echo.New()
	e.Debug = true

	estatehttp.NewEstateHandler(e, estateUsecase)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, helper.Response(http.StatusOK, "Pong", nil, nil))
	})

	return e.Start(":8080")
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	err = initRepo()
	if err != nil {
		log.Fatal(err)
	}
	err = initUsecase()
	if err != nil {
		log.Fatal(err)
	}
	err = initHTTP()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/routes"
	"github.com/RamazanZholdas/KeyboardistSV2/utils"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app, err := app.Intitialize(os.Getenv("MONGO_URI"), os.Getenv("DATABASE_NAME"))
	if err != nil {
		log.Panic(err)
	}

	routes.Setup(app.Fiber)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		utils.LogInfo("Shutting down server...")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to create a backup of the database? (y/n): ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if strings.ToLower(answer) == "y" {
			err := exec.Command("sh", "../../scripts/createMongoDbDump.sh").Run()
			if err != nil {
				utils.LogError(fmt.Sprintf("Error creating backup: %v", err))
			} else {
				utils.LogInfo("Backup created successfully")
			}
		}

		app.Close()
	}()

	if err := app.Fiber.Listen("0.0.0.0:3000"); err != nil {
		log.Panic(err)
	}

	utils.LogInfo("Running cleanning up tasks...")
}

package commands

import (
	"encoding/json"
	"log"
	"telebot-trading/app/helper"
	"telebot-trading/app/models"
	"telebot-trading/external/db"
	"time"

	"github.com/jaswdr/faker"
	"github.com/spf13/cobra"
)

func GenerateAdminCommand() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "generate:user"

	cmd.Short = "Run Database Migration"

	cmd.Long = `Run Database Migration`

	cmd.Run = func(cmd *cobra.Command, args []string) {
		admin := getFakeAdminData()
		result := db.GetDB().Create(&admin)
		if result.Error != nil {
			log.Print("Error generate admin data")
			log.Print(result.Error.Error())
		}

		log.Printf("Generate Admin data did run successfully")
		s, _ := json.MarshalIndent(admin, "", "\t")
		log.Print(string(s))
	}

	return cmd
}

func getFakeAdminData() models.Admin {
	admin := models.Admin{}
	faker := faker.New()
	admin.Name = faker.Person().Name()
	admin.Email = faker.Person().Contact().Email
	admin.Password, _ = helper.MakeHash("admin123")
	admin.Role = models.ROLE_SUPER
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()
	return admin
}

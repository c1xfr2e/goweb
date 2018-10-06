package command

import (
	"fmt"
	"os"

	"github.com/bluecover/lm/mocking"
	"github.com/bluecover/lm/models"
	"github.com/jinzhu/gorm"
	"github.com/spf13/pflag"
)

func RunCommand(db *gorm.DB) {
	switch os.Args[1] {
	case "pushfig":
		figcmd := pflag.NewFlagSet("pushfig", pflag.ExitOnError)
		figpath := figcmd.StringP("path", "p", "./figures", "figure file dir")
		dataset := figcmd.StringP("dataset", "d", "", "name of dataset")
		figure := figcmd.StringP("figure", "f", "", "path of figure file")
		table := figcmd.StringP("table", "t", "", "db table to push into")
		figcmd.Parse(os.Args[2:])

		err := db.AutoMigrate(models.Figure{}).Error
		if err != nil {
			fmt.Println("AutoMigrate failed: ", err)
			os.Exit(2)
		}

		if len(*figure) > 0 {
			if len(*table) == 0 {
				fmt.Println("use -t to set table to push")
				os.Exit(2)
			}
			err := PushFigure(*figure, db, *table)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			PushAllFigures(*figpath, *dataset, db)
		}

	case "mock":
		mocking.CreateMockingData(db)

	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

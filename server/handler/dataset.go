package handler

import (
	"encoding/json"
	"time"

	"github.com/bluecover/lm/business/auth"
	"github.com/bluecover/lm/business/timing"
	"github.com/bluecover/lm/figure_parser"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/middleware/authware"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
)

// DatasetList handles dataset list request
func DatasetList(db *gorm.DB) gin.HandlerFunc {

	type CoreIndex struct {
		Type            int    `json:"type"`
		Amount          string `json:"amount"`
		Name            string `json:"name"`
		UpdateFrequency string `json:"updateFrequency"`
	}

	type Overview struct {
		CoreIndex CoreIndex `json:"coreIndex"`
		Detail    []string  `json:"detail"`
	}

	type Dataset struct {
		ID         string      `json:"id"`
		Name       string      `json:"name"`
		Icon       string      `json:"icon"`
		Access     bool        `json:"access"`
		UpdatedAt  int64       `json:"updatedAt"`
		Overview   Overview    `json:"overview"`
		FigureSets interface{} `json:"figureSets"`
	}

	return func(c *gin.Context) {
		userID := authware.GetCurrentUserID(c)
		userDatasets := auth.GetAuthorizedDatasets(db, userID)

		resultDatasets := make([]Dataset, 0)
		for _, dataset := range userDatasets {
			datasetMessages := models.GetDatasetMessages(db, dataset.ID)
			msgs := make([]string, 0)
			for _, m := range datasetMessages {
				msgs = append(msgs, m.Msg)
			}

			var figureSet map[string]interface{}
			if len(dataset.FigureSet) > 0 {
				json.Unmarshal([]byte(dataset.FigureSet), &figureSet)
			}

			coreIndex := CoreIndex{}

			for {
				coreIndexQuery := dataset.CoreIndexQuery
				if len(coreIndexQuery) <= 0 {
					break
				}
				var fig map[string]interface{}
				err := json.Unmarshal([]byte(coreIndexQuery), &fig)
				if err != nil {
					break
				}

				args := figure_parser.ParseArgs{
					Start:  now.New(time.Now().UTC().Add(-24 * time.Hour)).BeginningOfDay(),
					End:    now.New(time.Now().UTC().Add(-24 * time.Hour)).BeginningOfDay(),
					Period: timing.PeriodDate,
				}
				result, err := figure_parser.ParseQuery(fig, args, db)
				if err != nil {
					break
				}

				coreIndex.Name = dataset.CoreIndexName
				coreIndex.Type = result["inc_or_dec"].(int)
				coreIndex.Amount = result["change"].(string)
				coreIndex.UpdateFrequency = dataset.CoreIndexUpdateFrequency
				break
			}

			resultDatasets = append(
				resultDatasets,
				Dataset{
					ID: dataset.Name, Icon: dataset.Icon, Access: true, Name: dataset.DisplayName,
					UpdatedAt: dataset.IndexUpdatedAt.UTC().Unix(),
					Overview: Overview{
						CoreIndex: coreIndex,
						Detail:    msgs,
					},
					FigureSets: figureSet["FigureSet"],
				},
			)
		}

		render.OK(c, gin.H{"datasets": resultDatasets})
	}
}

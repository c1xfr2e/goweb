package handler

import (
	"encoding/json"
	"strings"

	"github.com/bluecover/lm/business/auth"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/server/middleware/authware"
	"github.com/bluecover/lm/server/render"
	"github.com/bluecover/lm/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/oliveagle/jsonpath"
)

type SearchOp struct {
	DB        *gorm.DB
	SearchSet []SearchResult
}

type SearchResult struct {
	Dataset    string `json:"dataset"`    // dataset id
	IconID     string `json:"iconId"`     // dataset icon id, normally the same as dataset id
	FigurePage string `json:"figurePage"` // figure page id
	Text       string `json:"text"`       // figure page title, search target
	FigureID   string `json:"figureId"`
}

func (r *SearchOp) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := authware.GetCurrentUserID(c)
		word := c.Query("word")

		if strings.TrimSpace(word) == "" {
			r.GetLastNSearchHistory(c, userID, 6)
			return
		}

		userDatasets := auth.GetAuthorizedDatasets(r.DB, userID)
		if len(userDatasets) == 0 {
			render.OK(c, []SearchResult{})
			return
		}

		var results []SearchResult
		var datasetIDs []string
		for _, dataset := range userDatasets {
			datasetIDs = append(datasetIDs, dataset.Name)
		}

		for _, sr := range r.SearchSet {
			if !util.Contains(datasetIDs, sr.Dataset) {
				continue
			}
			if !strings.Contains(strings.ToLower(sr.FigureID), strings.ToLower(word)) {
				continue
			}
			results = append(
				results,
				SearchResult{
					Dataset: sr.Dataset, IconID: sr.IconID,
					FigurePage: sr.FigurePage, Text: sr.Text,
					FigureID: sr.FigureID,
				},
			)
		}

		render.OK(c, gin.H{"results": results})
	}
}

func (r *SearchOp) LoadAllFigures() {
	var figurePages []models.FigurePage
	r.DB.Find(&figurePages)

	for _, fp := range figurePages {
		var json_data interface{}
		json.Unmarshal([]byte(fp.Data), &json_data)

		figurePageID, _ := jsonpath.JsonPathLookup(json_data, "$.id")
		datasetID_iconID := strings.Split(figurePageID.(string), ".")[0]
		title, _ := jsonpath.JsonPathLookup(json_data, "$.title")
		ids, _ := jsonpath.JsonPathLookup(json_data, "$.dataView[:].figures[:].figures[:].id")
		var _ids []string
		for _, l1 := range ids.([]interface{}) {
			for _, l2 := range l1.([]interface{}) {
				for _, l3 := range l2.([]interface{}) {
					_ids = append(_ids, l3.(string))
				}
			}
		}
		for _, id := range _ids {
			_id := id
			r.SearchSet = append(r.SearchSet, SearchResult{
				Dataset:    datasetID_iconID,
				IconID:     datasetID_iconID,
				FigurePage: figurePageID.(string),
				Text:       title.(string),
				FigureID:   _id,
			})
		}
	}
}

func (r *SearchOp) SaveSearchHistory() gin.HandlerFunc {
	type ReqData struct {
		DatasetID    string `json:"dataset" valid:"required"`
		FigurePageID string `json:"figurePage" valid:"required"`
	}

	return func(c *gin.Context) {
		user_id := authware.GetCurrentUserID(c)

		op := &models.SearchHistoryDBOp{r.DB}
		op.SaveSearchHistory(user_id, c.PostForm("dataset"), c.PostForm("figurePage"))

		render.OK(c, nil)
	}
}

func (r *SearchOp) GetLastNSearchHistory(c *gin.Context, userId uint, n int) {
	userID := authware.GetCurrentUserID(c)
	op := &models.SearchHistoryDBOp{r.DB}
	searchHistories := op.GetLatestN(userID, n)
	if len(searchHistories) == 0 {
		render.OK(c, []SearchResult{})
		return
	}

	var results []SearchResult
	for _, historyItem := range searchHistories {
		for _, sr := range r.SearchSet {
			if sr.FigurePage == historyItem.FigurePageID && sr.Dataset == historyItem.DatasetID {
				results = append(results, sr)
				break
			}
		}
	}

	render.OK(c, gin.H{"results": results})
}

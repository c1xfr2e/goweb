package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bluecover/lm/models"
	"github.com/jinzhu/gorm"
)

const (
	figureTable = "figures"
)

func PushAllFigures(figRootDir string, datasetName string, db *gorm.DB) error {
	dir, err := ioutil.ReadDir(figRootDir)
	if err != nil {
		return fmt.Errorf("open dir %s failed %s", dir, err)
	}

	for _, f := range dir {
		if !f.IsDir() {
			continue
		}
		if len(datasetName) > 0 && f.Name() != datasetName {
			continue
		}

		datasetDir := path.Join(figRootDir, f.Name())
		PushFigureSet(path.Join(datasetDir, "dataset.json"), db)
		PushFiguresInDir(path.Join(datasetDir, "figures"), db, "figures")
		PushFiguresInDir(path.Join(datasetDir, "pages"), db, "figure_pages")
	}

	return nil
}

func PushFigureSet(datasetPath string, db *gorm.DB) error {
	bytes, err := ioutil.ReadFile(path.Join(datasetPath, "dataset.json"))
	if err != nil {
		return err
	}

	var dataset map[string]interface{}
	json.Unmarshal(bytes, &dataset)

	err = db.Model(models.Dataset{}).
		Where("name=?", dataset["name"].(string)).
		UpdateColumns(models.Dataset{FigureSet: string(bytes)}).Error
	if err != nil {
		return err
	}

	return nil
}

func PushFiguresInDir(dir string, db *gorm.DB, table string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("open dir %s failed %s", dir, err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		err := PushFigure(path.Join(dir, f.Name()), db, table)
		if err != nil {
			fmt.Println(f.Name(), err)
			continue
		}
	}

	return nil
}

func PushFigure(path string, db *gorm.DB, table string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if len(bytes) < 4 {
		return fmt.Errorf("figure file too small: %s", path)
	}

	var figure map[string]interface{}
	if json.Unmarshal(bytes, &figure) != nil {
		return err
	}

	figureID, ok := figure["id"].(string)
	if !ok {
		return err
	}

	sql := fmt.Sprintf("delete from %s where data->>'id'=?", table)
	err = db.Exec(sql, figureID).Error
	if err != nil {
		return fmt.Errorf("delete figure from %s failed: %s %s", table, figureID, err)
	}

	sql = fmt.Sprintf("insert into %s (data) values (?)", table)
	err = db.Exec(sql, string(bytes)).Error
	if err != nil {
		return fmt.Errorf("insert figure into %s failed: %s %s", table, figureID, err)
	}

	return nil
}

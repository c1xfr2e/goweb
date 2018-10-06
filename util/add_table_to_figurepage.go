package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

const (
	figuredir = "figures"
	pagedir   = "pages"
)

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

type tfig struct {
	Filepath string
	Obj      map[string]interface{}
}

func loadFigureFiles(dir string) map[string]tfig {
	files, err := ioutil.ReadDir(dir)
	panicerr(err)

	figset := make(map[string]tfig, 0)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		figFilepath := path.Join(dir, f.Name())
		bytes, err := ioutil.ReadFile(figFilepath)
		panicerr(err)
		var figure map[string]interface{}
		err = json.Unmarshal(bytes, &figure)
		panicerr(err)
		figset[figure["id"].(string)] = tfig{figFilepath, figure}
	}

	return figset
}

func datasetDirs(figRoot string) []string {
	dir, err := ioutil.ReadDir(figRoot)
	panicerr(err)
	dsdir := make([]string, 0)
	for _, f := range dir {
		if f.IsDir() {
			dsdir = append(dsdir, path.Join(figRoot, f.Name()))
		}
	}
	return dsdir
}

func getTableInQuery(query map[string]interface{}) string {
	for k, v := range query {
		switch v.(type) {
		case string:
			if k == "table" {
				return v.(string)
			}
		case map[string]interface{}:
			return v.(map[string]interface{})["table"].(string)
		}
	}
	panic(fmt.Errorf("can not find table in query %s", query))
	return ""
}

func WalkPagesForTable(figRoot string) {
	datasetDirs := datasetDirs(figRoot)
	for _, dir := range datasetDirs {
		figs := loadFigureFiles(path.Join(dir, figuredir))
		pages := loadFigureFiles(path.Join(dir, pagedir))
		for _, tpage := range pages {
			firstFigure := tpage.Obj["dataView"].([]interface{})[0].(map[string]interface{})["figures"].([]interface{})[0].(map[string]interface{})["figures"].([]interface{})[0].(map[string]interface{})
			figureID := firstFigure["id"].(string)
			tfig, ok := figs[figureID]
			if !ok {
				log.Printf("unkonw figure id %s", figureID)
				continue
			}
			figure := tfig.Obj

			table := getTableInQuery(figure["#query"].(map[string]interface{}))
			tpage.Obj["table"] = table

			bytes, err := json.MarshalIndent(tpage.Obj, "", "    ")
			panicerr(err)
			err = ioutil.WriteFile(tpage.Filepath, bytes, 0644)
			panicerr(err)
		}
	}
}

/*
   "dataView":[
       {
           "type":"chartView",
           "figures":[
               {
                   "type":"figureBox",
                   "figures":[
                       {
                           "type":"LineChart",
                           "id":"YRD.BorrowerProfile.PersonasOfBorrowers.NumberOfBorrowersServed"
                       },
                       {
                           "type":"kvCard",
                           "id":"YRD.BorrowerProfile.PersonasOfBorrowers.NumberOfBorrowersServed.KV"
                       }
                   ]
               }
           ]
       }
   ]
*/

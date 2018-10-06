package util

import (
	"testing"
	"fmt"
	"encoding/json"
	"github.com/oliveagle/jsonpath"
)

func TestFlatten(t *testing.T) {
	jsonString := `
{
	"id": "YRD.Overview",
	"title":"Overview",
	"dateView":0,
	"dataView":[
		{
			"type":"chartView",
			"figures":[
				{
					"type":"figureGroup",
					"title":"Overview",
					"figures":[
						{
							"type":"OverviewCard",
							"id":"YRD.Overview.Loans"
						},
						{
							"type":"OverviewCard",
							"id":"YRD.Overview.LoanFacilitation"
						}
					]
				}
			]
		},
		{
			"type":"chartView",
			"figures":[
				{
					"type":"figureGroup",
					"title":"Overview",
					"figures":[
						{
							"type":"OverviewCard",
							"id":"YRD.Overview.Loans"
						},
						{
							"type":"OverviewCard",
							"id":"YRD.Overview.LoanFacilitation"
						}
					]
				}
			]
		}
	]
}`

	var json_data interface{}
	json.Unmarshal([]byte(jsonString), &json_data)
	ids, _ := jsonpath.JsonPathLookup(json_data, "$.dataView[:].figures[:].figures[:].id")
	fmt.Println(ids)
	fmt.Println([]string{"1"})
	for _, l1 := range ids.([]interface{}) {
		for _, l2 := range l1.([]interface{}) {
			for _, l3 := range l2.([]interface{}) {
				fmt.Println(l3.(string))
			}
		}
	}
}

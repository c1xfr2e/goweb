package handler

import (
	"encoding/json"
	"testing"
)

var tableFigureTests = []struct {
	in  string
	out string
}{
	{
		`{
    "dataView": [
        {
            "figures": [
                {
                    "figures": [
                        {
                            "id": "HTHT.HotelsAndRooms.NumberOfIncreasedHotels",
                            "type": "LineChart"
                        },
                        {
                            "id": "HTHT.HotelsAndRooms.NumberOfIncreasedHotels.KV",
                            "type": "kvCard"
                        }
                    ],
                    "title": "Number of Increased Hotels",
                    "type": "figureGroup"
                },
                {
                    "figures": [
                        {
                            "id": "HTHT.HotelsAndRooms.NumberOfDecreasedHotels",
                            "type": "LineChart"
                        },
                        {
                            "id": "HTHT.HotelsAndRooms.NumberOfDecreasedHotels.KV",
                            "type": "kvCard"
                        }
                    ],
                    "title": "Number of Decreased Hotels",
                    "type": "figureGroup"
                }
            ],
            "type": "chartView"
        },
        {
            "figures": [
                {
                    "id": "HTHT.HotelsAndRooms.HotelDevelopment.Table",
                    "type": "table"
                }
            ],
            "type": "tableView"
        }
    ],
    "dateView": 14,
    "id": "HTHT.HotelsAndRooms.HotelDevelopment",
    "table": "HTHT.hotel",
    "title": "Hotel Development"
}`,
		"HTHT.HotelsAndRooms.HotelDevelopment.Table",
	},
	{
		`{
    "dataView": [
        {
            "figures": [
                {
                    "figures": [
                        {
                            "id": "YRD.LoanFacilitation.AverageLoanSize",
                            "type": "LineChart"
                        },
                        {
                            "id": "YRD.LoanFacilitation.AverageLoanSize.KV",
                            "type": "kvCard"
                        }
                    ],
                    "type": "figureBox"
                }
            ],
            "type": "chartView"
        },
        {
            "figures": [
                {
                    "id": "YRD.LoanFacilitation.AverageLoanSize.Table",
                    "type": "table"
                }
            ],
            "type": "tableView"
        }
    ],
    "dateView": 15,
    "id": "YRD.LoanFacilitation.AverageLoanSize",
    "table": "YRD.loans",
    "title": "Average Loan Size (RMB)"
}`,
		"YRD.LoanFacilitation.AverageLoanSize.Table",
	},
}

func TestGetTableFigureID(t *testing.T) {
	for _, tt := range tableFigureTests {
		var j interface{}
		err := json.Unmarshal([]byte(tt.in), &j)
		if err != nil {
			t.Fatal(err)
		}
		id := getTableFigureID(j)
		if id != tt.out {
			t.Error("wrong id:", id, "; want:", tt.out)
		}
	}
}

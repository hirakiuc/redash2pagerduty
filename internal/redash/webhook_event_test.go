package redash

import (
	"strings"
	"testing"
)

const ValidJson = `
{
	"url_base": "http://redash.sitedomain.jp",
	"event": "alert_state_change",
	"alert": {
		"state": "triggered",
		"query_id": 1,
		"name": "fuga: count(*) greater than 5000",
		"rearm": null,
		"updated_at": "2015-12-11T08:09:04.898736",
		"user_id": 1,
		"created_at": "2015-12-11T08:08:58.228976+00:00",
		"last_triggered_at": "2015-12-11T08:09:04.898690+00:00",
		"id": 6,
		"options": {
			"column": "count(*)",
			"value": 5000,
			"op": "greater than"
		}
	}
}
`

func TestParseWithValidJson(t *testing.T) {
	examples := []string{
		ValidJson,
	}

	for _, str := range examples {
		reader := strings.NewReader(str)

		result, err := Parse(reader)
		if err != nil {
			t.Errorf("Parse Failed. %v", err)
		}
		if result == nil {
			t.Errorf("Result should not nil")
		}
	}
}

package gorm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type WeeklyHours [7][2]int

func (w *WeeklyHours) String() string {
	return fmt.Sprintf(
		"{{%d,%d},{%d,%d},{%d,%d},{%d,%d},{%d,%d},{%d,%d},{%d,%d}}",
		w[0][0], w[0][1], w[1][0], w[1][1], w[2][0], w[2][1], w[3][0], w[3][1], w[4][0], w[4][1], w[5][0], w[5][1], w[6][0], w[6][1])
}

func (w WeeklyHours) Value() (driver.Value, error) {
	return w.String(), nil
}

func (w *WeeklyHours) Scan(val interface{}) error {
	pgArray := string(val.([]uint8))
	jsonArray := strings.Replace(pgArray, "{", "[", -1)
	jsonArray = strings.Replace(jsonArray, "}", "]", -1)
	err := json.Unmarshal([]byte(jsonArray), w)
	return err
}

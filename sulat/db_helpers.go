package sulat

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func scanJson(src any, dest any, typeName string) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), dest)
	case []byte:
		return json.Unmarshal(v, dest)
	default:
		return fmt.Errorf("cannot scan value of type %T into %s", v, typeName)
	}
}

func driverValueJson(src any) (driver.Value, error) {
	js, err := json.Marshal(src)
	return string(js), err
}

package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/ucard/global"
)

type AdminConfig struct {
	global.GVA_MODEL
	Kay          string    `gorm:"column:kay;unique" json:"kay"`
	ValueType    string    `gorm:"column:value_type;type:ENUM('string', 'number', 'boolean', 'json', 'date'); NOT NULL" json:"valueType"`
	StringValue  *string   `gorm:"column:string_value;type:text" json:"stringValue"`
	NumberValue  *int64    `gorm:"column:number_value;type:bigint" json:"numberValue"`
	BooleanValue *bool     `gorm:"column:boolean_value;type:boolean" json:"booleanValue"`
	JsonValue    *string   `gorm:"column:json_value;type:text" json:"jsonValue"`
	DateValue    time.Time `gorm:"column:date_value;type:datetime" json:"dateValue"`
	Operator     string    `gorm:"column:operator;type:varchar(255)" json:"operator"`
}

func (AdminConfig) TableName() string {
	return "admin_configs"
}

type MapStringString map[string]bool

func (m *MapStringString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (m MapStringString) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type MapStringBool map[string]bool

func (m *MapStringBool) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (m MapStringBool) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type ListString []string

func (m *ListString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (m ListString) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// SliceUint 用于存储 uint 数组到 JSON 字段
type SliceUint []uint

func (s *SliceUint) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}
	if len(b) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(b, s)
}

func (s SliceUint) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

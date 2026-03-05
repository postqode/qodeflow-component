package saphana

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the SAP HANA activity
type Settings struct {
	// DSN is the full SAP HANA data source name (e.g. hdb://user:pass@host:39017)
	// If DSN is provided, Host/Port/User/Password are ignored.
	DSN string `md:"dsn"`
	// Host is the SAP HANA server hostname or IP address
	Host string `md:"host"`
	// Port is the SAP HANA server port (default: 39017)
	Port int `md:"port"`
	// User is the database username
	User string `md:"user"`
	// Password is the database password
	Password string `md:"password"`
}

// Input represents the input for the SAP HANA activity
type Input struct {
	// Method is the operation to perform: QUERY, EXEC, or CALL
	Method string `md:"method"`
	// Query is the SQL query or stored procedure call
	Query string `md:"query"`
	// Args is an optional list of positional query arguments
	Args []any `md:"args"`
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.Query, err = coerce.ToString(values["query"])
	if err != nil {
		return err
	}
	if v, ok := values["args"]; ok && v != nil {
		switch t := v.(type) {
		case []any:
			i.Args = t
		default:
			i.Args = []any{v}
		}
	}
	return nil
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method": i.Method,
		"query":  i.Query,
		"args":   i.Args,
	}
}

// Output represents the output for the SAP HANA activity
type Output struct {
	// Result contains the rows returned by a QUERY or CALL operation
	Result []map[string]any `md:"result"`
	// RowsAffected is the number of rows affected by an EXEC operation,
	// or the number of rows returned by a QUERY/CALL operation
	RowsAffected int64 `md:"rowsAffected"`
}

func (o *Output) FromMap(values map[string]any) error {
	var err error
	o.RowsAffected, err = coerce.ToInt64(values["rowsAffected"])
	if err != nil {
		return err
	}
	if v, ok := values["result"]; ok && v != nil {
		switch t := v.(type) {
		case []map[string]any:
			o.Result = t
		case []any:
			rows := make([]map[string]any, 0, len(t))
			for _, item := range t {
				if m, ok := item.(map[string]any); ok {
					rows = append(rows, m)
				}
			}
			o.Result = rows
		}
	}
	return nil
}

func (o *Output) ToMap() map[string]any {
	return map[string]any{
		"result":       o.Result,
		"rowsAffected": o.RowsAffected,
	}
}

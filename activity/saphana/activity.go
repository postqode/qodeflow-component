package saphana

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/SAP/go-hdb/driver"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&SapHanaActivity{}, New)
}

/*
Integration with SAP HANA
inputs: {method, query, args}
outputs: {result, rowsAffected}
*/
type SapHanaActivity struct {
	settings *Settings
}

// New creates a new SAP HANA activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	if s.DSN == "" {
		if s.Host == "" || s.User == "" {
			return nil, fmt.Errorf("saphana activity: either 'dsn' or both 'host' and 'user' settings are required")
		}
		port := s.Port
		if port == 0 {
			port = 39017
		}
		s.DSN = fmt.Sprintf("hdb://%s:%s@%s:%d", s.User, s.Password, s.Host, port)
	}

	return &SapHanaActivity{settings: s}, nil
}

// Metadata returns the activity's metadata
func (a *SapHanaActivity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval implements activity.Activity.Eval – SAP HANA integration
func (a *SapHanaActivity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.Query == "" {
		return false, fmt.Errorf("saphana activity: 'query' input is required")
	}

	db, err := sql.Open("hdb", a.settings.DSN)
	if err != nil {
		ctx.Logger().Errorf("SAP HANA connection error: %v", err)
		return false, fmt.Errorf("saphana activity: failed to open connection: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			ctx.Logger().Errorf("Error closing SAP HANA connection: %v", closeErr)
		}
	}()

	if pingErr := db.PingContext(context.Background()); pingErr != nil {
		ctx.Logger().Errorf("SAP HANA ping error: %v", pingErr)
		return false, fmt.Errorf("saphana activity: failed to connect to SAP HANA: %w", pingErr)
	}

	// Build args slice
	args := make([]any, len(input.Args))
	copy(args, input.Args)

	switch strings.ToUpper(input.Method) {
	case "QUERY":
		rows, err := db.QueryContext(context.Background(), input.Query, args...)
		if err != nil {
			return false, fmt.Errorf("saphana activity: query failed: %w", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return false, fmt.Errorf("saphana activity: failed to get columns: %w", err)
		}

		var results []map[string]any
		for rows.Next() {
			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				return false, fmt.Errorf("saphana activity: failed to scan row: %w", err)
			}

			row := make(map[string]any)
			for i, col := range columns {
				row[col] = values[i]
			}
			results = append(results, row)
		}

		if err := rows.Err(); err != nil {
			return false, fmt.Errorf("saphana activity: rows iteration error: %w", err)
		}

		ctx.Logger().Debugf("QUERY returned %d rows", len(results))
		_ = ctx.SetOutputObject(&Output{
			Result:       results,
			RowsAffected: int64(len(results)),
		})

	case "EXEC":
		result, err := db.ExecContext(context.Background(), input.Query, args...)
		if err != nil {
			return false, fmt.Errorf("saphana activity: exec failed: %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			ctx.Logger().Warnf("Could not retrieve rows affected: %v", err)
		}

		ctx.Logger().Debugf("EXEC affected %d rows", affected)
		_ = ctx.SetOutputObject(&Output{RowsAffected: affected})

	case "CALL":
		rows, err := db.QueryContext(context.Background(), input.Query, args...)
		if err != nil {
			return false, fmt.Errorf("saphana activity: call failed: %w", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return false, fmt.Errorf("saphana activity: failed to get columns from procedure result: %w", err)
		}

		var results []map[string]any
		for rows.Next() {
			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				return false, fmt.Errorf("saphana activity: failed to scan procedure result row: %w", err)
			}

			row := make(map[string]any)
			for i, col := range columns {
				row[col] = values[i]
			}
			results = append(results, row)
		}

		if err := rows.Err(); err != nil {
			return false, fmt.Errorf("saphana activity: procedure rows iteration error: %w", err)
		}

		ctx.Logger().Debugf("CALL returned %d rows", len(results))
		_ = ctx.SetOutputObject(&Output{
			Result:       results,
			RowsAffected: int64(len(results)),
		})

	default:
		ctx.Logger().Errorf("unsupported method '%s'", input.Method)
		return false, fmt.Errorf("saphana activity: unsupported method '%s' (supported: QUERY, EXEC, CALL)", input.Method)
	}

	return true, nil
}

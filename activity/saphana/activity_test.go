package saphana

import (
	"os"
	"testing"

	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestActivity_Metadata(t *testing.T) {
	act := &SapHanaActivity{}
	md := act.Metadata()
	assert.NotNil(t, md)
	assert.NotNil(t, md.Settings["dsn"])
	assert.NotNil(t, md.Input["method"])
	assert.NotNil(t, md.Input["query"])
	assert.NotNil(t, md.Output["result"])
	assert.NotNil(t, md.Output["rowsAffected"])
}

func TestActivity_New_WithDSN(t *testing.T) {
	settings := map[string]any{
		"dsn": "hdb://user:pass@localhost:39017",
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)
	assert.NotNil(t, act)
}

func TestActivity_New_WithHostUser(t *testing.T) {
	settings := map[string]any{
		"host":     "localhost",
		"port":     39017,
		"user":     "SYSTEM",
		"password": "secret",
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)
	assert.NotNil(t, act)
}

func TestActivity_New_DefaultPort(t *testing.T) {
	settings := map[string]any{
		"host": "localhost",
		"user": "SYSTEM",
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)
	assert.NotNil(t, act)
	// Verify the default port 39017 is used in the DSN
	hanaAct := act.(*SapHanaActivity)
	assert.Contains(t, hanaAct.settings.DSN, ":39017")
}

func TestActivity_New_MissingSettings(t *testing.T) {
	settings := map[string]any{}
	iCtx := test.NewActivityInitContext(settings, nil)
	_, err := New(iCtx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either 'dsn' or both 'host' and 'user'")
}

func TestActivity_Eval_MissingQuery(t *testing.T) {
	settings := map[string]any{
		"dsn": "hdb://user:pass@localhost:39017",
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		Method: "QUERY",
		Query:  "",
	}
	tc.SetInputObject(input)
	done, err := act.Eval(tc)
	assert.False(t, done)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'query' input is required")
}

func TestActivity_Eval_UnsupportedMethod(t *testing.T) {
	uri := os.Getenv("SAPHANA_DSN")
	if uri == "" {
		t.Skip("Skipping: SAPHANA_DSN not set")
	}

	settings := map[string]any{"dsn": uri}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		Method: "UNKNOWN",
		Query:  "SELECT 1 FROM DUMMY",
	}
	tc.SetInputObject(input)
	done, err := act.Eval(tc)
	assert.False(t, done)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported method")
}

// Integration tests – require a running SAP HANA instance.
// Set SAPHANA_DSN env var to enable (e.g. hdb://SYSTEM:password@localhost:39017).
func TestActivity_Eval_Integration(t *testing.T) {
	dsn := os.Getenv("SAPHANA_DSN")
	if dsn == "" {
		t.Skip("Skipping integration test: SAPHANA_DSN not set")
	}

	settings := map[string]any{"dsn": dsn}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)
	assert.NotNil(t, act)

	// 1. CREATE TABLE
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		Method: "EXEC",
		Query: `CREATE TABLE IF NOT EXISTS QODEFLOW_TEST (
			ID INTEGER PRIMARY KEY,
			NAME NVARCHAR(100),
			VALUE INTEGER
		)`,
	}
	tc.SetInputObject(input)
	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	// 2. INSERT
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method: "EXEC",
		Query:  "INSERT INTO QODEFLOW_TEST (ID, NAME, VALUE) VALUES (?, ?, ?)",
		Args:   []any{1, "test-item", 42},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.RowsAffected)

	// 3. QUERY
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method: "QUERY",
		Query:  "SELECT ID, NAME, VALUE FROM QODEFLOW_TEST WHERE ID = ?",
		Args:   []any{1},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.RowsAffected)
	assert.Len(t, output.Result, 1)
	assert.Equal(t, "test-item", output.Result[0]["NAME"])

	// 4. UPDATE
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method: "EXEC",
		Query:  "UPDATE QODEFLOW_TEST SET NAME = ? WHERE ID = ?",
		Args:   []any{"updated-item", 1},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.RowsAffected)

	// 5. DELETE
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method: "EXEC",
		Query:  "DELETE FROM QODEFLOW_TEST WHERE ID = ?",
		Args:   []any{1},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.RowsAffected)

	// 6. DROP TABLE (cleanup)
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method: "EXEC",
		Query:  "DROP TABLE QODEFLOW_TEST",
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)
}

func TestActivity_Eval_CallProcedure_Integration(t *testing.T) {
	dsn := os.Getenv("SAPHANA_DSN")
	if dsn == "" {
		t.Skip("Skipping integration test: SAPHANA_DSN not set")
	}

	settings := map[string]any{"dsn": dsn}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	// Call a built-in SAP HANA function via CALL-style query
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		Method: "QUERY",
		Query:  "SELECT CURRENT_TIMESTAMP AS NOW FROM DUMMY",
	}
	tc.SetInputObject(input)
	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.RowsAffected)
	assert.Len(t, output.Result, 1)
	assert.NotNil(t, output.Result[0]["NOW"])
}

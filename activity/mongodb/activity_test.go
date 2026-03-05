package mongodb

import (
	"os"
	"testing"

	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestActivity_Metadata(t *testing.T) {
	act := &MongoDbActivity{}
	md := act.Metadata()
	assert.NotNil(t, md)
	assert.NotNil(t, md.Settings["uri"])
	assert.NotNil(t, md.Input["method"])
}

func TestActivity_New(t *testing.T) {
	settings := map[string]any{
		"uri":        "mongodb://localhost:27017",
		"dbName":     "test",
		"collection": "items",
	}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)
	assert.NotNil(t, act)
}

func TestActivity_Eval_Integration(t *testing.T) {
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB")
	collName := os.Getenv("MONGODB_COLLECTION")

	if uri == "" {
		t.Skip("Skipping integration test: MONGODB_URI not set")
	}
	if dbName == "" {
		dbName = "test_db"
	}
	if collName == "" {
		collName = "test_collection"
	}

	settings := map[string]any{
		"uri":        uri,
		"dbName":     dbName,
		"collection": collName,
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	// 1. INSERT
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		Method:  "INSERT",
		KeyName: "id",
		Data: map[string]any{
			"id":    "123",
			"name":  "test-item",
			"value": 42,
		},
	}
	tc.SetInputObject(input)
	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.NotNil(t, output.Output) // Should be the inserted ID (123 or MongoID)

	// 2. GET
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method:   "GET",
		KeyName:  "id",
		KeyValue: "123",
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.NotNil(t, output.Output)
	retrieved := output.Output.(map[string]any)
	assert.Equal(t, "test-item", retrieved["name"])

	// 3. UPDATE
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method:   "UPDATE",
		KeyName:  "id",
		KeyValue: "123",
		Data: map[string]any{
			"name": "updated-item",
		},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.Count)

	// 4. REPLACE
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method:   "REPLACE",
		KeyName:  "id",
		KeyValue: "123",
		Data: map[string]any{
			"id":    "123",
			"name":  "replaced-item",
			"value": 100,
		},
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.Count)

	// 5. DELETE
	tc = test.NewActivityContext(act.Metadata())
	input = &Input{
		Method:   "DELETE",
		KeyName:  "id",
		KeyValue: "123",
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output = &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.Count)
}

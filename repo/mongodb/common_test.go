package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 測試日期範圍選項
func TestWithDateRange(t *testing.T) {
	opts := options.Find()
	startDate, endDate := time.Now().Add(-24*time.Hour), time.Now()

	WithDateRange("createdAt", startDate, endDate)(opts)

	expected := bson.M{
		"createdAt": bson.M{
			"$gte": startDate.Unix(),
			"$lte": endDate.Unix(),
		},
	}

	assert.Equal(t, expected, opts.Hint)
}

// 測試投影選項
func TestWithProjection(t *testing.T) {
	opts := options.Find()
	trueFields := []string{"_id", "name"}
	falseFields := []string{"password"}

	WithProjection(trueFields, falseFields)(opts)

	expected := bson.M{"_id": 1, "name": 1, "password": 0}
	assert.Equal(t, expected, opts.Projection)
}

// 測試限制選項
func TestWithLimit(t *testing.T) {
	opts := options.Find()
	limit := int64(10)

	WithLimit(limit)(opts)

	assert.Equal(t, limit, *opts.Limit)
}

// 測試偏移選項
func TestWithOffset(t *testing.T) {
	opts := options.Find()
	offset := int64(5)

	WithOffset(offset)(opts)

	assert.Equal(t, offset, *opts.Skip)
}

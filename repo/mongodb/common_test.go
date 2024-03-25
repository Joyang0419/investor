package mongodb

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tools/intergrationtest"
)

var testMongoClient *mongo.Client

// TestMain 是一個特殊的測試，可以用來實現全局設置和清理
func TestMain(m *testing.M) {
	fmt.Println("設置測試")
	// setup code here

	// m.Run() 運行所有的測試
	var resourceMongo *dockertest.Resource
	_, resourceMongo, testMongoClient = intergrationtest.CreateMongoDBContainer("mongoDB_common_test")
	code := m.Run()

	fmt.Println("清理測試")
	// teardown code here
	_ = resourceMongo.Close()

	// 退出測試，返回適當的值
	os.Exit(code)
}

// TestFilterDateRange 是 FilterDateRange 函數的測試函數
func TestFilterDateRange(t *testing.T) {
	// 定義測試用的起始和結束時間
	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)

	// 執行 FilterDateRange 函數
	filter := FilterDateRange("createdAt", startDate, endDate)

	// 構造預期的結果
	expected := map[string]any{
		"createdAt": map[string]any{
			"$gte": startDate.Unix(),
			"$lte": endDate.Unix(),
		},
	}

	// 使用 reflect.DeepEqual 來檢查函數返回值是否與預期一致
	if !reflect.DeepEqual(filter, expected) {
		t.Errorf("FilterDateRange 返回的結果不符合預期，得到 %v，期望 %v", filter, expected)
	}
}

// 測試投影選項
func TestWithProjection(t *testing.T) {
	opts := options.Find()
	trueFields := []string{"_id", "name"}

	WithProjection(trueFields)(opts)

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

var testMongoStorage = MongoStorage{
	Database:   "Database",
	Collection: "Collection",
}

type testDoc struct {
	TestID   string `bson:"testID"`
	DateTime int64  `bson:"DateTime"`
}

func TestAll(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 0; i < docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: time.Now().Unix(),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		nil,
	)

	assert.Len(t, r, docCount)
	assert.NoError(t, err)
}

func TestAllWithLimit(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 0; i < docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: time.Now().Unix(),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	limit := 5
	opt := WithLimit(int64(limit))
	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		nil,
		opt,
	)

	assert.Len(t, r, limit)
	assert.NoError(t, err)
}

func TestAllWithSkip(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 0; i < docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: time.Now().Unix(),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	skip := 5
	opt := WithOffset(int64(skip))
	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		nil,
		opt,
	)

	assert.Len(t, r, docCount-skip)
	assert.NoError(t, err)
}

func TestAllWithOrderBy(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 0; i < docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: int64(i),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	opt := WithOrderBy(map[string]int{"DateTime": -1})
	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		nil,
		opt,
	)

	dateTimes := make([]int64, len(r))
	for idx := range r {
		dateTimes[idx] = r[idx].DateTime
	}

	assert.IsDecreasing(t, dateTimes)

	assert.NoError(t, err)
}

func TestAllWithDateRange(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 1; i <= docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: int64(i),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		FilterDateRange("DateTime", time.Unix(6, 0), time.Unix(10, 0)),
		WithOrderBy(map[string]int{"DateTime": 1}),
	)

	dateTimes := make([]int64, 5)
	for idx := range r {
		dateTimes[idx] = r[idx].DateTime
	}

	assert.IsIncreasing(t, dateTimes)

	expected := []int64{6, 7, 8, 9, 10}
	assert.Equal(t, expected, dateTimes)

	assert.NoError(t, err)
}

func TestAllWithProjection(t *testing.T) {
	ctx := context.TODO()

	var models []mongo.WriteModel
	docCount := 10
	for i := 1; i <= docCount; i++ {
		doc := testDoc{
			TestID:   fmt.Sprintf("testID_%d", i),
			DateTime: int64(i),
		}
		model := mongo.NewInsertOneModel().SetDocument(doc)
		models = append(models, model)
	}

	if _, err := testMongoClient.Database(testMongoStorage.Database).
		Collection(testMongoStorage.Collection).BulkWrite(ctx, models); err != nil {
		t.Fatalf("BulkWrite err: %v", err)
	}
	defer func() {
		_ = testMongoClient.Database(testMongoStorage.Database).
			Collection(testMongoStorage.Collection).Drop(context.TODO())
	}()

	r, err := All[[]testDoc](
		context.TODO(),
		testMongoClient,
		10*time.Second,
		testMongoStorage,
		nil,
		WithProjection([]string{"testID"}),
	)

	for idx := range r {
		assert.NotEmpty(t, r[idx].TestID)
		assert.Equal(t, int64(0), r[idx].DateTime)
	}
	assert.NoError(t, err)
}

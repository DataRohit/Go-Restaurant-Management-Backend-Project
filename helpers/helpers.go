package helpers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection      *mongo.Collection = database.OpenCollection(database.Client, "user")
	foodCollection      *mongo.Collection = database.OpenCollection(database.Client, "food")
	orderCollection     *mongo.Collection = database.OpenCollection(database.Client, "order")
	orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")
)

func RecordExists(ctx context.Context, collection *mongo.Collection, field string, value interface{}) (bool, error) {
	count, err := collection.CountDocuments(ctx, bson.M{field: value})
	return count > 0, err
}

func GetNonNilString(s *string, defaultValue string) string {
	if s == nil || *s == "" {
		return defaultValue
	}
	return *s
}

func GetPaginationParams(c *gin.Context) (int64, int64) {
	recordPerPage, err := strconv.ParseInt(c.Query("recordPerPage"), 10, 64)
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	return recordPerPage, page
}

func InTimeSpan(start, end, check time.Time) bool {
	return !check.Before(start) && !check.After(end)
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

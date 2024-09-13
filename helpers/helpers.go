package helpers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

package queryitem

import (
	"bytes"
	"encoding/json"
	"io"
)

func BuildQueryConditions(items ...QueryItem) io.Reader {
	return buildHttpRequestBody(buildQueryConditions(items))
}

func buildQueryConditions(items []QueryItem) []QueryCondition {
	conditions := make([]QueryCondition, len(items))
	for i, item := range items {
		conditions[i] = QueryCondition{
			PropertyName: item.PropertyName(),
			Value:        item.Value(),
			CompareType:  item.CompareType(),
		}
	}
	return conditions
}

func buildHttpRequestBody(conditions []QueryCondition) io.Reader {
	data := map[string]interface{}{
		"criteria": conditions,
	}

	marshal, _ := json.Marshal(data)
	return bytes.NewBuffer(marshal)
}

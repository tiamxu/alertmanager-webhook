package api

import (
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/alertmanager-webhook/model"
)

func Test(c *gin.Context) {
	alerts := model.Alerts{
		{
			Status: "firing",
			Labels: map[string]string{"level": "1"},
		},
		{
			Status: "firing",
			Labels: map[string]string{"level": "4"},
		},
		{
			Status: "firing",
			Labels: map[string]string{"level": "2"},
		},
		{
			Status: "firing",
			Labels: map[string]string{"level": "3"},
		},
	}

	fmt.Println("Before sorting:")
	for _, alert := range alerts {
		fmt.Printf("%+v\n", alert)
	}

	sort.Sort(alerts)

	fmt.Println("\nAfter sorting:")
	for _, alert := range alerts {
		fmt.Printf("%+v\n", alert)
	}
}

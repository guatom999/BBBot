package utils

import (
	"log"
	"time"
)

type ()

// func GetLocalBkkTime() time.Time {
// 	loc, err := time.LoadLocation("Asia/Bangkok")
// 	if err != nil {
// 		log.Printf("Error: Load Localtime Failed")
// 		panic(err)
// 	}

// 	bkkTime := time.Now().In(loc)
// 	log.Printf("Current time in Bangkok: %v, Zone: %v", bkkTime.Format("15:04:05"), bkkTime.Location())
// 	return bkkTime
// }

func GetLocalBkkTime() time.Time {
	now := time.Now().UTC()

	thaiTime, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error: Load Localtime Failed :%s", err.Error())
	}

	return now.In(thaiTime)
}

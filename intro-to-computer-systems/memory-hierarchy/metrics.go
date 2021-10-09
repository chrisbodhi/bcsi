package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type UserId uint16 // was an int
type UserMap map[UserId]*User

type Address struct {
	fullAddress string
	zip         uint32
}

type DollarAmount struct {
	dollars uint32 // was uint64
	cents   uint8  // was uint64
}

type Payment struct {
	amount DollarAmount
	time   time.Time
}

type User struct {
	id       UserId
	name     string
	age      uint8 // was an int
	address  Address
	payments []Payment
}

func AverageAge(users UserMap) float32 {
	var average float32
	var count float32
	for _, u := range users {
		count += 1
		average += (float32(u.age) - average) / count
	}
	return average
}

func AveragePaymentAmount(users UserMap) float64 {
	average, count := 0.0, 0.0
	for _, u := range users {
		for _, p := range u.payments {
			count += 1
			amount := float64(p.amount.dollars) + float64(p.amount.cents)/100
			average += (amount - average) / count
		}
	}
	return average
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs, count := 0.0, 0.0
	for _, u := range users {
		for _, p := range u.payments {
			count += 1
			amount := float64(p.amount.dollars) + float64(p.amount.cents)/100
			diff := amount - mean
			squaredDiffs += diff * diff
		}
	}
	return math.Sqrt(squaredDiffs / count)
}

func LoadData() UserMap {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make(UserMap, len(userLines))
	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		name := line[1]
		age, _ := strconv.Atoi(line[2])
		address := line[3]
		zip, _ := strconv.Atoi(line[3])
		users[UserId(id)] = &User{UserId(id), name, uint8(age), Address{address, zip}, []Payment{}}
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	for _, line := range paymentLines {
		userID, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		datetime, _ := time.Parse(time.RFC3339, line[1])
		users[UserId(userID)].payments = append(users[UserId(userID)].payments, Payment{
			DollarAmount{uint32(paymentCents / 100), uint8(paymentCents % 100)},
			datetime,
		})
	}

	return users
}

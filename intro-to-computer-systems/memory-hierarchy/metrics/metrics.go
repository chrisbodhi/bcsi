package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

// UserID does the thing
type UserID uint16 // was an int
// UserMap does the thing
type UserMap map[UserID]*User

// Address does it, too
type Address struct {
	fullAddress string
	zip         uint32 // TODO should this be a string or something else more constrained?
}

// DollarAmount does
type DollarAmount struct {
	dollars uint32 // was uint64
	cents   uint8  // was uint64
}

// Payment contains the amount
type Payment struct {
	amount uint32
}

// User maps to the contents in the users.csv
type User struct {
	age      uint      // 8 B
	id       UserID    // 16 B
	payments []Payment // 24 B
}

func centsConv(cents float64) float64 {
	return cents/100 + float64(int(cents)%100/100)
}

// AverageAge calculates the average age field
func AverageAge(users UserMap) float32 {
	var sum uint
	for _, u := range users {
		sum += u.age
	}
	return float32(sum) / float32(len(users))
}

// AveragePaymentAmount calculates the average payment across all payments
func AveragePaymentAmount(users UserMap) float64 {
	count, sum := 0, 0
	for _, u := range users {
		for _, p := range u.payments {
			count++
			sum += int(p.amount)
		}
	}
	totalCents := float64(sum) / float64(count)
	return centsConv(totalCents)
}

// StdDevPaymentAmount computes the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs := 0.0
	count := 0
	for _, u := range users {
		for _, p := range u.payments {
			count++
			diff := centsConv(float64(p.amount)) - mean
			squaredDiffs += diff * diff
		}
	}
	return math.Sqrt(squaredDiffs / float64(count))
}

// LoadData does the thing
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
		if len(line) > 3 {
			id, _ := strconv.Atoi(line[0])  // No longer uses bounds-checking
			age, _ := strconv.Atoi(line[2]) // No longer uses bounds-checking
			users[UserID(id)] = &User{uint(age), UserID(id), []Payment{}}

		}
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
		if len(line) == 3 {
			userID, _ := strconv.Atoi(line[2]) // No longer uses bounds-checking
			paymentCents, _ := strconv.Atoi(line[0])
			users[UserID(userID)].payments = append(users[UserID(userID)].payments, Payment{
				uint32(paymentCents)})
		}
	}

	return users
}

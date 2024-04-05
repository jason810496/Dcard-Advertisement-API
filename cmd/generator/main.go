package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/config"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/database"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/services"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
)

var Amount int

func init() {
	flag.IntVar(&Amount, "n", 100, "Amount of fake data")
	flag.Usage = func() {
		fmt.Println("Usage: go run main.go -n <amount> -config <config mode>")
		flag.PrintDefaults()
	}
}

func main() {
	config.Init()
	flag.Parse()

	database.Init()
	database.CheckConnection()
	srv := services.NewAdminService()

	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)

	// generate fake data
	for i := 0; i < Amount; i++ {
		req := randAd(r, i)
		err := srv.CreateAdvertisement(&req)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("CreateAdvertisement", req)
	}
}

func randAd(r *rand.Rand, i int) schemas.CreateAdRequest {
	req := schemas.NewCreateAdRequest()
	req.Title = fmt.Sprintf("AD-%d", i)
	st, ed := randAge(r)
	req.Conditions.AgeStart = st
	req.Conditions.AgeEnd = ed
	req.Conditions.Gender = randPick(r, randLen(r, 2), utils.GenderList)
	req.Conditions.Country = randPick(r, randLen(r, 5), utils.CountryList)
	req.Conditions.Platform = randPick(r, randLen(r, 3), utils.PlatformList)
	req.StartAt = time.Now()
	req.EndAt = time.Now().Add(time.Hour * 24 * 30)
	return req
}

func randAge(r *rand.Rand) (int, int) {
	a := r.Intn(100) + 1
	b := r.Intn(100) + 1
	if a > b {
		return b, a
	}
	return a, b
}

// return 1 <= return val <= n
func randLen(r *rand.Rand, n int) int {
	return r.Intn(n + 1)
}

// random pick k items from arr
func randPick(r *rand.Rand, k int, arr []string) []string {
	if k == 0 {
		return []string{}
	}

	n := len(arr)
	indices := r.Perm(n)
	ret := make([]string, k)
	for i := 0; i < k; i++ {
		ret[i] = arr[indices[i]]
	}
	return ret
}

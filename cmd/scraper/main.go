package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/anantadwi13/jadwal-sholat-scraper/pkg/scraper"
)

var lock sync.Mutex
var allProvinsi []*scraper.ResProvinsi

func main() {
	var err error
	wg := &sync.WaitGroup{}

	scraper.Init()
	allProvinsi, err = scraper.ScrapeProvinsi()

	if err != nil {
		panic(err)
	}

	for _, provinsi := range allProvinsi {
		wg.Add(1)
		go scrapeKotaKab(wg, provinsi, 1)
	}

	fmt.Println("Waiting ...")

	wg.Wait()

	dump()
}

func scrapeKotaKab(wg *sync.WaitGroup, provinsi *scraper.ResProvinsi, current int) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recover ", r)
		}
	}()

	kotaKabs, err := scraper.ScrapeKotaKab(provinsi)

	if err != nil {
		if current > 5 {
			panic(err)
		}

		time.Sleep(time.Second)
		wg.Add(1)

		scrapeKotaKab(wg, provinsi, current+1)
		return
	}

	for _, kotaKab := range kotaKabs {
		wg.Add(1)
		go scrapeJadwalSholat(wg, kotaKab, 1)
	}
}

func scrapeJadwalSholat(wg *sync.WaitGroup, kotaKab *scraper.ResKotaKab, current int) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recover ", r)
		}
	}()

	if _, err := scraper.ScrapeJadwalSholat(kotaKab, 3, 2021); err != nil {
		if current > 5 {
			panic(err)
		}

		time.Sleep(time.Second)
		wg.Add(1)

		scrapeJadwalSholat(wg, kotaKab, current+1)
	} else {
		dump()
	}
}

func dump() {
	lock.Lock()
	defer lock.Unlock()

	if data, err := json.MarshalIndent(allProvinsi, "", "  "); err == nil {
		if err := os.WriteFile("data.json", data, 0666); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error ", err)
	}
}

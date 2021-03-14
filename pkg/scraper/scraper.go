package scraper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	http_client "github.com/anantadwi13/jadwal-sholat-scraper/pkg/http-client"
)

type ResProvinsi struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	KotaKab []*ResKotaKab `json:"kotakab"`
}

type ResKotaKab struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	Data     *ResDataKotaKab `json:"data"`
	Provinsi *ResProvinsi    `json:"-"`
}

type ResDataKotaKab struct {
	Status  int                        `json:"status"`
	Message string                     `json:"message"`
	Bujur   string                     `json:"bujur"`
	Lintang string                     `json:"lintang"`
	KabKo   string                     `json:"kabko"`
	Prov    string                     `json:"prov"`
	Data    map[string]ResJadwalSholat `json:"data"`
	KotaKab *ResKotaKab                `json:"-"`
}

type ResJadwalSholat struct {
	Tanggal string `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Terbit  string `json:"terbit"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

func Init() {
	client := http_client.HttpClientInstance()

	client.Get("https://bimasislam.kemenag.go.id/jadwalshalat")
}

func ScrapeProvinsi() ([]*ResProvinsi, error) {
	client := http_client.HttpClientInstance()

	res, err := client.Get("https://bimasislam.kemenag.go.id/jadwalshalat")

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		return nil, err
	}

	provinsiElements := doc.Find("#search_prov").First().Children()

	provinsi := make([]*ResProvinsi, 0)

	for i := range provinsiElements.Nodes {
		provEl := provinsiElements.Eq(i)
		tempProv := &ResProvinsi{Id: provEl.AttrOr("value", ""), Name: provEl.Text()}
		provinsi = append(provinsi, tempProv)
	}

	return provinsi, nil
}

func ScrapeKotaKab(provinsi *ResProvinsi) ([]*ResKotaKab, error) {
	client := http_client.HttpClientInstance()
	res, err := client.Post("https://bimasislam.kemenag.go.id/ajax/getKabkoshalat", http_client.FormData{
		"x": provinsi.Id,
	})

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))

	if err != nil {
		return nil, err
	}

	kotaKabElements := doc.Find("option")

	kotaKab := make([]*ResKotaKab, 0)

	for i := range kotaKabElements.Nodes {
		kotaKabEl := kotaKabElements.Eq(i)
		tempKotaKab := &ResKotaKab{Id: kotaKabEl.AttrOr("value", ""), Name: kotaKabEl.Text(), Provinsi: provinsi}
		kotaKab = append(kotaKab, tempKotaKab)
	}

	provinsi.KotaKab = kotaKab

	return kotaKab, nil
}

func ScrapeJadwalSholat(kotaKab *ResKotaKab, month int, year int) (*ResDataKotaKab, error) {
	client := http_client.HttpClientInstance()

	fmt.Println("Running " + kotaKab.Name)

	res, err := client.Post("https://bimasislam.kemenag.go.id/ajax/getShalatbln", http_client.FormData{
		"x":   kotaKab.Provinsi.Id,
		"y":   kotaKab.Id,
		"bln": strconv.Itoa(month),
		"thn": strconv.Itoa(year),
	})

	if err != nil {
		return nil, err
	}

	data := &ResDataKotaKab{KotaKab: kotaKab}

	err = json.Unmarshal([]byte(res), data)

	if err != nil {
		return nil, err
	}

	fmt.Println("Done "+data.Prov+" - "+data.KabKo, data.KotaKab.Name)

	kotaKab.Data = data

	return data, nil
}

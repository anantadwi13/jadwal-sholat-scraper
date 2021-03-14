# Jadwal Sholat Scraper

## About
Web scraping jadwal sholat yang ditulis dalam bahasa Go. Memanfaatkan API https://bimasislam.kemenag.go.id/jadwalshalat yang nantinya disimpan dalam file JSON.

## Usage
- Running from source code
  ```shell
  make run
  # or
  go run cmd/scraper/main.go
  ```

- Compiling source code & running compiled version
  ```shell
  # Compile
  make build
  
  # Run
  ./build/scrapper
  
  # Read output
  cat ./data.json
  ```
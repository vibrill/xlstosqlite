package Proses

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

func getTitle(fileSisPtk string) (judul []string) {
	path := "/Temp" //set path
	f, err := excelize.OpenFile(path + "/" + fileSisPtk)
	if err == nil {
		f.SetActiveSheet(0)
		sheetName := f.GetSheetName(0)
		//rows, err := f.GetRows(sheetName)
		rows, err := f.Rows(sheetName)
		if err != nil {
			fmt.Println(err, "kedua")
		}
		rows.Next()
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err, "ketiga")
		}
		judul = append(judul, row...)
	}
	/*for _, colCell := range row {
		//fmt.Println(colCell)
		judul = append(judul, colCell)
	}*/
	return judul
}

//func getFiles() files[]

func Proses(db *sql.DB) {
	_, err := os.Stat("/Temp")
	if err == nil {
		path := "/Temp/"
		var guru, tendik, siswa string
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println(err, "keempat")
		}
		for _, item := range files {
			a := item.Name()
			if len(a) > 10 {
				switch a[:9] {
				case "daftar_gu":
					guru = a
				case "daftar-gu":
					guru = a
				case "daftar-te":
					tendik = a
				case "daftar_te":
					tendik = a
				case "daftar-pd":
					siswa = a
				case "daftar_pd":
					siswa = a
				}
			}
		}
		//fmt.Println(guru, tendik, siswa)
		judulTabelPTK := getTitle(guru)
		judulTabelSiswa := getTitle(siswa)
		if len(judulTabelPTK) != 0 {
			createTable(db, "PTK", judulTabelPTK) // Create Database Tables
			prosesPTK(db, guru)
			prosesPTK(db, tendik)
			cleanPTK(db)
			os.Remove(path + guru)
			os.Remove(path + tendik)
		}
		if len(judulTabelSiswa) != 0 {
			createTable(db, "SISWA", judulTabelSiswa)
			prosesSiswa(db, siswa)
			os.Remove(path + siswa)
		}
	}
}

func cleanPTK(db *sql.DB) {
	var totalROW, data, stt string
	var max string
	/* db, _ := sql.Open("sqlite3", "./"+namadb) // Open the created SQLite File
	defer db.Close()                          // Defer Closing the database */
	//query, err := db.Prepare("select count(notebook) from pages where notebook = ?")
	query, _ := db.Query("SELECT COUNT(*) FROM PTK")
	for query.Next() {
		query.Scan(&totalROW)
	}
	x, _ := strconv.Atoi(totalROW)
	for i := 0; i <= x-1; i++ {
		query, _ := db.Query("SELECT max(NO) FROM PTK")
		for query.Next() {
			query.Scan(&max)
		}
		query, _ = db.Query("SELECT NO FROM PTK where NO = " + (strconv.Itoa(i)))
		for query.Next() {
			query.Scan(&data)
		}
		stt = "UPDATE PTK SET NO = " + strconv.Itoa(i+1) + " WHERE ROWID = " + strconv.Itoa(i+1)
		statement, err := db.Prepare(stt)
		if err != nil {
			log.Fatal("kedua", err.Error())
		}
		statement.Exec() // Execute SQL Statements
		//fmt.Println("row " + strconv.Itoa(i+1) + " cleaned")
	}
	//fmt.Println("Tabel telah dirapikan dan siap digunakan")
}

func CreateDB(namadb string) {
	//createdb
	os.Remove("./" + namadb)              // I delete the file to avoid duplicated records.
	file, err := os.Create("./" + namadb) // Create SQLite file
	if err != nil {
		log.Fatal("ketiga" + err.Error())
	}
	file.Close()
	//fmt.Println(namadb + " telah dibuat")
	// SQLite is a file based database.
}

func createTable(db *sql.DB, SisPtk string, judul []string) {
	/* 	db, _ := sql.Open("sqlite3", "./"+namadb) // Open the created SQLite File
	   	defer db.Close()                          // Defer Closing the database */

	text := ""
	for _, item := range judul {
		text = text + `"` + item + `" ` + "TEXT,"
	}
	text = text[:len(text)-1]
	text = "CREATE TABLE " + SisPtk + " ( " + text + `);`
	createTableSQL := text

	//fmt.Println(text)
	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal("keempat" + err.Error())
	}
	statement.Exec() // Execute SQL Statements
	//fmt.Println("Tabel " + SisPtk + " telah dibuat")
}

func prosesPTK(db *sql.DB, ptk string) {
	path := "/Temp"
	f, err := excelize.OpenFile(path + "/" + ptk)
	if err != nil {
		log.Fatal("ERROR kelima", err.Error())
	}
	f.SetActiveSheet(0)
	sheetName := f.GetSheetName(0)
	//rows, err := f.GetRows(sheetName)
	rows, err := f.Rows(sheetName)
	if err != nil {
		fmt.Println(err, "kelima")
		return
	}

	rows.Next()
	_, err = rows.Columns()
	if err != nil {
		fmt.Println(err, "keenam")
	}

	var dataX [52]string //52 untuk ptk
	x := 0
	y := 0
	yd := 0
	//passing tittle
	for i := 0; i < 20; i++ {
		rows.Next()
		row, err := rows.Columns()
		//fmt.Println(row)
		if err != nil {
			fmt.Println(err, "ketujuh")
			return
		}

		x = 0
		for _, colCell := range row {
			//fmt.Println(colCell)
			dataX[x] = colCell
			x++
		}
		yd, err = strconv.Atoi(dataX[0])
		if err != nil {
			fmt.Println(err, "kedelapan")
			return
		}
		if y == yd {
			break
		} else {
			y, err = strconv.Atoi(dataX[0])
			if err != nil {
				fmt.Println(err, "kesembilan")
				return
			}

		}
		insertDataPTK(db, dataX)
	}
}

func prosesSiswa(db *sql.DB, siswa string) {
	path := "/Temp"
	f, err := excelize.OpenFile(path + "/" + siswa)
	if err != nil {
		log.Fatal("ERROR keenam", err.Error())
	}
	f.SetActiveSheet(0)
	sheetName := f.GetSheetName(0)
	//rows, err := f.GetRows(sheetName)
	rows, err := f.Rows(sheetName)
	if err != nil {
		fmt.Println(err, "kesepuluh")
		return
	}

	rows.Next()
	_, err = rows.Columns()
	if err != nil {
		fmt.Println(err, "kesebelas")
	}

	var dataX [66]string
	x := 0
	y := 0
	yd := 0
	//passing tittle
	for i := 0; i < 5000; i++ {
		rows.Next()
		row, err := rows.Columns()
		//fmt.Println(row)
		if err != nil {
			fmt.Println(err, "keduabelas")
			return
		}

		x = 0
		for _, colCell := range row {
			//fmt.Println(colCell)
			dataX[x] = colCell
			x++
		}
		yd, err = strconv.Atoi(dataX[0])
		if err != nil {
			fmt.Println(err, "ketigabelas")
			return
		}
		if y == yd {
			break
		} else {
			y, err = strconv.Atoi(dataX[0])
			if err != nil {
				fmt.Println(err, "keempatbelas")
				return
			}

		}
		insertDataSISWA(db, dataX)
	}
}

// We are passing db reference connection from main to our method with other parameters
func insertDataPTK(db *sql.DB, dataX [52]string) {
	/* db, _ := sql.Open("sqlite3", "./"+namadb) // Open the created SQLite File
	defer db.Close()                          // Defer Closing the database */
	//fmt.Println("Inserting data record ...")
	text := ""
	for _, item := range dataX {
		text = text + `"` + item + `"` + ","
	}
	text = text[:len(text)-1] //omit last coma
	//fmt.Println(text)

	insertSQL := "INSERT INTO PTK VALUES(" + text + ")"
	statement, err := db.Prepare(insertSQL) // Prepare statement.
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func insertDataSISWA(db *sql.DB, dataX [66]string) {
	/* db, _ := sql.Open("sqlite3", "./"+namadb) // Open the created SQLite File
	defer db.Close()                          // Defer Closing the database */
	//fmt.Println("Inserting data record ...")
	text := ""
	for _, item := range dataX {
		text = text + `"` + item + `"` + ","
	}
	text = text[:len(text)-1] //omit last coma
	//fmt.Println(text)

	insertSQL := "INSERT INTO SISWA VALUES(" + text + ")"
	statement, err := db.Prepare(insertSQL) // Prepare statement.
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	//iexcelize start
	u, _ := os.UserHomeDir()
	//mendapatkan path desktop
	_, err := ioutil.ReadDir(string(u) + `/Desktop`)
	ud := u + `/Desktop`
	if err != nil {
		_, err = ioutil.ReadDir(`E:/Desktop`)
		ud = `E:/Desktop`
		if err != nil {
			_, err = ioutil.ReadDir(`D:/Desktop`)
			ud = `D:/Desktop`
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	path := ud + "/DapoSniff"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	a := files[0].Name()
	f, err := excelize.OpenFile(path + "/" + a)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}
	f.SetActiveSheet(0)
	sheetName := f.GetSheetName(0)
	//rows, err := f.GetRows(sheetName)
	rows, err := f.Rows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows.Next()
	row, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return
	}
	var judul []string
	for _, colCell := range row {
		fmt.Println(colCell)
		judul = append(judul, colCell)
	}
	fmt.Println(judul)

	/*for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println("\n\n")
	}
	*/

	//database start
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase, judul)                               // Create Database Tables
	var dataX [52]string
	x := 0
	y := 0
	yd := 0
	for i := 0; i < 20; i++ {
		rows.Next()
		row, err := rows.Columns()
		if err != nil {
			fmt.Println(err)
			return
		}

		x = 0
		for _, colCell := range row {
			fmt.Println(colCell)
			dataX[x] = colCell
			x++
		}
		yd, err = strconv.Atoi(dataX[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		if y == yd {
			break
		} else {
			y, err = strconv.Atoi(dataX[0])
			if err != nil {
				fmt.Println(err)
				return
			}

		}
		insertData(sqliteDatabase, dataX)
	}
}

func createTable(db *sql.DB, judul []string) {
	text := ""
	for _, item := range judul {
		text = text + `"` + item + `" ` + "TEXT,"
	}
	text = text[:len(text)-1]
	text = "CREATE TABLE PTK ( " + text + `);`
	createTableSQL := text

	log.Println("Create PTK table...")
	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertData(db *sql.DB, judul [52]string) {
	log.Println("Inserting data record ...")
	text := ""
	for _, item := range judul {
		text = text + `"` + item + `"` + ","
	}
	text = text[:len(text)-1] //omit last coma
	fmt.Println(text)

	insertSQL := "INSERT INTO PTK VALUES(" + text + ")"
	fmt.Println("2")
	statement, err := db.Prepare(insertSQL) // Prepare statement.
	fmt.Println("3")
	fmt.Println("3")
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println("4")
		log.Fatalln(err.Error())

	}
	fmt.Println("8")
	_, err = statement.Exec()
	fmt.Println("7")
	if err != nil {
		log.Fatalln(err.Error())
		fmt.Println("5")
	}
	fmt.Println("6")
}

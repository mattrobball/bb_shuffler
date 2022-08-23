package main

import (
	// "encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mattrobball/bb_shuffler/local"
)

func filePath() (path string) {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Enter the csv file name.")
		log.Fatal("No file entered")
	}
	arg := args[0]

	if !strings.Contains(arg, "csv") {
		fmt.Println("This doesn't look like a csv file.")
		log.Fatal("Not a csv")
	}

	if arg[len(arg)-3:] != "csv" {
		fmt.Println("This doesn't look like a csv file.")
		log.Fatal("Not a csv")
	}

	_, err := os.Open(arg)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Can't find that file.")
		log.Fatal(err)
	}

	path, _ = filepath.Abs(arg)
	return
}

func check(err error) {
	if err != nil {
		log.Fatal("A generic error")
	}
}

func main() {
	
	path := filePath()

	file, err := os.Open(path) 
	check(err)

	buf := new(strings.Builder)
	_, err = io.Copy(buf,file)
	check(err)

	s := strings.Replace(buf.String(),"\"","",-1)
	t := strings.Split(s,"\n")
	t = t[1:]

	var records [][]string
	for _, record := range(t[:len(t)-1]) {
		broken_record := strings.Split(record,",")
		records = append(records,broken_record[:4])
	}

	var email_list []string 
	for _, student := range(records) {
		pieces := []string{student[2],"@email.sc.edu"}
		email := strings.Join(pieces,"")
		email_list = append(email_list, email)
	}

	mailing_list := strings.Join(email_list,", ")

	fmt.Printf("%s\n",mailing_list)

	var class_list []string
	for _, student := range(records) {
		var c cases.Caser = cases.Title(language.English)
		first_name := c.String(student[1])
		names := []string{first_name,student[0]}
		name := strings.Join(names," ")
		class_list = append(class_list, local.ChosenName(name))
	}

	fmt.Println(strings.Join(class_list,"\n"))

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(class_list), func(i, j int) { class_list[i], class_list[j] = class_list[j], class_list[i] })

	// class_list = append(class_list, "Some Dude")

	var total, offset int
	remainder := len(class_list) % 4 
	if remainder != 0 {
		offset = 1
	} 
	
	total = (len(class_list) / 4) + offset

	groups := make([]string,total)

	// fmt.Println(class_list)

	for j, name := range class_list {
		l := j % total
		if groups[l] == "" {
			groups[l] = name
		} else {
			groups[l] = groups[l] + ", " + name 
		}
	}

	for i := 0; i < total ; i++ {
		fmt.Printf("Group %s : %s \n",strconv.Itoa(i+1),groups[i])
	}

}

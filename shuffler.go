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
	"sort"
	"strings"
	"time"
	"strconv"

	"github.com/mattrobball/bb_shuffler/local"
)

type studentInfo struct {
	firstName string 
	lastName string 
	email string 
}

type infoList []studentInfo

func (students infoList) Len() int { return len(students) }

func (students infoList) Swap(i, j int) { students[i], students[j] = students[j], students[i] }

func (students infoList) Less(i, j int) bool { 
    if students[i].lastName != students[j].lastName {
        return students[i].lastName < students[j].lastName
    }
    return students[i].firstName < students[j].firstName 
}

func (students infoList) emailList() string {
	var emails string 
	for j, val := range(students) {
		if j == 0 {
			emails = val.email + "@email.sc.edu" 
		} else {
			emails = emails + ", " + val.email + "@email.sc.edu"
		}
	}
	return emails 
}

func (student studentInfo) name() string {
	name := student.firstName + " " +  student.lastName 
	name = local.ChosenName(name)
	return name 
}

func (students infoList) classList(sep string) string {
	sort.Sort(students)
	var classList string 
	for j, val := range(students) {
		if j != 0 {
			classList = classList + sep + val.name()
		} else {
			classList = classList + val.name()
		}
	}
	return classList
}

func (students infoList) groups(n int) []infoList {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(students), func(i, j int) { students[i], students[j] = students[j], students[i] })

	var total, offset int
	remainder := len(students) % n 
	if remainder != 0 {
		offset = 1
	} 
	
	total = (len(students) / n) + offset

	groups := make([]infoList,total)

	for j, student := range students {
		l := j % total
		groups[l] = append(groups[l],student)
	}

	return groups 
}

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

	var records infoList
	for _, record := range(t[:len(t)-1]) {
		broken_record := strings.Split(record,",")
		first := strings.Title(strings.ToLower(broken_record[1]))
		last := strings.Title(strings.ToLower(broken_record[0]))
		student := studentInfo{first,last,broken_record[2]}
		records = append(records,student)
	}

	mailingList := records.emailList()

	fmt.Println("Email List")
	fmt.Println("-----")
	fmt.Printf("%s\n\n",mailingList)

	classList := records.classList("\n")

	fmt.Println("Class List")
	fmt.Println("-----")
	fmt.Printf("%s\n\n",classList)

	groups := records.groups(4)
	total := len(groups)

	fmt.Println("Groups")
	fmt.Println("-----")
	for i := 0; i < total ; i++ {
		fmt.Printf("Group %s : %s \n",strconv.Itoa(i+1),groups[i].classList(", "))
	}

}

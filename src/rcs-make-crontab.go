package main

import (
    "fmt"
    "encoding/json"
//    "time"
    "io/ioutil"
//    "crypto/md5"
    "log"
    "os"
//    "bufio"
//    "strconv"
//    "reflect"
)

var conf_dir = "/var/lib/rcs/"
var works_dir = conf_dir + "works"
var archiv_dir = conf_dir + "archiv"
var run_dir = conf_dir + "run"

type Work struct {
    Id           string
    Title        string
    Path         string
    TimeTable    string
    MaxSnap      string
    Services     string
    Status       string
    Done         bool
}

func LoadWork(id string) Work {
    var w Work
    fname := fmt.Sprintf("%s/%s", works_dir, id)
    content, err := ioutil.ReadFile(fname)
    if err != nil { /*log.Fatal("Error when opening file: ", err)*/ return w }
    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil { /*log.Fatal("Error during Unmarshal(): ", err)*/ return w }
//    fmt.Println(reflect.TypeOf(data["Title"]), data["Title"])
    w.Id = id
    w.Title = fmt.Sprintf("%s", data["Title"])
    w.Path = fmt.Sprintf("%s", data["Path"])
    w.TimeTable = fmt.Sprintf("%s", data["TimeTable"])
    w.MaxSnap = fmt.Sprintf("%s", data["MaxSnap"])
    w.Services = fmt.Sprintf("%s", data["Services"])
    w.Status = fmt.Sprintf("%s", data["Status"])
    w.Done = false
    return w
}

func GetWorkList(works []Work) []Work {
    files, err := ioutil.ReadDir(works_dir)
    if err != nil { /*log.Fatal(err)*/ return works }
    for _, file := range files {
	w := LoadWork(file.Name())
	works = append(works, w)
	fmt.Printf("%s run-backup-zfs.sh  %s > /dev/null\n", w.TimeTable, w.Id)
    }
    return works
}

func main() {
    var WorkList []Work
    file, err := os.OpenFile(conf_dir + "rcs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }
    fmt.Printf("* * * * * /usr/bin/rcs-make-crontab > /var/spool/cron/crontabs/root\n")
    _ = GetWorkList(WorkList)
    log.Println("Обновление crontab")
}

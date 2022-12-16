package main

import (
    "fmt"
    "encoding/json"
    "time"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "strings"
)

func msys3(comstr string) string {
    arg := strings.Split(comstr, " ")
    out, err := exec.Command(arg[0], arg[1:]...).Output()
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}

var conf_dir = "/var/lib/rcs/"
var works_dir = conf_dir + "works"
var archiv_dir = conf_dir + "archiv"
var run_dir = conf_dir + "run"

type Conf struct {
    Title    	string
    Pool      	string
}

type Work struct {
    Id           string
    Title        string
    Path         string
    Login        string
    Password     string
    TimeTable    string
    MaxSnap      string
    Services     string
    Status       string
    Done         bool
}

func GetTime() string {
// utime := int32(time.Now().Unix())
   current_time := time.Now()
   return fmt.Sprintln(current_time.Format(time.UnixDate))
}

func LoadConf() Conf {
    var w Conf
    fname := fmt.Sprintf("%s%s", conf_dir, "rcs.conf")
    content, err := ioutil.ReadFile(fname)
    if err != nil { 
    	w.Title = fmt.Sprintf("%s", "")
    	w.Pool = fmt.Sprintf("%s", "")
	/*log.Fatal("Error when opening file: ", err)*/ 
	return w 
    }
    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil { /*log.Fatal("Error during Unmarshal(): ", err)*/ return w }
    w.Title = fmt.Sprintf("%s", data["Title"])
    w.Pool = fmt.Sprintf("%s", data["Pool"])
    return w
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
    w.Login = fmt.Sprintf("%s", data["Login"])
    w.Password = fmt.Sprintf("%s", data["Password"])
    w.TimeTable = fmt.Sprintf("%s", data["TimeTable"])
    w.MaxSnap = fmt.Sprintf("%s", data["MaxSnap"])
    w.Services = fmt.Sprintf("%s", data["Services"])
    w.Status = fmt.Sprintf("%s", data["Status"])
    w.Done = false
    return w
}

func GetZfsInfo(Pool, Id string) {
    com:=fmt.Sprintf("zfs list -H %s/%s", Pool, Id)
    rs := msys3(com)
    r := strings.Fields(strings.TrimSpace(rs))
    fmt.Println(r[1], r[2])
}

func GetZpoolInfo(Pool string) {
    com:=fmt.Sprintf("zfs list -H %s", Pool)
    rs := msys3(com)
    r := strings.Fields(strings.TrimSpace(rs))
    fmt.Println(r[1], r[2])
}

func GetWorkList(works []Work) []Work {
    conf := LoadConf()
    files, err := ioutil.ReadDir(works_dir)
    if err != nil { /*log.Fatal(err)*/ return works }
    fmt.Printf("%s\n", conf.Pool)
    GetZpoolInfo(conf.Pool)
    for _, file := range files {
	w := LoadWork(file.Name())
	works = append(works, w)
	fmt.Printf("%s/%s\n", conf.Pool, w.Id)
	GetZfsInfo(conf.Pool, w.Id)
    }
    return works
}

func main() {
    var WorkList []Work
    file, err := os.OpenFile(conf_dir + "rcs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }
    _ = GetWorkList(WorkList)
    log.Println("Обновление crontab")
}


package main

import (
    "fmt"
    "html/template"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "time"
    "crypto/md5"
    "log"
    "os"
    "bufio"
//    "strconv"
//    "reflect"
)

var conf_dir = ".rcs/"
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

type PageData struct {
    PageTitle    string
    Message      string
    WorkList     []Work
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
    }
    return works
}

func GetLog() string {
    Log := ""
    file, err := os.Open(".rcs/rcs-server.log")
    if err != nil { log.Fatal(err) }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        Log = Log + fmt.Sprintf("%s\n", scanner.Text())
    }
    if err := scanner.Err(); err != nil { log.Fatal(err) }
    return Log
}

func main() {
    file, err := os.OpenFile(conf_dir + "rcs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "Список ресурсов:",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles("form_list.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
//	log.Println("Просмотр списка ресурсов")
    })


    http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles("form_view.html"))
        if r.Method != http.MethodGet {
            tmpl_view.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	work := LoadWork(Id)
        tmpl_view.Execute(w, work)
//	if work.Status=="prepared" {
//	    n, _ := strconv.Atoi(work.MaxSnap)
//	    fmt.Fprintf(w, "<ul>")
//	    for i := 1; i<(n+1); i++ {
//		fmt.Fprintf(w, "<li><a href=\"/get-work-file?id=%s&file=v%d.xyz\" target=\"_blank\">v%d.xyz</a></li>", Id, i, i)
//	    }
//	    fmt.Fprintf(w, "</ul>")
//	}
	log.Println("Просмотр ресурса: ", Id)
    })


    http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "Журнал:",
            Message: "",
        }
	tmpl_list := template.Must(template.ParseFiles("form_log.html"))
	data.Message = GetLog()
        tmpl_list.Execute(w, data)
//	log.Println("Просмотр журнала")
    })


    http.HandleFunc("/set-status", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles("form_view.html"))
        if r.Method != http.MethodGet {
            tmpl_view.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	Status := r.FormValue("status")
	work := LoadWork(Id)
	work.Status = Status
        dat, err := json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s/%s", works_dir, Id)
        _ = ioutil.WriteFile(fname, dat, 0644)
        tmpl_view.Execute(w, work)
	log.Println("Set status work", Id, "to", Status)
    })


    http.HandleFunc("/get-work-file", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            return
        }
	Id := r.FormValue("id")
	wFile := r.FormValue("file")
	fname := fmt.Sprintf("%s/%s/%s", run_dir, Id, wFile)
	content, _ := ioutil.ReadFile(fname)
	fmt.Printf("file %s content:\n %s", fname, content)
	fmt.Fprintf(w, "%s", content)
	log.Println("Get file", wFile, " for", Id)
    })


    http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
	tmpl_auto := template.Must(template.ParseFiles("form_auto.html"))
        if r.Method != http.MethodGet {
            tmpl_auto.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	p1 := fmt.Sprintf("%s/%s", works_dir, Id)
	p2 := fmt.Sprintf("%s/%s", archiv_dir, Id)
	e := os.Rename(p1, p2)
	if e != nil { /* log.Fatal(e) */ }
	msg := fmt.Sprintf("Удаление ресурса: %s", Id)
        data := PageData { Message: msg, }
        tmpl_auto.Execute(w, data)
	log.Println("Удаление ресурса: ", Id)
    })


    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
	tmpl_add := template.Must(template.ParseFiles("form_add.html"))
        if r.Method != http.MethodPost {
            tmpl_add.Execute(w, nil)
            return
        }
        work := Work{
            Title:      r.FormValue("Title"),
            Path:       r.FormValue("Path"),
            TimeTable:  r.FormValue("TimeTable"),
            MaxSnap:    r.FormValue("MaxSnap"),
            Services:   r.FormValue("Services"),
            Status:     r.FormValue("Status"),
        }
        dat, err := json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
        utime := int32(time.Now().Unix())
	hmd5 := md5.Sum([]byte(dat))
	work.Id = fmt.Sprintf("%d-%x", utime, hmd5)
        dat, err = json.MarshalIndent(work, "", " ")
        if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s/%s", works_dir, work.Id)
        _ = ioutil.WriteFile(fname, dat, 0644)
//        fmt.Println(fname, string(dat))
        _ = work
        tmpl_add.Execute(w, struct{ Success bool }{true})
	log.Println("Добавление нового ресурса: ", work.Id)
    })


    http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
	    tmpl_edit := template.Must(template.ParseFiles("form_edit.html"))
	    Id := r.FormValue("id")
	    work := LoadWork(Id)
	    fmt.Println(work)
            tmpl_edit.Execute(w, work)
            return
        }
	tmpl_edit_save := template.Must(template.ParseFiles("form_edit_save.html"))
	work := Work{
	    Id :        r.FormValue("id"),
	    Title:      r.FormValue("Title"),
	    Path:       r.FormValue("Path"),
	    TimeTable:  r.FormValue("TimeTable"),
	    MaxSnap:    r.FormValue("MaxSnap"),
	    Services:   r.FormValue("Services"),
	    Status:     r.FormValue("Status"),
	}
	dat, err := json.MarshalIndent(work, "", " ")
	if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s/%s", works_dir, work.Id)
	_ = ioutil.WriteFile(fname, dat, 0644)
	_ = work
	fmt.Println(work)
	tmpl_edit_save.Execute(w, work)
	log.Println("Редактирование ресурса: ", work.Id)
//	fmt.Println(dat)
    })


    http.ListenAndServe(":8080", nil)
}


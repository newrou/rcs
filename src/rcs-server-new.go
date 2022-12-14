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
    "syscall"
//    "strconv"
//    "reflect"
)


var conf_dir = "rcs/"
//var conf_dir = "/var/lib/rcs/"
var works_dir = conf_dir + "works/"
var archiv_dir = conf_dir + "archiv/"
var run_dir = conf_dir + "run/"
var www_dir = conf_dir + "www-new/"

type Conf struct {
    Title    	string
    Pool      	string
}


type Work struct {
    Id          string
    Title       string
    Path        string
    Login       string
    Password    string
    TimeTable   string
    MaxSnap     string
    Services    string
    Status      string
    Done        bool
}


type PageData struct {
    PageTitle   string
    Message     string
    WorkList    []Work
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
    fname := fmt.Sprintf("%s%s", works_dir, id)
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
    file, err := os.Open(conf_dir + "rcs-server.log")
    if err != nil { log.Fatal(err) }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        Log = Log + fmt.Sprintf("%s\n", scanner.Text())
    }
    if err := scanner.Err(); err != nil { log.Fatal(err) }
    return Log
}


func GetTime() string {
//   utime := int32(time.Now().Unix())
   current_time := time.Now()
   return fmt.Sprintln(current_time.Format(time.UnixDate))
}


func GetState() string {
    sysinfo := syscall.Sysinfo_t{}
    err := syscall.Sysinfo(&sysinfo)
    if err != nil { fmt.Println("Error:", err) } 
    Stat := fmt.Sprintln("Ok!\n", GetTime())
    Stat += fmt.Sprintln("Uptime:", sysinfo.Uptime)
    Stat += fmt.Sprintln("Loads:", sysinfo.Loads)
    Stat += fmt.Sprintln("Totalram:", sysinfo.Totalram)
    Stat += fmt.Sprintln("Freeram:", sysinfo.Freeram)
    Stat += fmt.Sprintln("Sharedram:", sysinfo.Sharedram)
    Stat += fmt.Sprintln("Bufferram:", sysinfo.Bufferram)
    Stat += fmt.Sprintln("Totalswap:", sysinfo.Totalswap)
    Stat += fmt.Sprintln("Freeswap:", sysinfo.Freeswap)
    Stat += fmt.Sprintln("Procs:", sysinfo.Procs)
    Stat += fmt.Sprintln("Pad:", sysinfo.Pad)
    Stat += fmt.Sprintln("Pad_cgo_0:", sysinfo.Pad_cgo_0)
    Stat += fmt.Sprintln("Totalhigh:", sysinfo.Totalhigh)
    Stat += fmt.Sprintln("Freehigh:", sysinfo.Freehigh)
    Stat += fmt.Sprintln("Unit:", sysinfo.Unit)
    Stat += fmt.Sprintln("X_f:", sysinfo.X_f)
    Stat += fmt.Sprintln("Pad_cgo_1:", sysinfo.Pad_cgo_1)
//    Log := fmt.Sprintln("Ok!\n\n", GetTime())    
    return Stat
}


func main() {
    file, err := os.OpenFile(conf_dir + "rcs-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    // if err != nil { log.Fatal(err) }
    if err == nil { log.SetOutput(file) }


    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "???????????? RCS:",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles(www_dir + "form_main.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
//	log.Println("???????????????? ???????????? ????????????????")
    })


    http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "???????????? RCS: ????????",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles(www_dir + "form_menu.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
//	log.Println("???????????????? ???????????? ????????????????")
    })


    http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "???????????? RCS: ???????????? ????????????????:",
            WorkList: []Work{},
        }
	tmpl_list := template.Must(template.ParseFiles(www_dir + "form_list.html"))
	data.WorkList = GetWorkList(data.WorkList)
        tmpl_list.Execute(w, data)
//	log.Println("???????????????? ???????????? ????????????????")
    })


    http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles(www_dir + "form_view.html"))
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
	log.Println("???????????????? ??????????????: ", Id)
    })


    http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "?????????????????? ??????????????:",
            Message: "",
        }
	tmpl_list := template.Must(template.ParseFiles(www_dir + "form_state.html"))
	data.Message = GetState()
        tmpl_list.Execute(w, data)
//	log.Println("???????????????? ??????????????")
    })

 
    http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            PageTitle: "????????????:",
            Message: "",
        }
	tmpl_list := template.Must(template.ParseFiles(www_dir + "form_log.html"))
	data.Message = GetLog()
        tmpl_list.Execute(w, data)
//	log.Println("???????????????? ??????????????")
    })


    http.HandleFunc("/set-status", func(w http.ResponseWriter, r *http.Request) {
	tmpl_view := template.Must(template.ParseFiles(www_dir + "form_view.html"))
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
	fname := fmt.Sprintf("%s%s", works_dir, Id)
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
	fname := fmt.Sprintf("%s%s/%s", run_dir, Id, wFile)
	content, _ := ioutil.ReadFile(fname)
	fmt.Printf("file %s content:\n %s", fname, content)
	fmt.Fprintf(w, "%s", content)
	log.Println("Get file", wFile, " for", Id)
    })


    http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
	tmpl_auto := template.Must(template.ParseFiles(www_dir + "form_auto.html"))
        if r.Method != http.MethodGet {
            tmpl_auto.Execute(w, nil)
            return
        }
	Id := r.FormValue("id")
	p1 := fmt.Sprintf("%s%s", works_dir, Id)
	p2 := fmt.Sprintf("%s%s", archiv_dir, Id)
	e := os.Rename(p1, p2)
	if e != nil { /* log.Fatal(e) */ }
	msg := fmt.Sprintf("???????????????? ??????????????: %s", Id)
        data := PageData { Message: msg, }
        tmpl_auto.Execute(w, data)
	log.Println("???????????????? ??????????????: ", Id)
    })


    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
	tmpl_add := template.Must(template.ParseFiles(www_dir + "form_add.html"))
        if r.Method != http.MethodPost {
            tmpl_add.Execute(w, nil)
            return
        }
        work := Work{
            Title:      r.FormValue("Title"),
            Path:       r.FormValue("Path"),
	    Login:      r.FormValue("Login"),
	    Password:   r.FormValue("Password"),
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
	fname := fmt.Sprintf("%s%s", works_dir, work.Id)
        _ = ioutil.WriteFile(fname, dat, 0644)
//        fmt.Println(fname, string(dat))
        _ = work
        tmpl_add.Execute(w, struct{ Success bool }{true})
	log.Println("???????????????????? ???????????? ??????????????: ", work.Id)
    })


    http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
	    tmpl_edit := template.Must(template.ParseFiles(www_dir + "form_edit.html"))
	    Id := r.FormValue("id")
	    work := LoadWork(Id)
	    fmt.Println(work)
            tmpl_edit.Execute(w, work)
            return
        }
	tmpl_edit_save := template.Must(template.ParseFiles(www_dir + "form_edit_save.html"))
	work := Work{
	    Id :        r.FormValue("id"),
	    Title:      r.FormValue("Title"),
	    Path:       r.FormValue("Path"),
	    Login:      r.FormValue("Login"),
	    Password:   r.FormValue("Password"),
	    TimeTable:  r.FormValue("TimeTable"),
	    MaxSnap:    r.FormValue("MaxSnap"),
	    Services:   r.FormValue("Services"),
	    Status:     r.FormValue("Status"),
	}
	dat, err := json.MarshalIndent(work, "", " ")
	if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s%s", works_dir, work.Id)
	_ = ioutil.WriteFile(fname, dat, 0644)
	_ = work
	fmt.Println(work)
	tmpl_edit_save.Execute(w, work)
	log.Println("???????????????????????????? ??????????????: ", work.Id)
//	fmt.Println(dat)
    })


    http.HandleFunc("/conf", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
	    tmpl_edit := template.Must(template.ParseFiles(www_dir + "form_conf.html"))
	    conf := LoadConf()
	    fmt.Println(conf)
            tmpl_edit.Execute(w, conf)
            return
        }
	tmpl_edit_save := template.Must(template.ParseFiles(www_dir + "form_conf_save.html"))
	conf := Conf{
	    Title:      r.FormValue("Title"),
	    Pool:       r.FormValue("Pool"),
	}
	dat, err := json.MarshalIndent(conf, "", " ")
	if err != nil { fmt.Println(err) }
	fname := fmt.Sprintf("%s%s", conf_dir, "rcs.conf")
	_ = ioutil.WriteFile(fname, dat, 0644)
	_ = conf
	fmt.Println(conf)
	tmpl_edit_save.Execute(w, conf)
	log.Println("??????????????????")
//	fmt.Println(dat)
    })


    fileServer := http.FileServer(http.Dir("rcs/snapshot"))
    http.Handle("/snapshot/", http.StripPrefix("/snapshot/", fileServer))


    http.ListenAndServe(":8080", nil)
}


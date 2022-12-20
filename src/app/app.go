package main

import (
    "html/template"
    "database/sql"
    "net/http"
    "runtime"
    "time"
    "fmt"
    "os"

    "github.com/berta00/cloudNoteApp/utils"
    _ "github.com/go-sql-driver/mysql"
)

// debug mode
var DebugMode bool = true;

// ws global info
var WSdomain string = "localhost";
var WSport   string = os.Getenv("WS_PORT");
var JWTsec   string = os.Getenv("JWT_SECRET");

// db global info
var DBname string = os.Getenv("MYSQL_DB");
var DBuser string = os.Getenv("MYSQL_USER");
var DBpass string = os.Getenv("MYSQL_PASSWORD");
var DBaddr string = os.Getenv("MYSQL_HOST");
var DBport string = os.Getenv("MYSQL_PORT");

// gmail access info
var GMemail string = "";
var GMpass string = os.Getenv("GMAIL_PASSWORD");

// users
var TotalUsers int = 0;

// login page struct
type loginPage struct {
    NewAcc    bool
    WrongCred bool
    ConfMail  bool
    Section   string
}
// dash page struct
type dashPage struct {
    Nick      string
    Name      string
    Email     string
    EmailCk   bool
    Section   string
}

// db table struct
type UsersStruct struct {
    id       int
    nick     string
    name     string
    email    string
    password string
    emailCk  bool
    admin    string
    date     string
}
type BasicNoteStruct struct {
    id       int
    creator  string
    content  string
    crDate   string
    mfDate   string
}
//tallini
//sesso4
func main(){
    if DebugMode {
        fmt.Println("Setting/updating env varaible:");
        utils.EnvVarSet();
    }
    fmt.Println("\nWebserver start program:");
    current_time := time.Now();
    fmt.Println("- " + current_time.Format("2006-01-02 15:04:05") + " Starting webserver at port: " + WSport + "...");

    // initialize routes
    routesInit();

    // start webserver
    err := http.ListenAndServe(":" + WSport, nil);
    if err == nil {
        fmt.Print("\n\nError (can't start webserver): ");
        fmt.Println(err);
    }
}

func routesInit(){
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

    http.HandleFunc("/",          mainRoute);
    http.HandleFunc("/dash",      dashRoute);
    http.HandleFunc("/reg",       regRoute);
    http.HandleFunc("/emailceck", emailCecked);
}

func mainRoute(w http.ResponseWriter, r *http.Request){
    // users number track
    TotalUsers++;
    current_time := time.Now();
    fmt.Printf("- " + current_time.Format("2006-01-02 15:04:05") + " Total user: %d\r", TotalUsers);
    // query parsing
    // 1: account created, 2: wrong credentials 3: email confirmed
    msg := r.URL.Query().Get("msg");
    sec := r.URL.Query().Get("sec");
    
    newLogin := new(loginPage);
    newLogin.Section = sec;

    newLogin.NewAcc = false;
    newLogin.WrongCred = false;
    newLogin.ConfMail = false;
    switch(msg){
        case "1":
            newLogin.NewAcc = true;
            break;
        case "2":
            newLogin.WrongCred = true;
            break;
        case "3":
            newLogin.ConfMail = true;
    }
    // html template
    Cwd, _ := os.Getwd();
    Os := runtime.GOOS
    switch Os {
        case "windows":
            template, _ := template.ParseFiles(Cwd + "\\static\\pages\\login.html")
            template.Execute(w, newLogin)
            break
        default:
            template, _ := template.ParseFiles(Cwd + "/static/pages/login.html")
            template.Execute(w, newLogin)
    }
}

func dashRoute(w http.ResponseWriter, r *http.Request){
    userIsGet := false;
    userIsAuth := false;
    userNick := "";
    userName := "";
    userEmail := "";
    userEmailCk := false;
    dashSection := "";

    switch(r.Method){
        case "GET":
            userIsGet = true;
            break;

        case "POST":
            email := r.FormValue("email");
            encodedPass := r.FormValue("password");

            dashSection = r.URL.Query().Get("sec");;

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, err := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();
            // query the users
            dbQuery := "select nick, name, email, password, emailConfirmed, admin from users where email='" + email + "';";
            userQuery, err := DBConn.Query(dbQuery);
            if err != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err);
            } 
            // read query
            userStructQuery := new(UsersStruct);
            for userQuery.Next(){
                userQuery.Scan(&userStructQuery.nick, &userStructQuery.name, &userStructQuery.email, &userStructQuery.password, &userStructQuery.emailCk, &userStructQuery.admin);
            }
            if encodedPass == userStructQuery.password {
                userIsAuth = true;
            }
            
            userNick = userStructQuery.nick;
            userName = userStructQuery.name;
            userEmail = userStructQuery.email;
            userEmailCk = userStructQuery.emailCk;

            break;
    }
    if userIsAuth {
        // re-connect to db
        dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
        DBConn, err := sql.Open("mysql", dbConnstring);
        if err != nil {
            fmt.Print("Authentication err: ");
            fmt.Println(err);
        }
        defer DBConn.Close();

        // query the users data


        dashData := new(dashPage);
        dashData.Nick = userNick;
        dashData.Name = userName;
        dashData.Email = userEmail;
        dashData.EmailCk = userEmailCk;
        dashData.Section = dashSection;

        // html template
        Cwd, _ := os.Getwd();
        Os := runtime.GOOS
        switch Os {
            case "windows":
                template, _ := template.ParseFiles(Cwd + "\\static\\pages\\dashboard.html")
                template.Execute(w, dashData)
                break
            default:
                template, _ := template.ParseFiles(Cwd + "/static/pages/dashboard.html")
                template.Execute(w, dashData)
        }
    } else if userIsGet {
        // redirect to home
        http.Redirect(w, r, "/", http.StatusFound);
    } else {
        // redirect to home |err: wrong credentials
        http.Redirect(w, r, "/?msg=2", http.StatusFound);
    }
}

func regRoute(w http.ResponseWriter, r *http.Request){
    switch(r.Method){
        case "GET":
            // redirect to home
            http.Redirect(w, r, "/", http.StatusFound);
            break;

        case "POST":
            nick     := r.FormValue("nick");
            name     := r.FormValue("name");
            email    := r.FormValue("email");
            password := r.FormValue("password");

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, _ := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();

            // insert query to db (new user)
            userUploadQuery := "insert users (nick, name, email, password, emailConfirmed, admin) values ('" + nick + "', '" + name + "','" + email + "','" + password + "', false, false);";
            _, err1 := DBConn.Query(userUploadQuery);
            if err1 != nil {
                http.Redirect(w, r, "/?msg=1&sec=2", http.StatusFound);
            }
            // insert query to db (new token)
            newToken := utils.TokenGenerator(40);
            tokenUploadQuery := "insert emailConf (name, email, token, sndDate, expDate, done) values ('" + name + "','" + email + "','" + newToken + "', current_timestamp(), current_timestamp() + INTERVAL 1 DAY, false);";
            _, err2 := DBConn.Query(tokenUploadQuery);
            if err2 != nil {
                http.Redirect(w, r, "/?msg=1&sec=2", http.StatusFound);
            }

            utils.EmailSender(name, email, newToken);

            // redirect to home |msg: account created
            http.Redirect(w, r, "/?msg=1", http.StatusFound);

            break;
    }
}

func emailCecked(w http.ResponseWriter, r *http.Request){
    confEmail := r.URL.Query().Get("email");
    confToken := r.URL.Query().Get("tok");

    // connect to db
    dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
    DBConn, _ := sql.Open("mysql", dbConnstring);
    defer DBConn.Close();

    // ceck exp date

    // update query to db (token table)
    tokenDoneQuery := "update emailConf set done=true where email='" + confEmail + "' and token='" + confToken + "';";
    _, err1 := DBConn.Query(tokenDoneQuery);
    if err1 != nil {
        fmt.Print("Email confirm err: ");
        fmt.Println(err1);
    }

    // update query to db (users table)
    userDoneQuery := "update users set emailConfirmed=true where email='" + confEmail + "';";
    _, err2 := DBConn.Query(userDoneQuery);
    if err2 != nil {
        fmt.Print("Email confirm err: ");
        fmt.Println(err2);
    }

    // html template
    Cwd, _ := os.Getwd();
    Os := runtime.GOOS
    switch Os {
        case "windows":
            template, _ := template.ParseFiles(Cwd + "\\static\\pages\\confirmedEmail.html")
            template.Execute(w, "")
            break
        default:
            template, _ := template.ParseFiles(Cwd + "/static/pages/confirmedEmail.html")
            template.Execute(w, "")
    }
}

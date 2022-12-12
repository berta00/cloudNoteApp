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
}

// db table struct
type UsersStruct struct {
    id       int
    name     string
    email    string
    password string
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
    http.HandleFunc("/register",  regRoute);
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

    newLogin := new(loginPage);
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
    //userName := "";
    //userEmail := "";

    switch(r.Method){
        case "GET":
            userIsGet = true;
            break;

        case "POST":
            email := r.FormValue("email");
            decodedPass := r.FormValue("password");

            encodedPass := utils.MD5Converter([]byte(decodedPass));

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, err := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();
            // query the users
            dbQuery := "select name, email, password, admin from users where email='" + email + "';";
            userQuery, err := DBConn.Query(dbQuery);
            if err != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err);
            } 
            // read query
            userStructQuery := new(UsersStruct);
            for userQuery.Next(){
                userQuery.Scan(&userStructQuery.name, &userStructQuery.email, &userStructQuery.password, &userStructQuery.admin);
            }
            // if password is correct
            fmt.Println(userStructQuery.password);
            if userStructQuery.password == utils.Byte16ToString(encodedPass) {
                userIsAuth = true;
            }
            
            //userName := userStructQuery.name;
            //userEmail := userStructQuery.email;

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
        // query the 

        // html template
        Cwd, _ := os.Getwd();
        Os := runtime.GOOS
        switch Os {
            case "windows":
                template, _ := template.ParseFiles(Cwd + "\\static\\pages\\dashboard.html")
                template.Execute(w, "")
                break
            default:
                template, _ := template.ParseFiles(Cwd + "/static/pages/dashboard.html")
                template.Execute(w, "")
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
            name     := r.FormValue("name");
            email    := r.FormValue("email");
            password := r.FormValue("password");

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, _ := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();

            // insert query to db (new user)
            finalPassword := utils.Byte16ToString(utils.MD5Converter([]byte(password)));
            userUploadQuery := "insert users (name, email, password, admin) values ('" + name + "','" + email + "','" + finalPassword + "', false);";
            _, err1 := DBConn.Query(userUploadQuery);
            if err1 != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err1);
            }
            // insert query to db (new token)
            newToken := utils.TokenGenerator(40);
            tokenUploadQuery := "insert emailConf (name, email, token, sndDate, expDate, done) values ('" + name + "','" + email + "','" + newToken + "', current_timestamp(), current_timestamp() + INTERVAL 1 DAY, false);";
            _, err2 := DBConn.Query(tokenUploadQuery);
            if err2 != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err2);
            }

            utils.EmailSender(name, email, newToken);

            fmt.Println("new user registered");

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

    // update query to db (token)
    tokenDoneQuery := "update emailConf set done=true where email='" + confEmail + "' and token='" + confToken + "';";
    _, err := DBConn.Query(tokenDoneQuery);
    if err != nil {
        fmt.Print("Email confirm err: ");
        fmt.Println(err);
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

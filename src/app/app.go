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

// db table struct
type UsersStruct struct {
    Id       int
    Nick     string
    Name     string
    Email    string
    Password string
    EmailCk  bool
    Admin    string
    Date     string
}
type Folder struct {
    Id      int
    Name    string
    Creator string
    Path    string
    CrDate  string
    MfDate  string
}
type BasicNoteStruct struct {
    Id       int
    Name     string
    Creator  string
    Path     string
    Content  string
    CrDate   string
    MfDate   string
}

// login page struct
type loginPage struct {
    NewAcc     bool
    WrongCred  bool
    ConfMail   bool
    Section    string
}
// dash page struct
type dashPage struct {
    Nick       string
    Name       string
    Email      string
    EmailCk    bool
    Section    string
    BasicNotes []BasicNoteStruct
    Folder     []Folder
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
    // global variables declaration
    userIsGet := false;
    userIsAuth := false;
    userNick := "";
    userName := "";
    userEmail := "";
    userEmailCk := false;
    dashSection := "";
    var folderArr    []Folder;
    var basicNoteArr []BasicNoteStruct;

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
            if err != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err);
            }
            defer DBConn.Close();
            // query the users
            dbUserQueryStr := "select nick, name, email, password, emailConfirmed, admin from users where email='" + email + "';";
            dbUserQuery, err1 := DBConn.Query(dbUserQueryStr);
            if err1 != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err1);
            } 
            // read query
            userStructQuery := new(UsersStruct);
            for dbUserQuery.Next(){
                dbUserQuery.Scan(&userStructQuery.Nick, &userStructQuery.Name, &userStructQuery.Email, &userStructQuery.Password, &userStructQuery.EmailCk, &userStructQuery.Admin);
            }
            if encodedPass == userStructQuery.Password {
                userIsAuth = true;
            }
            // query the folder
            dbFolderQueryStr := "select name, creator, path, crDate, mfDate from folder";
            dbFolderQuery, err2 := DBConn.Query(dbFolderQueryStr);
            if err2 != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err2);
            }
            // read query
            folderStructQuery := new(Folder);
            for dbFolderQuery.Next(){
                dbFolderQuery.Scan(&folderStructQuery.Name, &folderStructQuery.Creator, &folderStructQuery.Path, &folderStructQuery.CrDate, &folderStructQuery.MfDate);
                folderArr = append(folderArr, *folderStructQuery);
            }
            // query the basic notes
            dbBasicNoteQueryStr := "select name, creator, path, content, crDate, mfDate from basicNote";
            dbBasicNoteQuery, err2 := DBConn.Query(dbBasicNoteQueryStr);
            if err2 != nil {
                fmt.Print("Authentication err: ");
                fmt.Println(err2);
            }
            // read query
            basicNoteStructQuery := new(BasicNoteStruct);
            for dbBasicNoteQuery.Next(){
                dbBasicNoteQuery.Scan(&basicNoteStructQuery.Name, &basicNoteStructQuery.Creator, &basicNoteStructQuery.Path, &basicNoteStructQuery.Content, &basicNoteStructQuery.CrDate, &basicNoteStructQuery.MfDate);
                basicNoteArr = append(basicNoteArr, *basicNoteStructQuery);
            }
            
            userNick = userStructQuery.Nick;
            userName = userStructQuery.Name;
            userEmail = userStructQuery.Email;
            userEmailCk = userStructQuery.EmailCk;

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

        dashData := new(dashPage);
        dashData.Nick = userNick;
        dashData.Name = userName;
        dashData.Email = userEmail;
        dashData.EmailCk = userEmailCk;
        dashData.Section = dashSection;
        dashData.BasicNotes = basicNoteArr;
        dashData.Folder = folderArr;

        // html template
        Cwd, _ := os.Getwd();
        Os := runtime.GOOS;
        switch Os {
            case "windows":
                template, _ := template.ParseFiles(Cwd + "\\static\\pages\\dashboard.html");
                template.Execute(w, dashData);
                break
            default:
                template, _ := template.ParseFiles(Cwd + "/static/pages/dashboard.html");
                template.Execute(w, dashData);
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

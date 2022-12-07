package main

import (
    "html/template"
    "database/sql"
    "net/http"
    "strings"
    "runtime"
    "fmt"
    "os"

    "github.com/berta00/cloudNoteApp/utils"
    _ "github.com/go-sql-driver/mysql"
)

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

// JWT struct
type JWT struct {
    header  JWTheader
    payload JWTpayload
    secret  string
}
type JWTheader struct {
    alg   string
    typ   string
}
type JWTpayload struct {
    name  string
    email string
    admin string
}

// db users table struct
type UsersStruct struct {
    id       int;
    name     string;
    email    string;
    password string;
    admin    string;
    date     string;
}

func main(){
    fmt.Println("Webserver start program");

    // initialize routes
    routesInit();

    // start webserver
    err := http.ListenAndServe(":" + WSport, nil);
    fmt.Println("Webserver started at port " + WSport);
    if err == nil {
        fmt.Print("\n\nError (can't start webserver): ");
        fmt.Println(err);
    }
}

func routesInit(){
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

    http.HandleFunc("/",          mainRoute);
    http.HandleFunc("/auth",      authRoute);
    http.HandleFunc("/login",     loginRoute);
    http.HandleFunc("/validate",  validateRoute);
    http.HandleFunc("/register",  regRoute);
    http.HandleFunc("/emailceck", emailCecked);
}

func mainRoute(w http.ResponseWriter, r *http.Request){
    fmt.Println("cloud note app home service!");
}

func loginRoute(w http.ResponseWriter, r *http.Request){

    // html template
    Cwd, _ := os.Getwd();
    Os := runtime.GOOS
    switch Os {
        case "windows":
            template, _ := template.ParseFiles(Cwd + "\\static\\pages\\login.html")
            template.Execute(w, "")
            break
        default:
            template, _ := template.ParseFiles(Cwd + "/static/pages/login.html")
            template.Execute(w, "")
    }
}

func authRoute(w http.ResponseWriter, r *http.Request){
    switch(r.Method){
        case "GET":
            // redirect to home
            http.Redirect(w, r, "/", http.StatusFound);
            break;

        case "POST":
            http.Redirect(w, r, "/", http.StatusFound);

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

            // if account dont exist, return with error 1
            if userStructQuery.password == "" {
                http.Redirect(w, r, "/?err=1", http.StatusFound);
            }
            // if password ok do some things
            fmt.Println(userStructQuery.password);
            if userStructQuery.password == utils.Byte16ToString(encodedPass) {
                JWTtoken := utils.JWTgenerator(userStructQuery.name, userStructQuery.email, userStructQuery.admin);
                // response
                fmt.Println(string(JWTtoken));

                //encodedJWTtoken := HS255Converter("encode", []byte(JWTtoken));

            } else {
                http.Redirect(w, r, "/?err=2", http.StatusFound);
            }
    }
}

func validateRoute(w http.ResponseWriter, r *http.Request){
    switch(r.Method){
        case "GET":
            // redirect to home
            http.Redirect(w, r, "/", http.StatusFound);
            break;
        case "POST":
            token := r.FormValue("tok");

            // parse the token
            parsedToken := strings.Split(token, ".");

            // decode token
            NewJWT := new(JWT);
            NewJWTheader := new(JWTheader);
            NewJWTpayload := new(JWTpayload);

            for jwtI := 0; jwtI < 2; jwtI++ {
                currentSection := utils.Base64Converter("decode", []byte(parsedToken[jwtI]));
                parsedSection := strings.Split(currentSection, "\"");

                newValueFlag := false;
                newValue := "";
                validValue := 0;
                for sectionI := 0; sectionI < len(parsedSection); sectionI++ {
                    if parsedSection[sectionI] == ", " || parsedSection[sectionI] == "}" {
                        if jwtI == 0 {
                            switch validValue {
                                case 1:
                                    NewJWTheader.alg = newValue;
                                case 2:
                                    NewJWTheader.typ = newValue;
                            }
                        } else if jwtI == 1 {
                            switch validValue {
                                case 1:
                                    NewJWTpayload.name = newValue;
                                case 2:
                                    NewJWTpayload.email = newValue;
                                case 3:
                                    NewJWTpayload.admin = newValue;
                            }
                        }
                        newValueFlag = false;
                    }
                    if newValueFlag {
                        newValue += parsedSection[sectionI]
                    }
                    if parsedSection[sectionI] == ":" {
                        validValue++;
                        newValue = "";
                        newValueFlag = true;
                    }
                }
            }

            NewJWT.header = *NewJWTheader;
            NewJWT.payload = *NewJWTpayload;
            NewJWT.secret = parsedToken[2];

            // result in NewJWT
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

            http.Redirect(w, r, "/", http.StatusFound);



            break;
    }
}

func emailCecked(w http.ResponseWriter, r *http.Request){
    // SAY THIS ONLY WHEN EMAIL AND TOKEN ARE OK
    fmt.Fprint(w, "Email cecked!");

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
}

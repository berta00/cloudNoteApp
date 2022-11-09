package main

import (
    "encoding/base64"
    "encoding/json"
    "database/mysql"
    "net/http"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// db global info
var DBname string = "auth"
var DBuser string = "auth"
var DBpass string = "auth123"
var DBaddr stirng = "127.0.0.1"
var DBport String = "3306"

// db users table struct
type UsersStruct struct {
    id       int,
    name     string,
    email    string,
    password string,
    admin    string,
    date     string
}

func main(){
    // connect to db
    dbConnstring := DBuser + ":" + DBpas + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
    DBConn, err := mysql.Open("mysql", dbConnstring);

    // initialize routes
    routesInit();

    // start webserver
    err := http.ListenAndServe(":8080", nil);
    fmt.Println("Webserver started at port 8080");
    if err == nil {
        fmt.Print("\n\nError (can't start webserver): ");
        fmt.Println(err);
    }
}

func routesInit(){
    http.HandleFunc("/",         mainRoute);
    http.HandleFunc("/auth",     authRoute);
    http.HandleFunc("/register", regRoute);
}

func mainRoute(w http.ResponseWriter, r *http.Request){
    fmt.Println("Sound cloud download home service!");
}

func authRoute(w http.ResponseWriter, r *http.Request){
    switch(r.Method){
        case "GET":
            // redirect to home
            http.Redirect(w, r, "/", http.StatusFound);
            break;

        case "POST":
            email := r.FormValue("email");
            encodedPass := r.FormValue("password");

            // query the users
            dbQuery := "select name, email, password, admin from user where email='" + email + "'";
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

            // ceck password
            if userStructQuery.password == encodedPass {
                JWTtoken = JWTgenerator(userStructQuery.name, userStructQuery.email, userStructQuery.password, userStructQuery.admin);
                // response
                fmt.Println(JWTtoken);
            } else {
                http.Redirect(w, r, "/", http.StatusFound);
            }

            break;
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
            password := r.FormValue("password"); // base64 encrypted on client side

            // connect to db
            userUploadQuery = "insert user (name, email, password, admin) values ('" + name + "','" + email + "','" + password + "', false)";
            _, err := dbConn.Query(userUploadQuery);
            if err != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err);
            }

            http.Redirect(w, r, "/", http.StatusFound);

            break;
    }
}

func base64Converter(action string, string string) string {
    var returnString string;

    switch(action){
        case "decode":
            // password         (ascii decimal arr)
            decimalString, _ := base64.StdEncoding.DecodeString(encodedString);
            // ascii dec arr => char arr
            var decodedString string;
            for nChar := 0; nChar < len(decimalString); nChar++ {
                decodedString += string(decimalString[nChar]);
            }

            returnString = decodedString;
            break;

        case "encode":
            encodedString, _ := base64.StdEncoding.EncodeString(encodedString);

            returnString = encodedString;
            break;

        default:
            returnString = "base64Converter err: function action parameter";
    }

    return returnString;
}

/*  Json Web Token format

    Header:                 Payload:                              Secret:
    {                       {                                     ziopera
        "alg": "HS256",         "name":  "GiovanniRossi",
        "typ": "JWT"            "email": "giovanni@google.com",
    }                           "admin": "true"
                            }

    REDY TO SEND: "[base64(header)].[base64(payload)].[secret]"
*/
func JWTgenerator(name string, email string, admin string) string {
    // create json element
    header  := "{'alg':'HS256','typ':'JWT'}";
    payload := "{'name':'" + name + "','email':'" + email + "','" + admin + "':'true'}";
    secret  := "ziopera";
    // create finale JWT
    JWTtoken := base64Converter("encode", header) + "." + base64Converter("encode", payload) + "." + secret;

    return JWTtoken;
}

package main

import (
    "encoding/base64"
    "database/sql"
    "net/http"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// db global info
var DBname string = "auth";
var DBuser string = "auth";
var DBpass string = "auth123";
var DBaddr string = "127.0.0.1";
var DBport string = "3306";

// db users table struct
type UsersStruct struct {
    id       int;
    name     string;
    email    string;
    password string;
    admin    string;
    date     string;
}

// JWT struct
type JWTheader struct {
    alg   string
    typ   string
}
type JWTpayload struct {
    name  string
    email string
    admin string
}

func main(){
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

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, err := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();

            // query the users
            dbQuery := "select name, email, password, admin from users where email='" + email + "'";
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
                JWTtoken := JWTgenerator(userStructQuery.name, userStructQuery.email, userStructQuery.admin);
                // response
                fmt.Println(JWTtoken);
                base64Converter("decode", "cm9vdDEyMw==");
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
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, _ := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();

            // connect to db
            userUploadQuery := "insert users (name, email, password, admin) values ('" + name + "','" + email + "','" + password + "', false)";
            _, err := DBConn.Query(userUploadQuery);
            if err != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err);
            }

            http.Redirect(w, r, "/", http.StatusFound);

            break;
    }
}

func base64Converter(action string, string string) string {
    returnString := "";

    switch(action){
        case "decode":/*   --RIVEDERE--
            //                    (ascii decimal arr)
            decimalString, err := base64.StdEncoding.DecodeString(string);
            if err != nil {
                fmt.Print("Base64 decode err: ");
                fmt.Println(err);
            }
            // ascii dec arr => char arr
            decodedString := strings.Join(decimalString, "");
            // return
            returnString = decodedString;   */
            break;

        case "encode":
            encodedString := base64.StdEncoding.EncodeToString([]byte(string));
            // return
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
    jsonHeader  := []byte(`{"alg":"HS256", "typ":"JWT"}`);
    jsonPayload := []byte(`{"name":"`+name+`", "email":"`+email+`", "admin":"`+admin+`"}`);
    secret  := "ziopera";
    // create finale JWT
    encodedHeader := base64.StdEncoding.EncodeToString(jsonHeader);
    encodedPayload := base64.StdEncoding.EncodeToString(jsonPayload);
    encodedSecret := base64.StdEncoding.EncodeToString([]byte(secret));

    var finalEncodedHeader string;
    var finalEncodedPayload string;
    var finalEncodedSecret string;

    for a := 0; a < len(encodedHeader); a++ {
        if string(encodedHeader[a]) != "=" {
            finalEncodedHeader += string(encodedHeader[a]);
        } else {
            a = len(encodedHeader);
        }
    }
    for b := 0; b < len(encodedPayload); b++ {
        if string(encodedPayload[b]) != "=" {
            finalEncodedPayload += string(encodedPayload[b]);
        } else {
            b = len(encodedPayload);
        }
    }
    for c := 0; c < len(encodedSecret); c++ {
        if string(encodedSecret[c]) != "=" {
            finalEncodedSecret += string(encodedSecret[c]);
        } else {
            c = len(encodedSecret);
        }
    }

    JWTtoken := finalEncodedHeader + "." + finalEncodedPayload + "." + finalEncodedSecret;

    return JWTtoken;
}

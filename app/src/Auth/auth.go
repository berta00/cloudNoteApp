package main

import (
    "encoding/base64"
	"crypto/sha256"
    "crypto/hmac"
    "database/sql"
    "math/rand"
    "net/http"
    "net/smtp"
    "strings"
    "time"
    "fmt"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

// ws global info
var WSdomain string = "localhost";
var WSport   string = "8080";
var JWTsec   string = os.Getenv("JWT_SECRET");

// db global info
var DBname string = os.Getenv("MYSQL_HOST");
var DBuser string = os.Getenv("MYSQL_USER");
var DBpass string = os.Getenv("MYSQL_PASSWORD");
var DBaddr string = os.Getenv("MYSQL_DB");
var DBport string = os.Getenv("MYSQL_PORT");

// gmail access info
var GMemail string = ;
var GMpass string = os.Getenv("GMAIL_PASSWORD");

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
    http.HandleFunc("/",          mainRoute);
    http.HandleFunc("/auth",      authRoute);
    http.HandleFunc("/validate",  validateRoute);
    http.HandleFunc("/register",  regRoute);
    http.HandleFunc("/emailceck", emailCecked);
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

            // ceck password
            if userStructQuery.password == encodedPass {
                JWTtoken := JWTgenerator(userStructQuery.name, userStructQuery.email, userStructQuery.admin);
                // response
                fmt.Println(string(JWTtoken));

                //encodedJWTtoken := HS255Converter("encode", []byte(JWTtoken));

            } else {
                http.Redirect(w, r, "/", http.StatusFound);
            }

            break;
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
                currentSection := base64Converter("decode", parsedToken[jwtI]);
                parsedSection := strings.Split(currentSection, "\"");

                newValueFlag := false;
                newValue := "";
                validValue := 0;
                for sectionI := 0; sectionI < len(parsedSection) ; sectionI++ {
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
            password := r.FormValue("password"); // base64 encrypted on client side

            // connect to db
            dbConnstring := DBuser + ":" + DBpass + "@tcp(" + DBaddr + ":" + DBport + ")/" + DBname;
            DBConn, _ := sql.Open("mysql", dbConnstring);
            defer DBConn.Close();

            // insert query to db (new user)
            userUploadQuery := "insert users (name, email, password, admin) values ('" + name + "','" + email + "','" + password + "', false);";
            _, err1 := DBConn.Query(userUploadQuery);
            if err1 != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err1);
            }
            // insert query to db (new token)
            newToken := TokenGenerator(40);
            tokenUploadQuery := "insert emailConf (name, email, token, sndDate, expDate, done) values ('" + name + "','" + email + "','" + newToken + "', current_timestamp(), current_timestamp() + INTERVAL 1 DAY, false);";
            _, err2 := DBConn.Query(tokenUploadQuery);
            if err2 != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err2);
            }

            emailSender(name, email, newToken);

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

    // update query to db (token)
    tokenDoneQuery := "update emailConf set done=true where email='" + confEmail + "' and token='" + confToken + "';";
    _, err := DBConn.Query(tokenDoneQuery);
    if err != nil {
        fmt.Print("Registration err: ");
        fmt.Println(err);
    }
}

func emailSender(name string, destinatioEmail string, token string){
    emailLink := WSdomain + ":" + WSport + "/emailceck?email=" + destinatioEmail + "&tok=" + token;
    // email variable
    senderEmail := GMemail;
    senderPass := GMpass;
    reciverMail := []string{destinatioEmail};
    // host variable
    host := "smtp.gmail.com"
    port := "587"
    address := host + ":" + port
    // message variable
    subject := "Subject: cloud note app email verification";
    body := "\nHi " + name + ",\n click here: " + emailLink + " to verify your email on the platform!";
    message := []byte(subject + body);

    // auth in mail service
    auth := smtp.PlainAuth("", senderEmail, senderPass, host);
    // send email
    err := smtp.SendMail(address, auth, senderEmail, reciverMail, message);
    if err != nil {
        fmt.Print("err (send email): ");
        fmt.Println(err);
    }
}

func base64Converter(action string, string string) string {
    returnString := "";

    switch(action){
        case "decode":
            //                    (ascii decimal arr)
            decimalString, err := base64.RawURLEncoding.DecodeString(string);
            if err != nil {
                fmt.Print("Base64 decode err: ");
                fmt.Println(err);
            }
            finalString := fmt.Sprintf("%s", decimalString);
            // return
            returnString = finalString;
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

func HS255Converter(action string, string []byte) []byte {
    var returnString []byte;

    switch(action){
        case "decode":

            break;
        case "encode":
            hasher := hmac.New(sha256.New, []byte(JWTsec));
            _, err := hasher.Write(string);
        	if err != nil {
        		return []byte("err encoding the string");
            }
        	r := hasher.Sum(nil);
            returnString = r;

            break;

        default:
            returnString = []byte("HS255Converter err: function action parameter");
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
func JWTgenerator(name string, email string, admin string) []byte {
    // create json element
    jsonHeader  := []byte(`{"alg":"HS256", "typ":"JWT"}`);
    jsonPayload := []byte(`{"name":"`+name+`", "email":"`+email+`", "admin":"`+admin+`"}`);
    secret := JWTsec;

    // create finale JWT
    encodedHeader := base64.RawURLEncoding.EncodeToString(jsonHeader);
    encodedPayload := base64.RawURLEncoding.EncodeToString(jsonPayload);

    JWTtoken := encodedHeader + "." + encodedPayload + "." + secret;

    return []byte(JWTtoken);
}

func TokenGenerator(secLen int) string {
    symbol := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9"};
    finalSecret := "";

    for a := 0; a < secLen; a++ {
        rand.Seed(time.Now().UnixNano())

        finalSecret += symbol[rand.Intn(len(symbol))];
    }

    return finalSecret;
}

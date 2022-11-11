package main

import (
    "encoding/base64"
    "database/sql"
    "math/rand"
    "net/http"
    "net/smtp"
    "time"
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
    http.HandleFunc("/",           mainRoute);
    http.HandleFunc("/auth",       authRoute);
    http.HandleFunc("/register",   regRoute);
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

            // insert query to db (new user)
            userUploadQuery := "insert users (name, email, password, admin) values ('" + name + "','" + email + "','" + password + "', false);";
            _, err1 := DBConn.Query(userUploadQuery);
            if err1 != nil {
                fmt.Print("Registration err: ");
                fmt.Println(err1);
            }
            // insert query to db (new token)
            newToken := SECRETgenerator(40);
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
    tokenDoneQuery := "update SET done=true where email='" + confEmail + "' and token='" + confToken + "';";
    _, err := DBConn.Query(tokenDoneQuery);
    if err != nil {
        fmt.Print("Registration err: ");
        fmt.Println(err);
    }
}

func emailSender(name string, destinatioEmail string, token string){
    emailLink := "/localhost/emailceck?email=" + destinatioEmail + "&tok=" + token;
    // email variable
    senderEmail := "soundclouddownloader00@gmail.com";
    senderPass := "";
    reciverMail := []string{destinatioEmail};
    // host variable
    host := "smtp.gmail.com"
    port := "587"
    address := host + ":" + port
    // message variable
    subject := "Subject: SoundCloud download service email verification";
    body := "Hi " + name + ",\nclick <a href='" + emailLink + "'>here</a> to verify your email on the platform!";
    message := []byte(subject+body);

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
    secret  := SECRETgenerator(20);

    // create finale JWT
    encodedHeader := base64.StdEncoding.EncodeToString(jsonHeader);
    encodedPayload := base64.StdEncoding.EncodeToString(jsonPayload);
    //encodedSecret := base64.StdEncoding.EncodeToString([]byte(secret));

    var finalEncodedHeader string;
    var finalEncodedPayload string;
    //var finalEncodedSecret string;

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
    /* NOT BASE64 ENCODED FOR NOW
    for c := 0; c < len(encodedSecret); c++ {
        if string(encodedSecret[c]) != "=" {
            finalEncodedSecret += string(encodedSecret[c]);
        } else {
            c = len(encodedSecret);
        }
    }
    */

    JWTtoken := finalEncodedHeader + "." + finalEncodedPayload + "." + secret;

    return JWTtoken;
}

func SECRETgenerator(secLen int) string {
    symbol := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9"};
    finalSecret := "";

    for a := 0; a < secLen; a++ {
        rand.Seed(time.Now().UnixNano())

        finalSecret += symbol[rand.Intn(len(symbol))];
    }

    return finalSecret;
}

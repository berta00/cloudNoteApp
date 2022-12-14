package utils;

import(
    "encoding/base64"
    "encoding/hex"
    "crypto/md5"
    "math/rand"
    "net/smtp"
    "net/mail"
    "os/exec"
    "runtime"
    "time"
    "fmt"
    "os"
	"github.com/scorredoira/email"
)

// ws global info
var WSdomain string = "localhost";
var WSport   string = os.Getenv("WS_PORT");
var JWTsec   string = os.Getenv("JWT_SECRET");

// db global info
var DBname string = os.Getenv("MYSQL_HOST");
var DBuser string = os.Getenv("MYSQL_USER");
var DBpass string = os.Getenv("MYSQL_PASSWORD");
var DBaddr string = os.Getenv("MYSQL_DB");
var DBport string = os.Getenv("MYSQL_PORT");

// gmail access info
var GMemail string = os.Getenv("GMAIL");
var GMpass string = os.Getenv("GMAIL_PASSWORD");

func EnvVarSet(){
    finalPath1 := "Users/marcobertagnolli/Desktop/Programmazione/cloudNoteApp/src/app/utils/testEnvSetup/config.sh";
    finalPath2 := "./utils/testEnvSetup/secret.sh";
    if runtime.GOOS != "windows" {
        out1, err1 := exec.Command("source", "/", finalPath1).Output();
        if err1 != nil {
            fmt.Println("- Err setting the config var: " + err1.Error());
        } else {
            fmt.Println("- " + string(out1[:]));
        }
        out2, err2 := exec.Command("source", finalPath2).Output();
        if err2 != nil {
            fmt.Println("- Err setting the secret var: " + err2.Error());
        } else {
            fmt.Println("- " + string(out2[:]));
        }
    } else {
        fmt.Println("windows is not compatible to run this program (for now)");
    }
}

func EmailSender(name string, destinatioEmail string, token string){
    emailLink := WSdomain + ":" + WSport + "/emailceck?email=" + destinatioEmail + "&tok=" + token;
    // email variable
    senderEmail := GMemail;
    senderPass := GMpass;
    reciverMail := []string{destinatioEmail};
    // host variable
    host := "smtp.gmail.com"
    port := "587"
    address := host + ":" + port
    // message
    finalBody := "Hi " + name + ",<br> click on this link: " + emailLink + "<br>to verify your email.";
    finalEmail := email.NewHTMLMessage("Email verification", finalBody)
    finalEmail.From = mail.Address{Name: "Cloud notes", Address: senderEmail};
    finalEmail.To = reciverMail;

    // auth in mail service
    auth := smtp.PlainAuth("", senderEmail, senderPass, host);
    // send email
	if err := email.Send(address, auth, finalEmail); err != nil {
		fmt.Print("email err: ")
        fmt.Println(err)
	}
}

func Base64Converter(action string, parString []byte) string {
    returnString := "";

    switch(action){
        case "decode":
            //                    (ascii decimal arr)
            decimalString, err := base64.RawURLEncoding.DecodeString(string(parString));
            if err != nil {
                fmt.Print("Base64 decode err: ");
                fmt.Println(err);
            }
            finalString := fmt.Sprintf("%s", decimalString);
            // return
            returnString = finalString;
            break;

        case "encode":
            encodedString := base64.StdEncoding.EncodeToString(parString);
            // return
            returnString = encodedString;
            break;

        default:
            returnString = "base64Converter err: function action parameter";
    }

    return returnString;
}

func MD5Converter(parString []byte) [16]byte {
    salt := os.Getenv("PASS_SALT");
    finalString := salt + string(parString);
    
    return md5.Sum([]byte(finalString));
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

func TokenGenerator(secLen int) string {
    symbol := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9"};
    finalSecret := "";

    for a := 0; a < secLen; a++ {
        rand.Seed(time.Now().UnixNano())

        finalSecret += symbol[rand.Intn(len(symbol))];
    }

    return finalSecret;
}

func Byte16ToString(parString [16]byte) string {
    return hex.EncodeToString([]byte(parString[:]));
}

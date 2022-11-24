package main;

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

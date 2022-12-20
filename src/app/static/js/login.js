let loadDiv = document.querySelector(".loadDiv");
let loginMain = document.querySelector("#login");
let registerMain = document.querySelector("#register");
let loginEmail = document.querySelector("#login .email");
let loginPwd = document.querySelector("#login .pwd");
let registerNickName = document.querySelector("#register .nickName");
let registerName = document.querySelector("#register .fullName");
let registerEmail = document.querySelector("#register .email");
let registerPwd = document.querySelector("#register .pwd");
let signButton = document.querySelector(".main div a");
let loginButton = loginMain.querySelector("div button");
let registerButton = registerMain.querySelector("div button");
let signInText = document.querySelector("#login .buttons a");
let logInText = document.querySelector("#register .buttons a");

if(section == "2"){
    loginMain.style.display = "none";
    registerMain.style.display = "flex";
}

loginButton.addEventListener("click", ()=> {
    loginEmail.style.display = "none";
    loginPwd.style.display = "none";
    signButton.style.display = "none";
    loginButton.style.display = "none";
    loadDiv.style.display = "block";

    let emailValue = loginEmail.value;
    let pwdValue = loginPwd.value;
    //perfavore non guardatelo grazie mille buonagiornata
    let encPwd = md5("kdabfaxjbcjkabldfasdfjlablfjbashbhakv" + pwdValue);

    // create form and submit
    loginMain.innerHTML = "<form class='loginForm' method='POST' action='/dash' style='display: none;'><input type='hidden' name='email' value='" + emailValue + "'><input type='hidden' name='password' value='" + encPwd + "'><input type='submit'></form>";
    loginMain.querySelector(".loginForm").submit();
});
registerButton.addEventListener("click", ()=> {
    registerName.style.display = "none";
    registerEmail.style.display = "none";
    registerPwd.style.display = "none";
    loginButton.style.display = "none";
    loadDiv.style.display = "block";

    let nickValue = registerNickName.value;
    let nameValue = registerName.value;
    let emailValue = registerEmail.value;
    let pwdValue = registerPwd.value;
    let encPwd = md5("kdabfaxjbcjkabldfasdfjlablfjbashbhakv" + pwdValue);

    // create form and submit
    registerMain.innerHTML = "<form class='regForm' method='POST' action='/reg' style='display: none;'><input type='hidden' name='nick' value='" + nickValue + "'><input type='hidden' name='name' value='" + nameValue + "'><input type='hidden' name='email' value='" + emailValue + "'><input type='hidden' name='password' value='" + encPwd + "'><input type='submit'></form>";
    registerMain.querySelector(".regForm").submit();
});

signInText.addEventListener("click", ()=> {
    loginMain.style.display = "none";
    registerMain.style.display = "flex";
});
logInText.addEventListener("click", ()=> {
    registerMain.style.display = "none";
    loginMain.style.display = "flex";
});

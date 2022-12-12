let main = document.querySelector(".main");
let loginEmail = document.querySelector(".main .email");
let loginPwd = document.querySelector(".main .pwd");
let signButton = document.querySelector(".main div a");
let loginButton = document.querySelector(".main div button");
let loadDiv = document.querySelector(".loadDiv");

// default sets
loadDiv.style.display = "none";

loginButton.addEventListener("click", ()=> {
    loginEmail.style.display = "none";
    loginPwd.style.display = "none";
    signButton.style.display = "none";
    loginButton.style.display = "none";
    loadDiv.style.display = "block";

    let emailValue = loginEmail.value;
    let pwdValue = loginPwd.value;
    let encPwd = md5("kdabfaxjbcjkabldfasdfjlablfjbashbhakv" + pwdValue);
    console.log(encPwd);
    // create form and submit
    main.innerHTML = "<form class='loginForm' method='POST' action='/dash' style='display: none;'><input type='hidden' name='email' value='" + emailValue + "'><input type='hidden' name='password' value='" + encPwd + "'><input type='submit'></form>";
    main.querySelector(".loginForm").submit();
});

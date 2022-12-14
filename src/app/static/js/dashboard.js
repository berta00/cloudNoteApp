// header variables
let pathP = document.querySelector(".header .path");
pathP.innerHTML = userNick + " /";
// browse section variables
let browseDiv = document.querySelector(".browseSection");
// file section variables
let fileDiv = document.querySelector(".fileSection");

// section manage
function switchSection(destinationSection, fileName, fileType){
    switch(destinationSection){
        case "browseSection":
            // header text
            pathP.innerHTML = userNick + " /";
            // main sections
            browseDiv.style.display = "flex";
            fileDiv.style.display = "none";
            break;
        case "fileSection":
            // header text
            pathP.innerHTML = userNick + " / " + fileName + " <a class='fileType'>[" + fileType + "]</a>";
            // main sections
            browseDiv.style.display = "flex";
            fileDiv.style.display = "none";
            break;
    }
}

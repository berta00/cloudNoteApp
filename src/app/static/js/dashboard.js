// header variables
let titlePath = document.querySelector(".header .path");
titlePath.innerHTML = userNick + " /";
// browse section variables
let browseMainDiv = document.querySelector(".browseMainSection");
// file section variables
let noteMainDiv = document.querySelector(".noteMainSection");
// get icon for note app
let noteToolDockDivs = document.querySelectorAll(".noteMainSection .main .toolDock div");
let noteToolDockList = ["marker", "txtColor", "txtSize", "txtFont", "txtBold", "txtItalic", "txtUnderline", "txtHeading", "txtNoStyle", "txtLink", "txtAlign", "txtList"];
let firstStartBrowseSection = true;
let firstStartNoteSection = true;
// default section
switchSection("browseSection");

// section manage
function switchSection(destinationSection, fileName, fileType){
    switch(destinationSection){
        case "browseSection":
            // header text
            titlePath.innerHTML = userNick + " /";
            // main sections
            browseMainDiv.style.display = "flex";
            noteMainDiv.style.display = "none";
            if(firstStartBrowseSection){
                // main table element
                let browseSectionTable = document.createElement("div");
                browseSectionTable.style.height = "auto";
                browseSectionTable.style.width = "calc(100% - 70px)";
                browseSectionTable.style.marginLeft = "35px";
                browseSectionTable.style.marginTop = "40px";
                browseSectionTable.style.display = "flex";
                browseSectionTable.style.flexDirection = "column";
                browseSectionTable.style.gap = "6px";
                browseMainDiv.querySelector(".main .centralSection").appendChild(browseSectionTable);
                // heading row
                newBrowseSectionRow(browseSectionTable, "", "Name:", "Creator:", "Creation date:", "Modify date:", "Size:");

                // separetion line
                let rowSeparationLine = document.createElement("div");
                rowSeparationLine.style.height = "4px";
                rowSeparationLine.style.width = "100%";
                rowSeparationLine.style.borderRadius = "2px";
                rowSeparationLine.style.background = "#000000";

                browseSectionTable.appendChild(rowSeparationLine);
                
                // fs deep
                let remainingFolder = folder;
                let remainingDeletedFolder = [];
                let previousFolders = 1;
                let previousFolders1 = ["/"];
                let previousFolders2 = [];
                let remainingFile = file;
                let remainingDeletedFile = [];
                let deepness = 0;
                let finishedFs = false;
                while(finishedFs == false){
                    console.log(remainingFolder);
                    // create folder
                    let folderRowIndex = 0;
                    let folderRowInterval = setInterval(()=>{
                        if(folderRowIndex >= remainingFolder.length){
                            clearInterval(folderRowInterval);
                        } else {
                            for(let previousIndex = 0; previousIndex < previousFolders1.length + previousFolders2.length; previousIndex++){
                                if(remainingFolder[folderRowIndex][2] == previousFolders1[previousIndex] || remainingFolder[folderRowIndex][2] == previousFolders2[previousIndex]){
                                    let folderSize = "18 item";
                                    if(previousFolders1[previousIndex] == "/" || previousFolders2[previousIndex] == "/"){
                                        newBrowseSectionRow(browseSectionTable, "/static/icons/browseSection-listFolder.svg", remainingFolder[folderRowIndex][0], remainingFolder[folderRowIndex][1], remainingFolder[folderRowIndex][3], remainingFolder[folderRowIndex][4], folderSize);
                                    } else {
                                        newBrowseSectionRow(browseSectionTable.querySelector("." + previousFolders[previousIndex]), "/static/icon/browseSection-listFolder.svg", remainingFolder[folderRowIndex][0], remainingFolder[folderRowIndex][1], remainingFolder[folderRowIndex][3], remainingFolder[folderRowIndex][4], folderSize);
                                    }
                                    switch(previousFolders){
                                        case 1:
                                            previousFolders2 = [];
                                            previousFolders1.push(remainingFolder[folderRowIndex]);
                                            previousFolders = 2;
                                        case 2:
                                            previousFolders1 = [];
                                            previousFolders2.push(remainingFolder[folderRowIndex]);
                                            previousFolders = 1;
                                    }
                                    remainingDeletedFolder.push(remainingFolder.splice(folderRowIndex, 1));
                                }
                            }
                        }
                        folderRowIndex++;
                    }, 50);
                    console.log(remainingFolder);
                    // create file
                    /*
                    let fileRowIndex = 0;
                    let fileRowInterval = setInterval(()=>{
                        if(fileRowIndex >= remainingFile.length){
                            clearInterval(fileRowInterval);
                        } else {
                            for(let previousIndex = 0; previousIndex < previousFolders.length; previousIndex++){
                                if(previousFolders[fileRowIndex][3] == previousFolders[previousIndex]){
                                    let fileSize = "34 byte";
                                    if(previousFolders[fileRowIndex] == "/"){
                                        newBrowseSectionRow(browseSectionTable, "/static/icon/browseSection-listFile.svg", remainingFile[fileRowIndex][1], remainingFile[fileRowIndex][2], remainingFile[fileRowIndex][4], remainingFile[fileRowIndex][5], fileSize);
                                    } else {
                                        newBrowseSectionRow(browseSectionTable.querySelector("." + previousFolders[fileRowIndex]), "/static/icon/browseSection-listFile.svg", remainingFile[fileRowIndex][1], remainingFile[fileRowIndex][2], remainingFile[fileRowIndex][4], remainingFile[fileRowIndex][5], fileSize);
                                    }
                                    remainingFile.splice(fileRowIndex, 1);
                                }
                            }
                        }
                        fileRowIndex++;
                    }, 50);*/
                    if(remainingDeletedFolder.length >= folder.length){
                        finishedFs = true;
                    }
                    deepness++;
                }
            }
            firstStartBrowseSection = false;
            break;
        case "noteSection":
            // header text
            titlePath.innerHTML = userNick + " / " + fileName + " <a class='fileType'>[" + fileType + "]</a>";
            // main sections
            browseMainDiv.style.display = "none";
            noteMainDiv.style.display = "flex";
            // animate dock toolbox
            if(firstStartNoteSection){
                let iconIndex = 0;
                let dockIconAnimationInterval = setInterval(()=>{
                    if(iconIndex >= noteToolDockDivs.length){
                        clearInterval(dockIconAnimationInterval);
                    } else {
                        noteToolDockDivs[iconIndex].style.opacity = "1";
                    }
                    iconIndex++
                }, 50);
            }
            // initialize the tool functions
            if(firstStartNoteSection){
                let toolIndex = 0;
                let dockToolInitializeInterval = setInterval(()=>{
                    if(toolIndex >= noteToolDockList.length){
                        clearInterval(dockToolInitializeInterval);
                    } else {
                        let currentNum = toolIndex;
                        noteToolDockDivs[currentNum].addEventListener("click", ()=>{
                            let currentTool = noteToolDockList[currentNum];
                            // border
                            for(let borderIndex = 0; borderIndex < noteToolDockDivs.length; borderIndex++){
                                if(borderIndex == currentNum){
                                    noteToolDockDivs[borderIndex].style.outline = "2px solid black";
                                    noteToolDockDivs[borderIndex].style.outlineOffset = "2px";
                                } else {
                                    noteToolDockDivs[borderIndex].style.outline = "0 solid";
                                    noteToolDockDivs[borderIndex].style.outlineOffset = "0";
                                }
                            }
                            // functioning

                        });
                        toolIndex++;
                    }
                }, 2);
            }
            firstStartNoteSection = false;
            break;
    }
}

function newBrowseSectionRow(parent, icon, name, creator, creation, lastModify, size){
    // row div
    let rowDiv = document.createElement("div");
    rowDiv.style.height = "auto";
    rowDiv.style.width = "100%";
    rowDiv.style.display = "flex";
    rowDiv.style.flexDirection = "row";
    rowDiv.style.alignItems = "center";
    rowDiv.className = name;
    // icon
    if(icon == ""){
        let rowIcon = document.createElement("div");
        rowIcon.style.height = "32px";
        rowIcon.style.width = "32px";
        rowDiv.appendChild(rowIcon);
    } else {
        let rowIcon = document.createElement("img");
        rowIcon.src = icon;
        rowIcon.style.height = "32px";
        rowDiv.appendChild(rowIcon);
    }
    // name
    let rowName = document.createElement("a");
    rowName.innerHTML = name;
    rowName.style.fontFamily = "'Roboto', sans-serif";
    rowName.style.fontSize = "100%";
    rowName.style.color = "#000000";
    rowName.style.marginLeft = "1.4%";
    rowName.style.width = "15%";
    rowDiv.appendChild(rowName);
    // creator
    let rowCreator = document.createElement("a");
    rowCreator.innerHTML = creator;
    rowCreator.style.fontFamily = "'Roboto', sans-serif";
    rowCreator.style.fontSize = "100%";
    rowCreator.style.color = "#000000";
    rowCreator.style.width = "30%";
    rowDiv.appendChild(rowCreator);
    // creation
    let rowCreation = document.createElement("a");
    rowCreation.innerHTML = creation;
    rowCreation.style.fontFamily = "'Roboto', sans-serif";
    rowCreation.style.fontSize = "100%";
    rowCreation.style.color = "#000000";
    rowCreation.style.width = "20%";
    rowDiv.appendChild(rowCreation);
    // last modify
    let rowLastModify = document.createElement("a");
    rowLastModify.innerHTML = lastModify;
    rowLastModify.style.fontFamily = "'Roboto', sans-serif";
    rowLastModify.style.fontSize = "100%";
    rowLastModify.style.color = "#000000";
    rowLastModify.style.width = "20%";
    rowDiv.appendChild(rowLastModify);
    // size
    let rowSize = document.createElement("a");
    rowSize.innerHTML = size;
    rowSize.style.fontFamily = "'Roboto', sans-serif";
    rowSize.style.fontSize = "100%";
    rowSize.style.color = "#000000";
    rowSize.style.width = "auto%";
    rowDiv.appendChild(rowSize);

    parent.appendChild(rowDiv);
}
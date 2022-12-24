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
                
                // file system
                let remainingFiles = fileSystenData;
                let preveousLayerFiles = ["/"];
                let deletedFiles = [];
                let layer = 0;
                let finishedFs = false;

                // loop throught all layers
                let browseFsReader = setInterval(()=>{
                    if(finishedFs){
                        clearInterval(browseFsReader);
                    } else {
                        console.log(remainingFiles);
                        // create element
                        for(let createI = 0; createI < remainingFiles.length; createI++){
                            for(let prevI = 0; prevI < preveousLayerFiles.length; prevI++){
                                if(preveousLayerFiles[prevI] == remainingFiles[createI][3]){
                                    // choose icon
                                    let icon = "";
                                    switch(remainingFiles[createI][0]){
                                        case "folder":
                                            icon = "/static/icons/browseSection-listFolder.svg";
                                            break;
                                        case "basicNote":
                                            icon = "/static/icons/browseSection-listFile.svg";
                                            break;
                                    }
                                    // generate element
                                    let size = "34 byte";
                                    if(preveousLayerFiles[prevI] == "/"){
                                        newBrowseSectionRow(browseSectionTable, icon, remainingFiles[createI][1], remainingFiles[createI][2], remainingFiles[createI][4], remainingFiles[createI][5], size);
                                    } else {
                                        newBrowseSectionRow(browseSectionTable.querySelector(preveousLayerFiles[prevI]), icon, remainingFiles[createI][1], remainingFiles[createI][2], remainingFiles[createI][4], remainingFiles[createI][5], size);
                                    }
                                    // push created element in deleted
                                    deletedFiles.push(remainingFiles[createI]);
                                    remainingFiles.splice(createI, 1);
                                }
                            }
                        }
                        // delete element from remainingFiles and put it in preveous
                        for(let remainI = 0; remainI < deletedFiles.length; remainI++){
                            preveousLayerFiles.push(deletedFiles[remainI][1]);
                        }
                        // reset arrays
                        preveousLayerFiles = [];
                        deletedFiles = [];
                        // loop exit contition
                        if(remainingFiles.length <= 0){
                            finishedFs = true;
                        }
                        layer++;
                    }
                }, 50);
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
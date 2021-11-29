"use strict"

console.log("hello world")

//TODO  send the lock (integer) status to /updateLock/:lockMode
function sendLockUpdate(newLockMode){

}


// logs time
function convertEpochToSpecificTimezone(timeEpoch, offset) {
    let d = new Date(timeEpoch);
    let utc = d.getTime() + (d.getTimezoneOffset() * 60000);  //This converts to UTC 00:00
    let nd = new Date(utc + (3600000 * offset));
    return nd.toLocaleString();
}

//logs table
function generateTable() {

    axios({
        method: "get",
        url: "http://localhost:8080/userLogs",
        data: null,
    })
        .then(function (response) {
            //handle success

            let string1 = JSON.stringify(response);
            let parsed = JSON.parse(string1);

            let teachTotal = 0
            let xTotal = 0
            let vTotal = 0
            let currentCount = 0

            for (let i = 0; i < parsed.data.length; i++) {

                // This could have been in the database
                if (parsed.data[i].user === "teach") {
                    teachTotal++
                    currentCount = teachTotal
                }
                if (parsed.data[i].user === "x") {
                    xTotal++
                    currentCount = xTotal
                }
                if (parsed.data[i].user === "v") {
                    vTotal++
                    currentCount = vTotal
                }

                if (parsed.data[i].user === "INCORRECT CREDENTIALS") {
                    currentCount = 0
                }

                let tr = document.createElement('tr')

                let td1 = document.createElement('th')
                let td2 = document.createElement('td')
                let td3 = document.createElement('td')
                let td4 = document.createElement('td')
                let text1 = document.createTextNode((i + 1).toString())
                let text2 = document.createTextNode(parsed.data[i].user)
                let text3 = document.createTextNode(convertEpochToSpecificTimezone(parsed.data[i].time, +3))
                let text4 = document.createTextNode(currentCount.toString())
                td1.appendChild(text1)
                td2.appendChild(text2)
                td3.appendChild(text3)
                td4.appendChild(text4)
                tr.appendChild(td1)
                tr.appendChild(td2)
                tr.appendChild(td3)
                tr.appendChild(td4)

                document.getElementById("logContent").appendChild(tr)

            }

        })
        .catch(function (response) {
            //handle error
            console.log(response);
        });


}


// clear database when log out
async function clearDB() {
    try {
        return await axios({
            url: 'http://localhost:8080/clearDatabase',
            method: 'Delete',
            timeout: 8000,
            headers: {
                'Content-Type': 'application/json',
            }
        })
    } catch (err) {
        console.error(err);
    }
}



window.addEventListener('DOMContentLoaded', (event) => {
    console.log("DOM LOADED")
    const lockImage = document.getElementById("lockImage")
    // const doorImage = document.getElementById("doorImage")
    const lockRadio = document.getElementsByName("mode")
    const applyButton = document.getElementById("applyButton")
    //const changeButton = document.getElementById("changeButton")
    const modeContainer = document.getElementById("modeContainer")
    // const statusContainer = document.getElementById("statusContainer")

    const homeButton = document.getElementById("homeButton")
    const historyButton = document.getElementById("historyButton")

    const logoutButton = document.getElementById("logoutButton")

    //const statusButton = document.getElementById("statusButton")

    const homeDiv = document.getElementById("homeContent")
    const historyDiv = document.getElementById("historyContent")
    //const statusDiv = document.getElementById("statusContent")

    homeDiv.style.display = "block"
    historyDiv.style.display = "none"
    //statusDiv.style.display = "none"

    let currentLockStatus = {
        UNLOCKED: 0,
        SOFT: 1,
        HARD: 2,
    };

    //  Home button handling
    homeButton.addEventListener("click", () => {
        console.log("homeButton clicked.")
        homeDiv.style.display = "block"
        historyDiv.style.display = "none"
        statusDiv.style.display = "none"
    });


    //  history button handling
    historyButton.addEventListener("click", () => {
        console.log("historyButton clicked.")
        homeDiv.style.display = "none"
        historyDiv.style.display = "block"
        statusDiv.style.display = "none"

    });


    // Database button
    logout.addEventListener("click", () => {
        console.log("database button pressed")
        clearDB().then(r => console.log(r))
    });



    applyButton.addEventListener("click", function () {


        for (let i = 0; i < lockRadio.length; i++) {
            if (lockRadio[i].checked) {
                console.log(i)

                switch (i) {
                    case 0:     // OPEN

                        console.log("OPEN")
                        if (currentLockStatus !== 0) {
                            lockImage.src = "img/lockOpenAnim.png"
                            modeContainer.innerHTML ="Current Mode: UNLOCKED"
                            currentLockStatus = 0

                        }
                        sendLockUpdate(currentLockStatus)
                        break;
                    case 1:     // SOFT

                        console.log("SOFT")

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"
                            currentLockStatus = 1
                        }
                        else if (currentLockStatus === 2) {
                            currentLockStatus = 1
                        }
                        modeContainer.innerHTML ="Current Mode: SOFT"
                        sendLockUpdate(currentLockStatus)
                        break;
                    case 2:     // HARD

                        console.log("HARD")

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"

                            currentLockStatus = 2
                        } else if (currentLockStatus === 1) {
                            currentLockStatus =  2
                        }
                        modeContainer.innerHTML ="Current Mode: HARD"
                        sendLockUpdate(currentLockStatus)
                        break;
                }

            }
        }
    });

});


function sendUserSettings(value = 0) {
    let json;

    console.log("sending user settings")

    if (UNLOCKED) {
        json = JSON.stringify({currentLockStatus: 0});
    }
    if (SOFT) {
        json = JSON.stringify({currentLockStatus: 1});
    }

    else {
        json = JSON.stringify({currentLockStatus: 2});
    }


    axios({
        method: "post",
        url: "http://localhost:8080/getUserSettings",
        data: json,
        headers: {
            "Content-Type": "application/json",
        },

    })
        .then(function (response) {
            //handle success
            console.log(response);
        })
        .catch(function (response) {
            //handle error
            console.log(response);
        });
}
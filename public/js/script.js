"use strict"


console.log("hello world")

//TODO  send the lock (integer) status to /updateLock/:lockMode
function sendLockUpdate(newLockMode){
  // console.log(newLockMode.toString())

    let headers = {}
    if (localStorage.token) {
        headers = { 'Authorization': `Bearer ${localStorage.token}`, }
    }

    axios({
        method: "put",
        withCredentials: true,
        crossDomain: true,
        headers: headers,
        url: "/updateLock/"+ newLockMode.toString(),
        data: null,
    })

}


// logs time
function convertEpochToSpecificTimezone(timeEpoch, offset) {
    let d = new Date(timeEpoch);
    let utc = d.getTime() + (d.getTimezoneOffset() * 60000);  //This converts to UTC 00:00
    let nd = new Date(utc + (3600000 * offset));
    return nd.toLocaleString();
}

function DoorTable() {
    // console.log("asking for users")

    let headers = {}
    if (localStorage.token) {
        headers = { 'Authorization': `Bearer ${localStorage.token}`, }
    }

    axios({
        method: "get",
        xhrFields: {
            withCredentials: true
        },
        headers: headers,
        crossDomain: true,
        url: "/statistics/keycardUsed",
        data: null,
    })
        .then(function (response) {
            //handle success

            let string1 = JSON.stringify(response);
            let parsed = JSON.parse(string1);

           // console.log(parsed + "here is data")

            for (let i = 0; i < parsed.data.length; i++) {

                let tr = document.createElement('tr')
                let td1 = document.createElement('th')
                let td2 = document.createElement('td')
                let td3 = document.createElement('td')
                let text1 = document.createTextNode((i + 1).toString())
                let text2 = document.createTextNode(parsed.data[i].name)
                let text3 = document.createTextNode(convertEpochToSpecificTimezone(parsed.data[i].time, +3))
                td1.appendChild(text1)
                td2.appendChild(text2)
                td3.appendChild(text3)
                tr.appendChild(td1)
                tr.appendChild(td2)
                tr.appendChild(td3)


                document.getElementById("doorTable").appendChild(tr)

            }

        })
        .catch(function (response) {
            //handle error
            console.log(response);
        });

}

function LockTable() {
    // console.log("asking for users")

    let headers = {}
    if (localStorage.token) {
        headers = { 'Authorization': `Bearer ${localStorage.token}`, }
    }

    axios({
        method: "get",
        xhrFields: {
            withCredentials: true
        },
        crossDomain: true,
        headers: headers,
        url: "statistics/modeChanged",
        data: null,
    })
        .then(function (response) {
            //handle success

            let string1 = JSON.stringify(response);
            let parsed = JSON.parse(string1);

          //  console.log(parsed + "here is data")

            for (let i = 0; i < parsed.data.length; i++) {

                let tr = document.createElement('tr')
                let td1 = document.createElement('th')
                let td2 = document.createElement('td')
                let td3 = document.createElement('td')
                let td4 = document.createElement('td')
                let text1 = document.createTextNode((i + 1).toString())
                let text2 = document.createTextNode(parsed.data[i].mode)
                let text3 = document.createTextNode(parsed.data[i].name)
                let text4 = document.createTextNode(convertEpochToSpecificTimezone(parsed.data[i].time, +3))
                td1.appendChild(text1)
                td2.appendChild(text2)
                td3.appendChild(text3)
                td4.appendChild(text4)
                tr.appendChild(td1)
                tr.appendChild(td2)
                tr.appendChild(td3)
                tr.appendChild(td4)



                document.getElementById("lockTable").appendChild(tr)

            }

        })
        .catch(function (response) {
            //handle error
            console.log(response);
        });

}




function  updateUserVars() {


        let headers = {}
        if (localStorage.token) {
            headers = { 'Authorization': `Bearer ${localStorage.token}`, }
        }

    fetch("/getInitData", { headers: headers })
        .then(function(response) {
            return response.json();
        })
        .then(function(myJson) {
            console.log(JSON.stringify(myJson));
            document.getElementById("userName").innerText = "Welcome " + myJson.name
            document.getElementById("modeContainer").innerText = "Current Lock Mode: " + myJson.mode
        });


}




window.addEventListener('DOMContentLoaded', (event) => {

    document.getElementById("modeContainer").innerText ="new"
    updateUserVars()




    console.log("DOM LOADED")
    const lockImage = document.getElementById("lockImage")
    const lockRadio = document.getElementsByName("mode")
    const applyButton = document.getElementById("applyButton")
    const modeContainer = document.getElementById("modeContainer")

    const homeButton = document.getElementById("homeButton")
    const historyButton = document.getElementById("historyButton")

    const logoutButton = document.getElementById("logoutButton")


    const homeDiv = document.getElementById("homeContent")
    const DoorStatusDiv = document.getElementById("DoorStatusContent")



    homeDiv.style.display = "block"
    DoorStatusDiv.style.display = "none"



    //  Home button handling
    homeButton.addEventListener("click", () => {
        console.log("homeButton clicked.")
        homeDiv.style.display = "block"
        DoorStatusDiv.style.display = "none"

    });


    //  history button handling
    historyButton.addEventListener("click", () => {
        console.log("historyButton clicked.")
        homeDiv.style.display = "none"
        DoorStatusDiv.style.display = "block"

    });


    // Log Out
    logoutButton.addEventListener("click", () => {
        console.log("logout button pressed")
        localStorage.clear()
        window.location.href = "http://127.0.0.1:8080"
    });


    let currentLockStatus = {
        UNLOCKED: 0,
        SOFT: 1,
        HARD: 2,
    };


    DoorTable()
    LockTable()

    applyButton.addEventListener("click", function () {
        for (let i = 0; i < lockRadio.length; i++) {
            if (lockRadio[i].checked) {
              //  console.log(i)

                switch (i) {
                    case 0:     // OPEN

                        console.log("OPEN")
                        if (currentLockStatus !== 0) {
                            lockImage.src = "img/lockOpenAnim.png"
                            modeContainer.innerHTML ="Current Mode: UNLOCKED"
                            currentLockStatus = 0

                            console.log("changed")
                        }
                        sendLockUpdate(currentLockStatus)
                        break;
                    case 1:     // SOFT

                        console.log("SOFT")

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"
                            currentLockStatus = 1
                            console.log("changed")
                        }
                        else if (currentLockStatus === 2) {
                            currentLockStatus = 1
                            console.log("changed")
                        }
                        modeContainer.innerHTML ="Current Mode: SOFT"
                        sendLockUpdate(currentLockStatus)
                        break;
                    case 2:     // HARD

                        console.log("HARD")

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"

                            console.log("changed")

                            currentLockStatus = 2
                        } else if (currentLockStatus === 1) {
                            currentLockStatus =  2

                            console.log("changed")
                        }
                        modeContainer.innerHTML ="Current Mode: HARD"
                        sendLockUpdate(currentLockStatus)
                        break;
                }

            }
        }
    });

});

"use strict"

console.log("hello world")

//TODO  send the lock (integer) status to /updateLock/:lockMode
function sendLockUpdate(newLockMode){

}


window.addEventListener('DOMContentLoaded', (event) => {
    console.log("DOM LOADED")
    const lockImage = document.getElementById("lockImage")
    const doorImage = document.getElementById("doorImage")
    const lockRadio = document.getElementsByName("mode")
    const applyButton = document.getElementById("applyButton")
    const changeButton = document.getElementById("changeButton")
    const modeContainer = document.getElementById("modeContainer")
    const statusContainer = document.getElementById("statusContainer")

    const homeButton = document.getElementById("homeButton")
    const historyButton = document.getElementById("historyButton")
    const statusButton = document.getElementById("statusButton")

    const homeDiv = document.getElementById("homeContent")
    const historyDiv = document.getElementById("historyContent")
    const statusDiv = document.getElementById("statusContent")

    homeDiv.style.display = "block"
    historyDiv.style.display = "none"
    statusDiv.style.display = "none"

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



    //  login button handling
    statusButton.addEventListener("click", () => {
        console.log("StatusButton clicked.")
        homeDiv.style.display = "none"
        historyDiv.style.display = "none"
        statusDiv.style.display = "block"
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
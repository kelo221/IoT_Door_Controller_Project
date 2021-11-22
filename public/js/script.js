"use strict"

console.log("hello world")


window.addEventListener('DOMContentLoaded', (event) => {
    console.log("DOM LOADED")
    const lockImage = document.getElementById("lockImage")
    const lockRadio = document.getElementsByName("mode")
    const applyButton = document.getElementById("applyButton")
    const modeContainer = document.getElementById("modeContainer")

    const homeButton = document.getElementById("homeButton")
    const historyButton = document.getElementById("historyButton")
    const statusButton = document.getElementById("statusButton")

    const homeDiv = document.getElementById("homeContent")
    const historyDiv = document.getElementById("historyContent")
    const statusDiv = document.getElementById("statusContent")

    historyDiv.style.display = "none"
    statusDiv.style.display = "none"

    let currentLockStatus = 0


    //  Home button handling
    homeButton.addEventListener("click", () => {
        console.log("homeButton clicked.")
        homeDiv.style.display = "block"
        historyDiv.style.display = "none"
        statusDiv.style.display = "none"
    });

/*

    //  history button handling
    historyButton.addEventListener("click", () => {
        console.log("historyButton clicked.")
        homeDiv.style.display = "none"
        historyDiv.style.display = "block"
        statusDiv.style.display = "none"

    });
*/


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

                        if (currentLockStatus !== 0) {
                            lockImage.src = "img/lockOpenAnim.png"
                            modeContainer.innerHTML ="Current Mode: Open"
                            currentLockStatus = 0
                        }
                        break;
                    case 1:     // SOFT

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"
                            modeContainer.innerHTML ="Current Mode: Soft Lock"
                            currentLockStatus = 1
                        }
                        break;
                    case 2:     // HARD

                        if (currentLockStatus === 0) {
                            lockImage.src = "img/lockCloseAnim.png"
                            modeContainer.innerHTML ="Current Mode: Hard Lock"
                            currentLockStatus = 2
                        } else if (currentLockStatus === 1) {
                            currentLockStatus = 2
                        }
                        break;
                }

            }
        }
    });


});
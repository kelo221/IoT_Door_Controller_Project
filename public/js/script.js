"use strict"

console.log("hello world")


window.addEventListener('DOMContentLoaded', (event) => {
    console.log("DOM LOADED")
    const lockImage = document.getElementById("lockImage")
    const lockRadio = document.getElementsByName("mode")
    const applyButton = document.getElementById("applyButton")
    const modeContainer = document.getElementById("modeContainer")

    let currentLockStatus = 0

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

                        if (currentLockStatus !== 1) {
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
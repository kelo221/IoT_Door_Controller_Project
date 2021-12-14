# IoT Door Project
## About

A Fullstack for door control system, using an Arduino, NFC, and a servo motor.
The web interface controls the lock state of the door and passes requests to the backend to send and receive information.

<img width="550" alt="截屏2021-12-14 02 29 25" src="https://user-images.githubusercontent.com/56063237/145910718-a256f146-35a1-4dfb-8ebc-9c43f8bedfbe.png">



## 1. Authentication page

<img width="1177" alt="截屏2021-12-14 02 27 33" src="https://user-images.githubusercontent.com/56063237/145910625-f2c29b7a-c1b1-46bd-8ba4-dc7ca4b94511.png">


The first time connecting to the website the user is prented with a login prompt.
Entering "User" for both password and username the server replies with an JWT authenctication key.

## 2. Home Page

<img width="1065" alt="截屏2021-12-14 02 25 45" src="https://user-images.githubusercontent.com/56063237/145910447-1f6ae905-1b43-43ab-b4bd-4c813c9b528b.png">

#### 1 .User name
Here shows the username.

#### 2. Lock mode
Here shows the current lock mode.

#### 3. Animation shows the current mode
if user changes the mode, there will show the animation to the mode user select.

#### 4. Users’ setting
Here are 3 modes user can choose, and each choice has the info icon to explain different mode.

#### 5. Apply setting
save the changes

#### 6. Log out
clear the data.



## 3. History page
logs pic

Login logs are stored here.


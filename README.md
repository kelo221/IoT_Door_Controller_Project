# IoT Door Project
## About

A Fullstack for door control system, using an Arduino, NFC, and a servo motor, the door opens when either the GridEye sees something or the lock requirements are met.
The web interface controls the lock state of the door and passes requests to the backend to send and receive information.
REST API documentation can be found at http://localhost:8080/swagger/index.html#/ , once the server is running.

<img width="565" alt="截屏2021-12-16 22 03 16" src="https://user-images.githubusercontent.com/56063237/146441339-f56d6abb-c3ba-4deb-ae5f-6b517e92e5ed.png">

![Untitled Diagram drawio](https://user-images.githubusercontent.com/56063237/146441074-f5c40f9c-63c3-4f5e-b0fe-fff94e4ff140.png)



## 1. Authentication page

<img width="1177" alt="截屏2021-12-14 02 27 33" src="https://user-images.githubusercontent.com/56063237/145910625-f2c29b7a-c1b1-46bd-8ba4-dc7ca4b94511.png">


The first time connecting to the website the user is prented with a login prompt.
Entering "User" for both password and username the server replies with an JWT authenctication key.

## 2. Home Page

<img width="1072" alt="截屏2021-12-15 23 31 09" src="https://user-images.githubusercontent.com/56063237/146269003-fb5208a5-5a98-4299-bedd-480bdec7aefc.png">

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

#### 7. Open door manually
when the mode is hard lock, user will need to open the door manually.




## 3. History page
<img width="1035" alt="截屏2021-12-15 23 34 18" src="https://user-images.githubusercontent.com/56063237/146268494-e68a0b49-88e4-4b64-8d28-f92acefe413d.png">


Login logs are stored here.

#### 1. Read NFC
The first table shows the NFC history. 

#### 2. Mode change history
everytime the door lock mode changes, the server will send to database and saved the history. The table shows the user's name, mode changed and time.



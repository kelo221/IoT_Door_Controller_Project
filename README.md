# IoT Door Project
## About

A Fullstack door control system, using an Arduino, NFC, and a servo motor. The door opens when either the GridEye sensor detects a rise in temperature or the lock requirements are met.
The web interface controls the lock state of the door and passes requests to the backend in order to send and receive information.
REST API documentation can be found at http://localhost:8080/swagger/index.html#/ , once the server is running.

![Door drawio](https://user-images.githubusercontent.com/56063237/146462973-2939beef-df30-462d-ac75-b00a9305da29.png)


![Untitled Diagram drawio](https://user-images.githubusercontent.com/56063237/146441074-f5c40f9c-63c3-4f5e-b0fe-fff94e4ff140.png)



## 1. Authentication page

<img width="1177" alt="截屏2021-12-14 02 27 33" src="https://user-images.githubusercontent.com/56063237/145910625-f2c29b7a-c1b1-46bd-8ba4-dc7ca4b94511.png">


If it is the user's first time connecting to the website, the user is prented with a login prompt.
By entering "User" for both password and username, the server replies with an JWT authenctication key.

## 2. Home Page

<img width="1072" alt="截屏2021-12-15 23 31 09" src="https://user-images.githubusercontent.com/56063237/146269003-fb5208a5-5a98-4299-bedd-480bdec7aefc.png">

#### 1 .User name
The username is displayed.

#### 2. Lock mode
Here, the current lock mode is shown.

#### 3. Animation shows the current mode
By changing the mode, an animation will be displayed to show the selected mode.

#### 4. Users’ setting
There are three differnt lock modes available. Each mode has an info icon attached, which explains the behaviour of the mode.

#### 5. Apply setting
Saves the changes.

#### 6. Log out
Clears the data.

#### 7. Open door manually
When the mode is hard lock, the user will need to open the door manually.




## 3. History page
<img width="1035" alt="截屏2021-12-15 23 34 18" src="https://user-images.githubusercontent.com/56063237/146268494-e68a0b49-88e4-4b64-8d28-f92acefe413d.png">


Login logs are stored here.

#### 1. Read NFC
The first table shows the NFC history. 

#### 2. Mode change history
Everytime the door lock mode changes, the server send this information to the database and save in in the history. The table shows the user's name, mode changes and time.



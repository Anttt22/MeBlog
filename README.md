# App MeBlog  

App for writing and viewing posts  

Introduction:    
App for viewing posts and adding posts  

Main page for viewing posts  
Login page for entering credentials   
Create post page for creaing posts in timeline for autorized users  

Tech:  
Go, html  

Authentication - bcrypt for hashing passwords  
Authorization - Jwt in cookie with an expiration time 2 min  
db - MySQL, db containd 2 tables: users and articles  


Launch  

go run main.go  

login: den  
password: supersecet2  
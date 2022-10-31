# Response programe
I used Domain Driven Design (DDD) for this project. it's a fantastic architecture that focuses on business logic.

The business domain doesn't care what framework you use for controllers or the databases to which the repository is hooked.


The separation of concern is quite resourceful because it's made possible by the use of interfaces, which means you can switch frameworks or databases very easily.

## how to run the app

the easiest way to check the functionality of the load test visit 

https://github.com/myrachanto/loader/actions - github actions

and check the test 

## or

depending on how you have set up your machine 

you can use either to pull the code to your local machine

@ ssh git clone git@github.com:myrachanto/loader.git

or

https git clone https://github.com/myrachanto/loader.git

then you can navigate to the project folder run "go run main.go"

"I am assuming you have golang already set in you machine!"

# usage  
## make sure "," separates the urls
go run . getulr -urls "https://www.google.com,https://www.chantosweb.co.ke"

 ### or
 go run main getulr -urls "https://www.google.com,https://www.chantosweb.co.ke"


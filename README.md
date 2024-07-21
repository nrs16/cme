# Coding Challenge task

## Some Disclaimers
- I assumed that messaging should be done among users registered on the platform, and for the sake of this assingment I did one to one chat only
- I did not put any restrictions on password , for the sake of ease of use
- The Database schema on some of the code are meant to do beyond the scope of the assumptions above, they can ce extended for group chatting, but this is not fully done. So if you see something that is not fullt needed or used, it's for this reason, disregard it
- I did some quick documention in ```/api_specs/api_specs.yaml```, I suggest you copy those into 3.0 swaager editor and use them as reference



## How to get and run the code
### Through github

- clone the project using ```git clone https://github.com/nrs16/echo-challenge.git```
- run the command ```go mod download```
- run the command ```go run main.go```
- use this curl to test the code: 

```
curl --location 'http://localhost:8080/routes' \
--header 'x-correlation-id: jweygfjkegdf' \
--header 'Content-Type: application/json' \
--data '[["LAX","DXB"],["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]'

```
You can remove x-correlation-id and Content-Type headers but you must send the body


### Through docker

- get the image using ```git pull nrs16/echoserver```
this might take a while to download because it has golang:1.21 image as base

- run the image inside a container using ```docker run -p 8000:8080 --name routes nrs16/echoserver```
    - Note that you can change port 8000 to whichever port you want on your machine.
    - To run the container in the background add -d to the commad so ```docker run -d -p 8000:8080 --name routes nrs16/echoserver```
- use the below curl to test the code:
```
curl --location 'http://localhost:8000/routes' \
--header 'x-correlation-id: jweygfjkegdf' \
--header 'Content-Type: application/json' \
--data '[["LAX","DXB"],["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]'

```
You can remove x-correlation-id and Content-Type headers but you must send the body
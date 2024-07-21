# Coding Challenge task

## Some Disclaimers
- jwt encryption key is harcoded in code, this should be moved and read from secret file.
- I assumed that messaging should be done among users registered on the platform, and for the sake of this assingment I did one to one chat only fully
- I did not put any restrictions on password , for the sake of ease of use
- The Database schema on some of the code are meant to do beyond the scope of the assumptions above, they can be extended for group chatting.
- API documention is in ```/api_specs/api_specs.yaml```, I suggest you copy those into 3.0 swagger editor and use them as reference
- register 2 users so you can send message from one to the other
- database initialization should be improved, and migrations should be implemented, but I did the simplest thing for the task
- database and redis persistance volumes can be added
- for sake of task no authentication done on redis and cassandra

## Deployment
### Docker compose

- clone the **main** branch (to access the docker-compose.yml file), if you don't want to clone the project you can just download the docker-compose.yml file from github
- go to project directory / if you only downloaded the docker-compose.yml file go to the directory that contains the file
- run ```docker-compose up``` Note that this will take a couple of minutes to pull the images, if you have connectivity issues, open the docker-compose file and pull each individually the run ```docker-compose up``` again

### system run

if you want to run project locally ,you also need to have cassandra and redis running:
- clone the **main** branch
- change config hosts for redis and cassandra to "localhost"
- make
- ./main


- Service should be good to go on port 8099: you can use the below curl to test 
also refer to documentation in ```/api_apecs/api_specs.yaml```

```
curl --location 'http://localhost:8099/api/v1/register' \
--data-raw '{
    "username":"jh0",
    "password":"Testing@123",
    "first_name":"John",
    "last_name":"Doe",
    "email_address":"johndoe@mail.comm"
}
'
```

```
curl --location 'http://localhost:8099/api/v1/message' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5yczEifQ.slPVn9nJOtUER8Eoevcr2ueuLqmloomLkBKtcPcYv_A' \
--header 'Content-Type: application/json' \
--data '{
    "to_id": "nrs16",
    "message": "Holaaaa nrs1666666"
}'
```


## Architercture
### Database
refer to ```schema.sql``` for full schema
- **user** : this will contain user information and used to authenticate,includes password salt and hash, I kept as simple as possible for now
- **chat** table: for one on one chats: id, participant_1, participant_2, ts_created
- **group_chat**: table for group chats (the API for this is not implemented fully on API level): id, participants, ts_created
- **message** table: id, chat_id, from_id, message, ts_created
    message is linked to chat, and contains the sender username.


### API

- first you need to register by providing your info, username and password, this will return a token that you can use to send to retrieve message. I checked cached usernames from redis, to make sure of uniqueness
- you can also login, this will also return a token to be used on message sending an retrieval
- to send a message the user should:
    -  send recipient username in payload  
    ```
        {
            "to_id": "nrs16",
            "message": "Holaaaa nrs1666666"
        }
    ```
    OR
    - send chatId in payload  
    ```
        {
            "chat_id": "02887379-7b6f-49fb-94e7-cb8df1e21555",
            "message": "Holaaaa nrs1666666"
        }
    ```

    if recipient username is sent, it means this is a one to one message.
    I check db to see if a chat exists between these 2.
        if a chat exists: I insert the message and link it to the chat
        if a chat does not exist: I create the chat and link the message to it

    if a chatid is sent, I just insert the message and link it to chat id (so chatid  works for group and one to one messaging)

- to view messages, the client needs to get the chats first using ```get /chat```,(I get user chat_id from redis to avoid filtering on db) and the retrieve the messages of that chat using   ```get /chat/{chatId}/message```
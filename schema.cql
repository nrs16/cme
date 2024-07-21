CREATE KEYSPACE authentication WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};
create table authentication.user( username VARCHAR PRIMARY KEY, first_name VARCHAR, last_name VARCHAR, password_hash VARCHAR, password_salt VARCHAR, email VARCHAR, ts_created TIMESTAMP );
create table authentication.chat ( id UUID, participant_1  VARCHAR, participant_2  VARCHAR, ts_created TIMESTAMP, PRIMARY KEY (id,participant_1, participant_2));
create table authentication.groupchat (id UUID PRIMARY KEY,participants LIST<VARCHAR>,ts_created TIMESTAMP);
create table authentication.message(  id UUID, chat_id UUID, from_id VARCHAR,  message VARCHAR, ts_created TIMESTAMP, PRIMARY KEY (chat_id, ts_created, id))WITH CLUSTERING ORDER BY (ts_created DESC);


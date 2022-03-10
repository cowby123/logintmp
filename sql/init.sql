drop database uploaddb;
create database uploaddb;
ALTER DATABASE uploaddb SET timezone TO 'Asia/Shanghai';
\c uploaddb
CREATE TABLE UserData(
    ID  SERIAL PRIMARY KEY,
    UserName TEXT NOT NULL,   
    Password TEXT NOT NULL,   
    ChineseName TEXT NOT NULL,   
    Email TEXT NOT NULL,   
    Address TEXT NOT NULL,   
    Phone TEXT NOT NULL,   
    State INT NOT NULL,
    LastLogin TEXT NOT NULL);


INSERT INTO UserData(UserName,Password,ChineseName,Email, Address ,Phone,  State,LastLogin)

VALUES('cowby123','EBAF985E5D8380FD374351EE89752EC5','腿骨','cowby123@gmail.com','嘉義縣民雄鄉','0912345678',0,now());


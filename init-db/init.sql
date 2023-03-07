CREATE DATABASE main_db;
USE main_db;
CREATE TABLE articles(
   ID   VARCHAR (40) NOT NULL,
   TITLE VARCHAR (20) NOT NULL,
   CONTENT  VARCHAR (100),
   AUTHOR   VARCHAR (20),       
   PRIMARY KEY (ID)
);
INSERT INTO articles VALUES 
(1, "Title 1", "lorem ipsum some content etcetc", "author 1"),
(2, "Title 2", "lorem ipsum some content etcetc", "author 2"),
(3, "Title 3", "lorem ipsum some content etcetc", "author 3");

CREATE DATABASE test_db;
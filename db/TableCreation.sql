# Order Matters
DROP TABLE IF EXISTS Body, Food, Snakes, MoveReq, Games;

CREATE TABLE Games
(
id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
gameid CHAR(36) UNIQUE,
width INT,
height INT
);

CREATE TABLE MoveReq
(
id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
turn int,
g_id int,
FOREIGN KEY (g_id) REFERENCES Games(id) ON DELETE CASCADE
);

CREATE TABLE Snakes
(
id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
name_ varchar(15),
health int,
m_id int,
len int,
FOREIGN KEY (m_id) REFERENCES MoveReq(id) ON DELETE CASCADE
);


CREATE TABLE Food
(
x int,
y int,
m_id int,
FOREIGN KEY (m_id) REFERENCES MoveReq(id) ON DELETE CASCADE,
PRIMARY KEY (x,y,m_id)
);

CREATE TABLE Body
(
x int,
y int,
m_id int,
FOREIGN KEY (m_id) REFERENCES MoveReq(id) ON DELETE CASCADE,
PRIMARY KEY (x,y,m_id)
);






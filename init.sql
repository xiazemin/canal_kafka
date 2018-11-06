GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
FLUSH PRIVILEGES;
create database test;
use test;
CREATE TABLE `xdual` (     ->   `ID` int(11) NOT NULL AUTO_INCREMENT,     ->   `X` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,     ->   PRIMARY KEY (`ID`)     -> ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
insert into xdual(id,x) values(null,now());
update xdual set x='2018-11-06 20:38:43' where id<10;
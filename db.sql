Create database EG;
Use EG

DROP TABLE IF EXISTS users;
create table users(
  student_id   varchar(100)  not null UNIQUE,
  password     varchar(100)  not null ,
  user_picture varchar(100)  not null ,
  name         varchar(100)  not null ,
  summary      varchar(1999) not null ,
  sex          int           not null ,
  god          int           not null ,
  primary key (student_id)
)ENGINE=InnoDB;

DROP TABLE IF EXISTS backpads;
create table backpads(
  id         int          not null auto_increment ,
  student_id varchar(100) not null ,
  name       varchar(100) not null ,
  time       datetime     not null ,
  state      int          not null ,
  hours      int ,
  minutes    int ,       
  primary key (id)
)ENGINE=InnoDB;


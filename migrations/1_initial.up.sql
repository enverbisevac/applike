create table if not exists todo_item(
    id int auto_increment primary key,
    description varchar(100) not null,
    due_date datetime not null
);
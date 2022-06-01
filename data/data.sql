create table if not exists products (
  id varchar(40) not null,
  productName varchar(120),
  description varchar(120),
  price varchar(45),
  status varchar(45),
  primary key (id)
);

insert into products (id, productName, description, price, status) values ('P001', 'Iron Man', 'toys', '1000', 'available');
insert into products (id, productName, description, price, status) values ('P002', 'Scram411', 'bike', '2000', 'available');
insert into products (id, productName, description, price, status) values ('P003', 'Ikea 4025', 'furniture', '3000', 'not available');

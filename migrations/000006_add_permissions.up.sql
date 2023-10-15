create table if not exists premissions (
	id bigserial primary key,
	code text not null
);

create table if not exists user_premissions (
	user_id bigint not null references users on delete cascade,
	premission_id bigint not null references premissions on delete cascade,
	primary key (user_id, premission_id)
);

insert into premissions (code)
values
	('movies:read'),
	('movies:write');

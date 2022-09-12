create table tb_endereco
(
    id         integer      not null
        constraint pk_id_tb_entedereco
            primary key,
    logradouro varchar(255) not null,
    bairro     varchar(255),
    cep        integer      not null,
    cidade     varchar(255),
    uf         varchar(2)   not null
);

alter table tb_endereco
    owner to postgres;



create table tb_cliente
(
    id            integer      not null
        constraint pk_id_tb_cliente
            primary key,
    primeiro_nome varchar(100) not null,
    sobrenome     varchar(255),
    endereco_id   integer
        constraint fk_endereco_id_tb_cliente
            references tb_endereco
            on update cascade on delete cascade
);

alter table tb_cliente
    owner to postgres;


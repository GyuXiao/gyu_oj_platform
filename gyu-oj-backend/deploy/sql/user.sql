-- 创建库
create database if not exists oj_db;

-- 切换库
use oj_db;

-- 用户表
create table if not exists `user`
(
    `id`         bigint                                                                                         not null auto_increment comment '唯一 id' primary key,
    `username`   varchar(256)                                                                                   not null comment '用户昵称',
    `password`   varchar(512)                                                                                   not null comment '用户密码',
    `avatarUrl`  varchar(1024) default 'https://gyu-pic-bucket.oss-cn-shenzhen.aliyuncs.com/gyustudio_icon.jpg' not null comment '用户头像',
    `email`      varchar(256)                                                                                   null comment '用户邮箱',
    `phone`      varchar(256)                                                                                   null comment '手机号',
    `userRole`   tinyint       default 0                                                                        not null comment '用户角色 0 - 普通用户 1 - 管理员',
    `isDelete`   tinyint       default 0                                                                        not null comment '是否删除 0 - 未删除 1- 删除',
    `createTime` datetime      default CURRENT_TIMESTAMP                                                      not null comment '创建时间',
    `updateTime` datetime      default CURRENT_TIMESTAMP                                                      not null on update CURRENT_TIMESTAMP comment '更新时间',
    UNIQUE KEY `username` (`username`)
) comment '用户表';

-- 插入初始用户
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('方志强', '0r4', 'www.grazyna-feest.org', 'lynna.davis@hotmail.com', '15009759625', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('周瑞霖', '1yJ', 'www.angel-rice.io', 'rodrick.ruecker@hotmail.com', '17767824701', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('赖博超', 'ePGc', 'www.annamaria-koepp.co', 'mercy.lynch@gmail.com', '17385852712', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('于嘉熙', 'BOrlI', 'www.mitsue-pollich.io', 'javier.donnelly@gmail.com', '18081785994', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('何伟泽', 'BdD', 'www.cheryle-white.org', 'erasmo.medhurst@yahoo.com', '18368755044', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('尹炎彬', 'pl', 'www.anthony-dickens.com', 'myrtle.bins@hotmail.com', '15008868928', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('陶靖琪', 'L8', 'www.stephan-harris.net', 'hisako.johnston@yahoo.com', '17745109837', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('王立果', 'Cm2M', 'www.norberto-hessel.co', 'anton.mills@gmail.com', '14792623727', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('侯晓博', 'r40c', 'www.araceli-collins.com', 'tamie.bashirian@yahoo.com', '15663192241', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('莫致远', 'ds', 'www.perry-quigley.net', 'raymond.klein@yahoo.com', '17610814290', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('陶昊然', '2bSwI', 'www.laurence-towne.co', 'vincenzo.mraz@hotmail.com', '17111179430', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('王思聪', 'PNS', 'www.arlene-wolf.name', 'maritza.turcotte@yahoo.com', '13626299256', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('崔文', 'dQbZT', 'www.quinton-johnston.org', 'dong.osinski@hotmail.com', '13365490026', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('熊昊天', '49', 'www.jospeh-considine.net', 'carl.haley@gmail.com', '15982953501', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('崔靖琪', 'hPY', 'www.ellis-barton.com', 'zachary.hoppe@gmail.com', '13196148866', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('张弘文', 'g5', 'www.robena-stanton.net', 'rod.mayert@gmail.com', '15623773067', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('黎鹏煊', 'H7e2z', 'www.juliet-stroman.info', 'tomas.paucek@gmail.com', '17038978569', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('贾语堂', 'xhV50', 'www.adam-cormier.co', 'aisha.konopelski@yahoo.com', '17840798968', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('林浩宇', 'gScm', 'www.jenelle-kilback.io', 'lorie.west@hotmail.com', '15311512259', 0);
insert into `user` (`username`, `password`, `avatarUrl`, `email`, `phone`, `userRole`) values ('秦胤祥', '1T', 'www.marvin-kirlin.name', 'terence.abernathy@gmail.com', '17298874783', 0);

# ALTER TABLE `user` DROP INDEX `phone`;
# ALTER TABLE `user` DROP INDEX `email`;

#
# ALTER TABLE `user` ADD INDEX `phone` (`phone`);
# ALTER TABLE `user` ADD INDEX `email` (`email`);

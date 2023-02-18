/* 基本的用户表*/
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `user_id` bigint(20) NOT NULL,
                        `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_username` (`username`) USING BTREE,
                        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*用户信息表*/
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`(
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL,
    `follow_num` smallint NOT NULL DEFAULT 0 COMMENT '关注的数量',
    `fans_num` smallint NOT NULL DEFAULT 0 COMMENT '粉丝的数量',
    `praise_num` INT NOT NULL DEFAULT 0 COMMENT '得到赞的数量',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`)USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '记录用户关注和喜欢的数量表';

/*用户的 关注/粉丝 列表*/
DROP TABLE IF EXISTS `user_follow`;
CREATE TABLE `user_follow`(
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `user_id` bigint(20) NOT NULL COMMENT '用户id',
                              `follower_id` bigint(20) NOT NULL COMMENT '关注的用户id',
                              `is_follow` bool NOT NULL DEFAULT false COMMENT 'false不关注，true关注',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `idx_user_follower` (`user_id`,`follower_id`)USING BTREE,
                              KEY `idx_follower_id` (`follower_id`)USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '用户之间的关系表';

/*视频表*/
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`(
                         `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
                         `author_id` bigint(20) NOT NULL COMMENT '视频作者id',
                         `play_url` varchar(255) NOT NULL COMMENT '播放url',
                         `cover_url` varchar(255) NOT NULL COMMENT '封面url',
                         `publish_time` datetime NOT NULL COMMENT '发布时间戳',
                         `title` varchar(255) DEFAULT NULL COMMENT '视频标题',
                         PRIMARY KEY (`id`),
                         KEY `idx_time`(`publish_time`) USING BTREE,
                         KEY `idx_author`(`author_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '视频表';

/*消息表*/

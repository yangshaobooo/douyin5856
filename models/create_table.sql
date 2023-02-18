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


/*基本的用户信息表*/
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
                         `id` bigint(20) NOT NULL COMMENT '视频唯一id',
                         `author_id` bigint(20) NOT NULL COMMENT '视频作者id',
                         `play_url` varchar(255) NOT NULL COMMENT '播放url',
                         `cover_url` varchar(255) NOT NULL COMMENT '封面url',
                         `publish_time` datetime NOT NULL COMMENT '发布时间戳',
                         `favorite_count` bigint(20) NOT NULL  COMMENT '视频的点赞数量',
                         `comment_count` bigint(20) NOT NULL  COMMENT '视频评论的数量',
                         `title` varchar(255) DEFAULT NULL COMMENT '视频标题',
                         PRIMARY KEY (`id`),
                         KEY `idx_time`(`publish_time`) USING BTREE,
                         KEY `idx_author`(`author_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '视频表';

/*用户点赞视频表*/
DROP TABLE IF EXISTS `user_favorite_video`;
CREATE TABLE `user_favorite_video`(
                                      `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                      `user_id` bigint(20) NOT NULL COMMENT '用户id',
                                      `video_id` bigint(20) NOT NULL COMMENT '视频id',
                                      `is_favorite` bool NOT NULL DEFAULT false COMMENT 'false不喜欢，true喜欢',
                                      PRIMARY KEY (`id`),
                                      UNIQUE KEY `idx_user_video` (`user_id`,`video_id`)USING BTREE,
                                      KEY `idx_user` (`user_id`)USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '用户喜不喜欢视频';

/*评论表*/
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`(
                           `id` bigint(20) NOT NULL COMMENT '评论id',
                           `user_id` bigint(20) NOT NULL COMMENT '评论发布用户id',
                           `video_id` bigint(20) NOT NULL COMMENT '评论视频id',
                           `comment_text` varchar(255) NOT NULL COMMENT '评论内容',
                           `create_date` datetime NOT NULL COMMENT '评论发布时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_videoId`(`video_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '评论表';

/*消息表*/
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`(
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '消息主键，自增',
                           `to_user_id` bigint(20) NOT NULL COMMENT '消息接受者的id',
                           `from_user_id` bigint(20) NOT NULL COMMENT '消息发送者id',
                           `content` VARCHAR(255) NOT NULL COMMENT '消息内容',
                           `create_time` datetime NOT NULL COMMENT '消息创建时间',
                           PRIMARY KEY (`id`),
                           KEY `idx_to_fromId`(`to_user_id`,`from_user_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '消息表';
1 -- 用户表：存储账号及头像路径CREATETABLE users (
2     id INTEGERPRIMARY KEY AUTOINCREMENT,
3     username TEXT NOTNULLUNIQUE,
4     password TEXT NOTNULL,          -- 实际开发中应存储 Hash 后的密文
5     avatar_url TEXT,                 -- 头像在服务器的相对路径
6     created_at DATETIME DEFAULTCURRENT_TIMESTAMP
7 );
8 
9 -- 相册表：用户创建的文件夹CREATETABLE albums (
10     id INTEGERPRIMARY KEY AUTOINCREMENT,
11     user_id INTEGERNOTNULL,
12     name TEXT NOTNULL,
13     description TEXT,
14     is_public INTEGERDEFAULT0,     -- 0: 私有, 1: 公有
15     created_at DATETIME DEFAULTCURRENT_TIMESTAMP,
16     FOREIGN KEY (user_id) REFERENCES users(id)
17 );
18 
19 -- 照片表：相册内的具体图片CREATETABLE photos (
20     id INTEGERPRIMARY KEY AUTOINCREMENT,
21     album_id INTEGERNOTNULL,
22     file_path TEXT NOTNULL,         -- 图片存储路径
23     file_size INTEGER,               -- 单位: Byte
24     created_at DATETIME DEFAULTCURRENT_TIMESTAMP,
25     FOREIGN KEY (album_id) REFERENCES albums(id)
26 );
# coze-chat-proxy

## 接口

#### /v1/chat/completions

###### 支持返回stream和json

```
http://<ip>:<port>/v1/chat/completions
```

##### 示例

```
curl --location --request POST 'http://127.0.0.1:8080/v1/chat/completions' \
--header 'Authorization: Bearer ****' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--data-raw '{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "西红柿炒钢丝球怎么做？"
        }
    ],
    "stream": false,
}'
```

#### /v1/images/generations

###### 仅支持返回json

```
http://<ip>:<port>/v1/images/generations
```

##### 示例

```
curl --location --request POST 'http://127.0.0.1:8080/v1/images/generations' \
--header 'Authorization: Bearer ****' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--data-raw '{
    "model": "dall-e-3",
    "prompt": "A cute dog",
    "n": 1,
    "size": "1024x1024"
}'
```

###### 注：此接口coze bot 需特殊设置，参考 [how_to_create_coze_agent](https://github.com/Feiyuyu0503/free-dall-e-proxy/blob/main/docs/how_to_create_coze_agent.md) 此配置稳定出图

## 如何使用

1. 打开 [discord开发者平台](https://discord.com/developers/applications) 。
2. 创建bot-01,并记录bot专属的`token`和`id(COZE_BOT_ID)`,此bot为被coze托管的bot。
3. 创建bot-02,并记录bot专属的`token(BOT_TOKEN)`,此bot为向bot-01发送与接收bot-01返回的消息。
4. 两个bot开通对应权限(`Send Messages`,`Read Message History`等)并邀请进服务器,记录服务器ID(`GUILD_ID`) 。
5. 打开 [coze官网](https://www.coze.com) 创建自己bot。
6. 创建好后推送，配置discord-bot的`token`,即bot-01的`token`,点击完成后在discord的服务器中可看到bot-A在线并可以@使用。
7. 配置环境变量，并启动本项目。
8. 访问接口地址即可开始调试。

## 配置

### 环境变量

1. `LOG_LEVEL：info`  日志等级

2. `SERVER_PORT：8080`  服务端口

3. `BOT_CONFIG：bot.json` discord-bot配置文件

4. `AUTH_TOKEN`:`123456` 请求头校验的值（前后端统一）,配置此参数后，每次发起请求时请求头加上`Authorization`
   参数，即`header`中添加 `Authorization：Bearer 123456`

###### 也可使用与程序同目录下 `.env` 文件配置上述字段


### bot.json配置

###### 支持多个机器人，添加到json列表即可，同模型（model）轮训请求

```
[
  {
    "model": "gpt-3.5-turbo", // bot模型可自定义，与请求接口保持一致 可自定义
    "bot_token": "MTI************", //见如何使用.3
    "coze_bot_id": "120**********", //见如何使用.2
    "guild_id": "103************", //见如何使用.4
    "channel_id": "120*********" //在所在服务器 创建频道 记录id
  },
  {
    "model": "dall-e-3", // bot模型可自定义，与请求接口保持一致 可自定义
    "bot_token": "MTI************", //见如何使用.3
    "coze_bot_id": "120**********", //见如何使用.2
    "guild_id": "103************", //见如何使用.4
    "channel_id": "120*********" //在所在服务器 创建频道 记录id
  }
]
```


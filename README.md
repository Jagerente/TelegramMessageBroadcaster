## Information

Telegram bot for message broadcasting in different languages across chats.
The service uses PostgreSQL database to store the list of chats, languages, groups, and administrators.
The list of commands and interaction is available only to administrators.
Data management is done through commands, which are only accessible to IsMaster administrators.
IsMaster status is granted only through the database.
Broadcast management is available to administrators.

## Broadcast configuration

Messages for broadcasting are configured by sending the bot a message with the pattern ${msg_id;lang_id}, where:
- `msg_id` - Message ID (positive integer);
- `lang_id` - Language ID (`/listlanguage` to list existing languages)

To edit a message, simply send a new message with the necessary data in the pattern.

### Commands and examples

#### Add admin

```
# Input
/addadmin

# Output
Input user_id;user_name

# Input
123456789;Username

# Output
User [123456789] Username has been added!
```

#### List all admins

```
# Input
/listadmin

# Output
Admins List:
 [123456789] Username1
 [123456780] Username2
```

#### Remove admin

```
# Input
/removeadmin

# Output
Input user_id

# Input
123456789

# Output
User 123456789 has been removed!
```

#### Add language

```
# Input
/addlanguage

# Output
Input language_name

# Input
English

# Output
English has been added!
```

#### List all languages

```
# Input
/listlanguage

# Output
Language List:
 [1] English
 [2] Russian
 [3] Spain
```

#### Remove language

```
# Input
/removelanguage

# Output
Input language_id

# Input
1

# Output
Language 1 has been removed!
```

#### Add group

```
# Input
/addgroup

# Output
Input group_name

# Input
Group1

# Output
Group1 has been added!
```

#### List all groups

```
# Input
/listgroup

# Output
Group List:
 [1] Group1
 [2] Group2
```

#### Remove group

```
# Input
/removegroup

# Output
Input group_id

# Input
2

# Output
Group 2 has been removed!
```

#### Add chat

```
# Input
/addchat

# Output
Input chat_id;chat_name;language_id;group_id

# Input
-123456789;Chat 1;1;1

# Output
Russian Test has been added!
```

#### List all chats

```
# Input
/listchat

# Output
Chat List:
 [-123456789] Chat 1
 [-123456780] Chat 2
```

#### Remove chat

```
# Input
/removechat

# Output
Input chat_id

# Input
-123456789

# Output
-123456789 has been removed!
```

#### Configure for Broadcasting Message with ID 1 for language ID 1

```
# Input
${1;1}
English language message

# Output
Message stashed
Message ID: 1
Language: [1] English
```

#### Edit for Broadcasting Message with ID 1 for language ID 1

```
# Input
${1;1}
New english language message

# Output
Message stashed
Message ID: 1
Language: [1] English
```

#### Add for Broadcasting Message with ID 1 for language ID 2

```
# Input
${1;2}
Сообщение на русском языке

# Output
Message stashed
Message ID: 1
Language: [2] Russian
```

#### Preview Broadcasting message

```
# Input
/testmessage

# Output
Input message_id

# Input
1

# Output 1 
New english language message

# Output 2
Сообщение на русском языке
```

#### Broadcast message

```
# Input
/sendmessages

# Output
Input message_id;group_id

# Input
1;2

# Output
2 messages sent to group [2] Group2.
```

## Build

To run the application, .env file or environment variables is required.

```
# App related + Docker partial related
TG_TOKEN=1234:ABCDEFG         # Telegram bot token
TG_MASTER_ID=123456789        # Telegram master id
# POSTGRES_HOST=localhost     # Database host
# POSTGRES_PORT=5432          # Database port
POSTGRES_HOST=localhost       # Database host
POSTGRES_PORT=8310            # Database port
POSTGRES_USER=postgres        # Database username
POSTGRES_PASSWORD=83          # Database password
POSTGRES_DB=database          # Database name
DEBUG=true                    # Debug. Affects logging.

# Docker related
POSTGRES_PORT_OUT=8310              # Database port external
DOCKER_IMAGE=username/imagename:tag # TG bot image
```

Building
```
go build ./cmd/main.go
```
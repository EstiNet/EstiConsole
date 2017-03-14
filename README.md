#EstiConsole
EstiNet's Minecraft server console wrapper with remote access support.

#Protocol
Communicating with EstiConsole is rather easy. Use the Socket.io library to initialize a connection.

##What EstiConsole recieves:

###hello [password]
Must be sent before any other messages are sent. Will return error 401 if the password is incorrect.

###curlogs
EstiConsole will return all the logs up to that point in bytes.

###command [command]
EstiConsole will execute the command onto the server.

###curdir [directory]
Request the file list of the directory. "./" requests the root directory.

###upload [directory] [bytestream]
Uploads a file to the directory.

###mkdir [directory]
Creates a folder at the directory.

###delete [directory]
Deletes a file at the directory.

###download [directory]
Downloads a file at the directory.

##What EstiConsole sends:

###authed
Acknowledges that your client has been authenticated after a hello query.

###curlogs [bytes]
Sent after the client requests the logs with curlogs. Returns bytes containing all the logs.

###log [bytes]
Sent when there there is console output.

###ecerror [error code]
Sent back if something went wrong with input. Check below for what the error code means.

###curdir [filename:sizekb:ifdir filename2:size2kb:ifdir etc.]
This is a callback, only sent when curdir is sent to the server. Format: 

###download [bytestream]
This is a callback, only sent when download is sent the server.

#Error Codes

##1xx

Errors that are associated with improper syntax.

###100

Sent back to client when the function is not recognized.

###101

Sent back to client if there are too many, or not enough arguments.

##2xx

Errors associated with files and file transfers.

###200

Sent back to client if the directory cannot be found.

###201

Sent back to client if the file cannot be found.

###202

Sent back to client if the directory requested is a file.

###203

Sent back to client if the directory already exists.

##3xx

Errors that are associated with Cliote creation.

###300

Sent back to client if the category is not recognized.

###301

Sent back to client if the Cliote name is already used.

##4xx

Errors associated with the login process.

###400

Sent back to client if the Cliote is already logged in.

###401

Sent back to client when the "password" is incorrect.

##9xx

Other errors.

###900

Sent back to client when the client tries to execute a query before authenticating.

###901

Sent back to client if a general error has occured.
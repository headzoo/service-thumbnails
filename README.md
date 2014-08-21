thumbnailer
===========
Used to create thumbnails from videos. Both normal thumbnails, and spites. Two apps are provided: main, which is a command line application, and server, which is an HTTP server.


### Requirements
* FFmpeg
* ImageMagick
* libmagic-dev


### Installation
`go install github.com/dulo-tech/thumbnailer`


### Command Line Usage
`main -t sprite -i video.mp4 -o thumb.jpg`


### Server Usage
First start the server using `server -h 127.0.0.1 -p 8888`. Then upload video files to the server. For example using curl `curl --form video=@video.mp4 -o thumb.jpg http://127.0.0.1:8888`.

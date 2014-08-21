thumbnailer
===========
Used to create thumbnails from videos. Both normal thumbnails, and spites. Two apps are provided: main, which is a command line application, and server, which is an HTTP server.


### Requirements
* FFmpeg
* ImageMagick
* libmagic-dev


### Installation
First make sure the requirements are installed, and then install the thumbnailer using:  
`go install github.com/dulo-tech/thumbnailer`


### Command Line Usage
Generating a sprite:  
`main -t sprite -i video.mp4 -o thumb.jpg`

Generating a big thumbnail:  
`main -t big -i video.mp4 -o thumb.jpg`

Generating thumbnails from several videos at once:  
`main -t big -i video1.mp4,video2.mp4,video3.mp4 -o thumb%02.jpg`


### Server Usage
First start the server using `server -h 127.0.0.1 -p 8888`. Then upload video files to the server. For example using curl `curl --form video=@video.mp4 -o thumb.jpg http://127.0.0.1:8888`. The server returns the thumbnail, which curl writes to thumb.jpg.

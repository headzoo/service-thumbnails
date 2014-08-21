thumbnailer
===========
Used to create thumbnails from videos. Both normal thumbnails, and sprites. Two apps are provided: thumbnailer, which is a command line application, and thumbnailer-server, which is an HTTP server.


### Requirements
* FFmpeg
* ImageMagick
* libmagic-dev


### Installation
First make sure the requirements are installed, and then install the thumbnailer using:  
`go install github.com/dulo-tech/thumbnailer`


### Command Line Usage
Generating a sprite:  
`thumbnailer -t sprite -i video.mp4 -o thumb.jpg`

Generating a big thumbnail:  
`thumbnailer -t big -i video.mp4 -o thumb.jpg`

Generating thumbnails from several videos at once:  
`thumbnailer -t big -i video1.mp4,video2.mp4,video3.mp4 -o thumb%02.jpg`


### Server Usage
First start the server using:  
`thumbnailer-server -h 127.0.0.1 -p 8888`

Then upload video files to the server. For example using curl:  
`curl --form video=@video.mp4 -o thumb.jpg http://127.0.0.1:8888/thumbnail/big`

The server returns the thumbnail, which curl writes to thumb.jpg. The server also has a help page which can be viewed at `http://127.0.0.1:8888/help`, and it implements the [Pulse Protocol](https://github.com/dulo-tech/amsterdam/wiki/Specification:-Pulse-Protocol).

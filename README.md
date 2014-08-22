service-thumbnails
==================
Used to create thumbnails from videos. Both normal thumbnails, and sprites. The app can run in one of two modes: cli and http server. Service Thumbnails is used to generate thumbnails from the command line when in cli mode, which is the default. In http mode service-thumbnails runs as an HTTP server capable of handling requests to generate thumbnails.


### Requirements
* FFmpeg
* ImageMagick
* libmagic-dev


### Installation
First make sure the requirements are installed, and then install the service-thumbnails using:  
`go install github.com/dulo-tech/service-thumbnails`

Linux binaries may be downloaded from the [releases page](https://github.com/dulo-tech/service-thumbnails/releases).


### Thumbnail Types
Thumbnailer currently creates two kinds of thumbnails: simple and sprite.

##### Simple
A simple thumbnail is a single frame from the video. By default the size (width/height) of the thumbnail is the size of the video frame. A video with frames 640x480 will result in a thumbnail that is 640x480. The size can be adjusted by using the 'width' option from the command line app or HTTP server.

Example:  
![Example Simple](http://i.imgur.com/HZUEppZ.jpg)


##### Sprite
A sprite thumbnail is two or more thumbnails from two or more video frames that have been stitched together into a single image. By default each thumbnail will be 180px wide, but can be changed using the 'width' option from the command line app or HTTP server. By default the sprite will always include 30 frames from the video, which have been chosen evenly from the full length of the video. That can be changed using the 'count' option from the command line or HTTP server.

Example:  
![Example Sprite](http://i.imgur.com/xSRxNbs.jpg)


### CLI Usage
Generating a sprite:  
`service-thumbnails -t sprite -i video.mp4 -o thumb.jpg`

Generating a simple thumbnail:  
`service-thumbnails -t simple -i video.mp4 -o thumb.jpg`

Generating thumbnails from several videos at once:  
`service-thumbnails -t simple -i video1.mp4,video2.mp4,video3.mp4 -o thumb%02.jpg`


### HTTP Usage
Start thumbnailer using the `-m http` switch:  
`service-thumbnails -m http -h 127.0.0.1 -p 8888`

Then upload video files to the server. For example using curl:  
`curl --form video=@video.mp4 -o thumb.jpg http://127.0.0.1:8888/thumbnail/simple`

The server returns the thumbnail, which curl writes to thumb.jpg. The server also has a help page which can be viewed at `http://127.0.0.1:8888/help`, and it implements the [Pulse Protocol](https://github.com/dulo-tech/amsterdam/wiki/Specification:-Pulse-Protocol).


### Configuration File
The command line options can also be specified in a configuration file. Pass the path to the file using the -conf switch. When the -conf switch isn't used the app will try to read from `$HOME/.service-thumbnails.conf`. Finally the app will try to read from `/etc/service-thumbnails.conf`.

See the [example configuration file](https://github.com/dulo-tech/service-thumbnails/blob/master/thumbnails.conf) for format and possible values.


### TODO
* The server needs to validate upload mime types.

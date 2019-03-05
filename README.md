# Writing Software for the Kindle Keyboard

It might be 9 years old, but the Kindle Keyboard is a pretty neat device. It runs Linux,
has a 600x800 eink Pearl screen that can display 16 shades of grey, and has a small keyboard,
five-point control, page buttons, and a speaker.

It also turns out that you can compile Go code for the device using the flags `env GOOS=linux GOARCH=arm GOARM=6 go build program.go`.

Background info is in paragraphs; steps I take will be in unordered lists.

## Day 1

The first day. I spent a bunch of time doing research and seeing what options there for writing software for the Kindle 3.

There was a Java SDK called the Kindle Development Kit or KDK) released by Amazon not long after the device was released, but it was EOL'ed in 2013 and is not easily available. The device runs Java 1.4, which is pretty old.

The central place for finding out information about Kindle development is the [Kindle Developer's Corner](https://www.mobileread.com/forums/forumdisplay.php?f=150) on the MobileRead forums. Interestingly, much of the code that folks have written hasn't made it to Github or other code sharing sites.

Most third-party code is distributed in the same format as official Kindle firmware updates by placing it in the root of the Kindle USB device and restarting the Kindle. The format of these is some kind of obfuscated tar file.

Someone made a [neat weather display](https://mpetroff.net/2012/09/kindle-weather-display/) that is nice inspiration.

I have seen many links to [this site](http://cowlark.com/kindle/getting-started.html) that describes extracting the Jars from the Kindle and writing a Java app using the KDK.

## Day 2

I need to be able to get a shell and poke around. I'm a little worried as this thing is almost a decade old, and I'm not sure what will be run on it. Evidently it's libc is from 2006.

* Applied the jailbreak and USBNetworking hacks from [Font, ScreenSaver & USBNetwork Hacks for Kindle 2.x, 3.x & 4.x - MobileRead Forums](https://www.mobileread.com/forums/showthread.php?t=88004).
* Turned on USBNetworking on the Kindle by entering the following codes on the homescreen in the box that appears when you press `del`: `;debugOn`, `~useNetwork`, `;debugOff`. There are [all kinds of commands that can be entered here](https://ebooks.stackexchange.com/questions/152/what-commands-can-be-given-in-the-kindles-search-box), depending on the Kindle version and if it has been jailbroken. You'll know that USBNetworking is enabled because the connection will appear as green in the Network panel in the Mac's System Preferences.
* I had to [use this configuration](https://www.mobileread.com/forums/showpost.php?p=2895606&postcount=13) on my Mac to setup the USB network interface.

There is a command called `eips` that will perform various functions with the screen, but all I've been able to do is get it to clear the screen using `eips -c` and print info about the screen using `eips -i`.

* I tried to display a PNG, and something happened, but the image only partially was displayed and was stretched. [This image from the weather project](https://github.com/mpetroff/kindle-weather-display/blob/master/kindle/weather-image-error.png), however, worked, so it's possible. There is probably something wrong with the image format I am using.

## Day 3

Today I wanted to try to draw something to the screen.

* Disabled built-in kindle UI with:`/etc/init.d/framework stop`. I also looked around at the other stuff in `/etc/init.d`.

I looked back at the weather display project and noticed that the PNGs need to be created with a `color_type` value in the PNG header to 0. `pngcrush -c 0` will do this.

* I was able to display arbitrary graphics on the device using `eips -g`!

I had previous cross-compiled Go and run it on other ARM devices, so I figured it was worth a shot.

* I compiled the example program from [this article](https://dave.cheney.net/2015/03/03/cross-compilation-just-got-a-whole-lot-better-in-go-1-5):

  ```golang
  package main

  import "fmt"
  import "runtime"

  func main() {
          fmt.Printf("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)
  }
  ```

  `env GOOS=linux GOARCH=arm go build test.go` produced an executable that resulted in an `Illegal instruction` when run on the Kindle. Adding `GOARM=6` fixed this, though, and allowed the program to run!
* So far, compiling on the device, `scp`ing the program over to the Kindle, and running it in SSH session has been working pretty well. At some point it would be nice to automate this.

## Day 4

Today I wanted to draw some graphics. The Kindle Keyboard has a framebuffer device available at `/dev/fb0`, so that is where I'm going to start. Once I saw that Go would run on the device, I became a lot less interested in using Java and the KDK.

A 4-bit number will hold values from 0-15, which is all that's needed to represent one of the 16 levels of grey supported by the display.

The framebuffer on the Kindle packs two 4-bit pixels into each byte. This means that each row of the display uses 300 bytes to hold the values of its 600 pixels. This is enough information to start working with the raw framebuffer data. Writing to a pixel will involve setting either the higher or lower 4 bits of the appropriate byte to a value representing the grey value to display.

* I tried using some code I had written for another experiment to draw to the framebuffer, leveraging [this framebuffer library](https://github.com/kaey/framebuffer). Unfortunately, I ran into an issue because the value of `Smem_len` is zero, so this library thinks there is nowhere to write data to. I hardcoded `483328` (which I got from running `eips -i`), which allowed something to be drawn to the screen. Looking at what I would need to do to modify that library to handle 4-bit pixels packed two to a byte instead of 4-byte RGBA pixels, I decided it would be easy enough to `Mmap` the framebuffer myself and write a few utility functions. I'm pretty sure that I only need to map 240,000 bytes to do what I need to do, and since this code is only ever going to run on this device (any maybe a DX if I find one), hardcoding these values is easier than troubleshooting why the `FBIOGET_FSCREENINFO` ioctl syscall isn't returning the expected value for `Smem_len`.
* After writing to the framebuffer, you have to tell the device to update the display. So far I have been using `echo 1 > /proc/eink_fb/update-display` to do this. There is probably a better way.
* I was able to draw some simple graphics on the display at the end of this session.

## Day 5

Today's another day for research and poking around.

`/etc/init.d/battcheck` returns a bunch of interesting battery statistics:

```
system: I battcheck:def:running
Sat Jan 19 19:52:48 2019  INFO:battery voltage: 4136 mV
Sat Jan 19 19:52:48 2019  INFO:battery charge: 86%
system: I battcheck:def:current voltage = 4133mV
Sat Jan 19 19:52:48 2019  INFO:battery charge: 86%
Sat Jan 19 19:52:48 2019  INFO:battery voltage: 4136 mV
Sat Jan 19 19:52:48 2019  INFO:battery current: 337 mA
system: I battcheck:def:gasgauge capacity=86% volts=4136 mV current=337 mA
system: I battcheck:def:Waiting for 3460mV or 4%
system: I battcheck:def:battery sufficient, booting to normal runlevel
```

Here's the output of `lsmod`:

```
Module                  Size  Used by
ar6000                161076  0
g_ether                21096  0
eink_fb_shim          116732  0
eink_fb_hal_broads    397532  0
eink_fb_hal            59764  5 eink_fb_shim,eink_fb_hal_broads
volume                  8900  0
fiveway                23552  0
mxc_keyb               15904  0
uinput                  7776  0
fuse                   48348  2
arcotg_udc             38628  1 g_ether
mwan                    7324  0
```

There isn't a `tree` command on the Kindle, and I've been using `ls -R` a lot to explore the filesystem. I'm considering `scp`ing the entire disk to the Mac so I can use my usual tools on it.

The Kindle is running [alsa](https://www.alsa-project.org/main/index.php/Main_Page) version 1.0.13:

```
$ alsactl -v
alsactl version 1.0.13`
```

To run `alsamixer`, `TERM` needs to be set to `xterm` (mine was `xterm-256color`).

## Day 6

I spent some timing looking at graphics packages, with an eye toward rendering text to the screen in a scalable way. I originally planned to use bitmap fonts due to their ease of use, but then I found a [pure-Go implementation of freetype](https://github.com/golang/freetype), and then, a while later, [gg](https://github.com/fogleman/gg), which provides a nice API for drawing graphics and text. It even includes wrapping text to lines, which is something that isn't provided by `freetype`.

`man` isn't available on the Kindle, which makes figuring out the arguments for the old versions of everything a little more challenging.

One of the twists of working with the eink display is that values that would appear on light-emitting displays as dark colors instead appear as light colors on eink displays. This in effect inverses everything, which needs to be compensated for somewhere in the graphics stack of a program.

* Wrote a simple program `circle` to draw a black circle on the center of the Kindle screen, centering the string `Hello!` within it. There are 7 concentric circles with decreasing stroke width surrounding it.

* Finally setup ssh keys so I could ssh into the Kindle without hitting enter at the password prompt by following [these directions](https://www.mobileread.com/forums/showthread.php?t=204942). I copied the public key over using  `scp kindle_rsa.pub root@192.168.2.2:/mnt/us/usbnet/etc/authorized_key`.

* Added `build_and_run`, a script to automate the process of compiling a program, copying it to the Kindle, and running it.

In the back of my mind I was a little worried about memory on the Kindle, but it's looking like things are going to be just fine:

```
             total       used       free     shared    buffers     cached
Mem:           250        219         30          0         82         97
-/+ buffers/cache:         40        210
Swap:            0          0          0
```

## Day 7

* Added a new executable, `screengrab`, that generates a PNG from the current state of the framebuffer. Wrapped this in a script to capture a screenshot on the device, `scp` it back to the host, and open it in the default program.
* Moved all scripts into `scripts` directory to keep the root of the repo tidy.

## Day 8

* I attempted to play a variety of sound files, including those I found on the Kindle via `aplay`, but all that came out of the speaker was a deafening static that came and went in waves. Not completely unlike an ocean, really, but also fairly unpleasant and not at all resembling the [piano music](http://www.pianosociety.com/pages/fieldnocturnes/) I was hoping to hear. [This page](), however, prompted me to attempt to play a file that was created on the device using `arecord`, which worked! So converting audio to the format of that file should allow for custom sound to be played. 8khz mono isn't going to sound great, but it's better than nothing.

* I later learned that the issue is that `aplay` expects raw audio data (say, from a WAV file). A 44.1khz stereo WAV played fine. I was previously trying to play an MP3, which isn't raw audio data and needs to be decoded first.

`gasgauge` returns stats related to the battery and charging status of the device:

```
[root@kindle root]# gasgauge-info -c
100%
Tue Jan 22 03:01:11 2019  INFO:battery charge: 100%
```

The Kindle has a `say` command that will speak arbitrary text ([source](https://grenville.wordpress.com/2011/09/26/kindle-3-playing-with-text-to-speech/)). I confirmed that this works!

`evtest` is available on the Kindle, which makes exploring the hardware keyboard pretty straight forward. The main keyboard, five-way control, and paging buttons are all seperate devices.

* Wrote a small program to read key events from the main keyboard (`/dev/input/event0`) and print the raw bytes to the screen:
  ```
  $ script/build_and_run keys
  ...16 [71 174 70 92 195 96 12 0 1 0 52 0 1 0 0 0]
  16 [72 174 70 92 232 199 0 0 1 0 52 0 0 0 0 0]
  16 [72 174 70 92 31 218 10 0 1 0 52 0 1 0 0 0]
  16 [72 174 70 92 232 252 12 0 1 0 52 0 0 0 0 0]
  16 [72 174 70 92 175 170 14 0 1 0 52 0 1 0 0 0]
  16 [73 174 70 92 24 61 1 0 1 0 52 0 0 0 0 0]
  ```

## Day 9

The Kindle Keyboard Linux kernel is `2.6.26`:

```
[root@kindle root]# uname -r
2.6.26-rt-lab126
```

[This site](https://elixir.bootlin.com/linux/v2.6.26/source/include/linux/time.h) is super useful for looking up the definitions of `input_type` on a specific version of the Linux kernel, which I need to do in order to load the raw data I'm receiving from `/dev/input/event0` into a Go struct.

* Wrote code to handle processing the events coming from `/dev/input/event0` (the main keyboard) and push Go structs representing them onto a channel. Used [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer) to generate the `String()` method for these types.

The Kindle Keyboard hardware or drivers (not sure which) are interesting. Look at what events are sent when the the _shift_ key is pressed and released, followed by the _z_ key:

```
{Time:2019-01-23 03:30:09.780995 +0015 GMT-00:20 Type:KeyDown Key:KeyShift}
{Time:2019-01-23 03:30:09.871003 +0015 GMT-00:20 Type:KeyUp Key:KeyShift}
{Time:2019-01-23 03:30:10.70096 +0015 GMT-00:20 Type:KeyDown Key:KeyZ}
{Time:2019-01-23 03:30:10.860955 +0015 GMT-00:20 Type:KeyUp Key:KeyZ}
```

Compare this to the same thing, but for the _alt_ key:

```
{Time:2019-01-23 03:30:15.191005 +0015 GMT-00:20 Type:KeyDown Key:KeyAlt}
{Time:2019-01-23 03:30:15.191294 +0015 GMT-00:20 Type:KeyDown Key:KeyZ}
{Time:2019-01-23 03:30:15.191521 +0015 GMT-00:20 Type:KeyUp Key:KeyZ}
{Time:2019-01-23 03:30:15.191528 +0015 GMT-00:20 Type:KeyUp Key:KeyAlt}
```

Notice how the _alt_ `KeyUp` event isn't sent until after another key is pressed, which is the same thing you see if the modifier was held down by the user:

```
{Time:2019-01-23 03:30:17.881007 +0015 GMT-00:20 Type:KeyDown Key:KeyShift}
{Time:2019-01-23 03:30:18.180979 +0015 GMT-00:20 Type:KeyDown Key:KeyZ}
{Time:2019-01-23 03:30:18.400987 +0015 GMT-00:20 Type:KeyUp Key:KeyZ}
{Time:2019-01-23 03:30:18.630976 +0015 GMT-00:20 Type:KeyUp Key:KeyShift}

{Time:2019-01-23 03:30:22.240992 +0015 GMT-00:20 Type:KeyDown Key:KeyAlt}
{Time:2019-01-23 03:30:22.241281 +0015 GMT-00:20 Type:KeyDown Key:KeyZ}
{Time:2019-01-23 03:30:22.400988 +0015 GMT-00:20 Type:KeyUp Key:KeyZ}
{Time:2019-01-23 03:30:22.700979 +0015 GMT-00:20 Type:KeyUp Key:KeyAlt}
```

This also means that (at least using the technique I am and listening to `dev/input/event0`) it's impossible to detect a keypress of only the `alt` key.

* Wrote a utility program, `simulate_eink`, to convert a PNG by mapping the shades of gray to a pallete that is perceptually much closer to what the eink screen looks like to a human. It also adds a bit of random noise for realism.

* Updated the `draw` command to use the latest `FrameBuffer` code.

* Added example images for `draw` and `circle`.

## Day 10

Today I wanted to run the built-in `say` whenever a key on the keyboard was pressed. Unfortunately, I hit an issue with using `exec.Command`:

```
$ script/build_and_run letters
...Q
goroutine 1 [running]:
runtime/debug.Stack(0x1045a000, 0xe8fa0, 0x104481e0)
	/usr/local/Cellar/go/1.10.2/libexec/src/runtime/debug/stack.go:24 +0x80
main.main()
	/Users/jimb/go/src/github.com/jim/kindleland/cmd/letters/letters.go:25 +0x124

panic: fork/exec /usr/bin/say: function not implemented

goroutine 1 [running]:
main.main()
	/Users/jimb/go/src/github.com/jim/kindleland/cmd/letters/letters.go:26 +0x230
```

After some searching online, I was worried that `exec.Command` might be relying on glibc, which the Kindle has an ancient version of:

```
[root@kindle root]# /lib/libc.so.6
GNU C Library stable release version 2.5, by Roland McGrath et al.
Copyright (C) 2006 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.
There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.
Compiled by GNU CC version 4.1.2.
Compiled on a Linux 2.6.15 system on 2008-06-10.
...
```

That software was compiled over a decade ago. And indeed there are threads about Golang not running on old versions of glibc, and I couldn't find a specific glibc version requirement, but I wanted to get a better look at what was going on. To examine what syscalls the program was making and what/which errors it was getting back, I ran the program with `strace` and saw the following:

```
pipe2(0x10431dac, O_CLOEXEC)            = -1 ENOSYS (Function not implemented)
pipe2(0x10431dac, 0)                    = -1 ENOSYS (Function not implemented)
```

The [pipe2 syscall](http://man7.org/linux/man-pages/man2/pipe.2.html) was added in Linux 2.6.27, and the Kindle has 2.6.23 (which is currently listed as Go's minimum supported Linux version).

`os.Pipe` is called several layers within the `exec.Command` code when `Command.Run()` is executed. There is a [fallback in the 
code](https://github.com/golang/go/blob/master/src/os/pipe_linux.go#L17-L27) to handle `pipe2` not being supported, and you can see how the calls to `syscall.Pipe2` and `syscall.Pipe` map to the two syscalls shown above.

```golang
// Pipe returns a connected pair of Files; reads from r return bytes written to w.
// It returns the files and an error, if any.
func Pipe() (r *File, w *File, err error) {
	var p [2]int

	e := syscall.Pipe2(p[0:], syscall.O_CLOEXEC)
	// pipe2 was added in 2.6.27 and our minimum requirement is 2.6.23, so it
	// might not be implemented.
	if e == syscall.ENOSYS {
		// See ../syscall/exec.go for description of lock.
		syscall.ForkLock.RLock()
    e = syscall.Pipe(p[0:])
    ...
```

Except, of course, that the second syscall should be `pipe`, not `pipe2` based on the fallback code. The problem is that `syscall.Pipe` is implemented using the `pipe2` syscall on linux/arm! This change happened [here](
https://github.com/golang/go/commit/9b3ccc082f6bda01727fc98096f7d197bba830db#diff-4457f0b3f1fea7981d43c6657258588bR31).

By restoring the previous definition of `syscall.Pipe` in `syscall/syscall_linux_arm.go`, I was able to get my code that uses `exec.Command` to run properly on the Kindle.

In the process of working through this, I attemped to compile Delve for some on-device debugging before discovering that [Delve doesn't support ARM](https://github.com/go-delve/delve/issues/118).

And of course, after digging into these syscalls, I read that you can also [just write text to `/var/tmp/ttsUSFifo`](https://grenville.wordpress.com/2011/09/26/kindle-3-playing-with-text-to-speech/) to use the Kindle's text to speech.

* Wrote a little program, `speak`, that uses `say` to speak the name of each key as it is pressed.

## Day 11

I'm having an issue using `gg` to render text to an image and display it on the screen. There is a memory leak somewhere, and after a certain number of updates the devices bogs down and things get weird. I am going to do some troubleshooting and see if the issue is in my code and in the way I am reusing the `gg.Context` multiple times instead of creating a new one each time I want to draw something.

[FBink](https://github.com/NiLuJe/FBInk) is a C library that does a lot of what I want to do. It was originally designed for the Kobo but now also supports Kindles and other Eink devices (they tend to have very similar hardware and software stacks). There are [Go bindings for it](https://github.com/shermp/go-fbink-v2), but to use them you have to enable `cgo`, which is something I would like to avoid to keep the build process as simple as possible. However, these libraries are excellent references as I move forward with a pure-Go approach.

Specifically interesting are the parts of FBInk that expose how to update only part of the screen. So far I have been doing entire device updates, which are slow and provide a jarring experience for the user.

Amazon posts the [source code they are required to release](https://www.amazon.com/gp/help/customer/display.html?nodeId=200203720) for all Kindle devices and apps. `linux-2.6.26/include/linux/einkfb.h` shows a lot of the details used by FBInk to do its work (and is actually used directly in that project).

## Day 12

Today was a day spent learning about `cgo`, linux headers, and `go generate`.

`go tool cgo -godefs` ignores the special `// #cgo` comments, so options that need to be passed to the C compiler have to be passed on the command line. The docs aren't super clear about this, but I was able to sort it out by using the `-gcc-debug` flag to `cgo` and then running that output through `clanng` myself, adding the `-v` flag so I could experiment with which options needed to be passed to get the lookup paths correct.

I ended up using the following to generate Go code from the `einkfb.h` file included in the GPL source distribution for the Kindle:

```
go tool cgo -godefs -- -Ivendor/linux-2.6.26-lab126/include -D__KERNEL__ constant_defs.go > constants.go
```

By putting this line in a shell script, `script/generate_constants`, I was able to invoke it by adding a special comment to `constant_defs.go` and then running `go generate`.

## Day 13

I went down a rabbit hole learning more about linux syscalls, ioctls, and how to interact with them in Go. I got the basic screen update working by making an IOCTL call on `/dev/fb/0` and was able to specify "fast" or "slow", although I haven't yet sorted out what the difference between the different options are.

My next task is to sort out how to do a partial screen update. To do so, I need to pass a pointer to a struct into the syscall which contains information about how to do the update: what areas to (not) update, what _FX_ to use, etc. There are a few levels of software running on the Kindle 3, which makes keeping everything straight a little bit harder.

[This article](https://blog.gopheracademy.com/advent-2017/unsafe-pointer-and-system-calls/) was the best thing I found while trying to figure out how I would pass an area of the screen to update.

It appears that `cgo -godefs` doesn't support most macros, which makes the way I was trying to define the constants from a pervious day a dead end as the header I am working with includes a lot of stuff like this:

```c
#define FBIO_EINK_UPDATE_DISPLAY            _IO(FBIO_MAGIC_NUMBER, 0xdb) // 0x46db (fx_type)
#define FBIO_EINK_UPDATE_DISPLAY_AREA       _IO(FBIO_MAGIC_NUMBER, 0xdd) // 0x46dd (update_area_t *)
```

I am probably just going to define the values I need in Go as they come up instead of attempting to autogenerate things from the C header. `godefs` may prove to be a useful tool, though, because it will automate the conversion when I need new values.

I also saw that the built-in `syscall` package is considered deprecated and that you are supposed to use [sys](https://godoc.org/golang.org/x/sys/unix) instead.

## Day 14

Today I decided to try to use the FBInk library a try. I did end up getting it to compile in a Docker container, although it took a lot of time to get everything working and for the cross-compiling toolchain it expects to install itself. I'm still pretty sold on avoiding C, but I am close to being able to use this library to see if that is useful.

```docker
FROM ubuntu:latest
RUN apt-get update
RUN apt-get -y install gperf help2man bison texinfo flex gawk git build-essential autoconf libncurses5-dev curl wget file
WORKDIR /root/src
RUN git clone https://github.com/koreader/koxtoolchain.git
WORKDIR /root/src/koxtoolchain
ENV CT_EXPERIMENTAL=y
ENV CT_ALLOW_BUILD_AS_ROOT=y
ENV CT_ALLOW_BUILD_AS_ROOT_SURE=y
RUN ./gen-tc.sh kindle
WORKDIR /root/src
RUN git clone https://github.com/NiLuJe/FBInk
WORKDIR /root/src/FBInk
RUN git submodule update --init
```

## Day 15

Spent some time getting the various screen update functions to work from Go without using FBInk. Ran `strace eips -c` to see how that program cleared the screen. The interesting part:

```
...
open("/dev/fb/0", O_RDWR)               = 3
ioctl(3, FBIOGET_VSCREENINFO or PF_IOCTL_INIT, 0xbeb52ae0) = 0
mmap2(NULL, 240000, PROT_READ|PROT_WRITE, MAP_SHARED|MAP_LOCKED, 3, 0) = 0x40146000
msync(0x40146000, 240000, MS_SYNC)      = 1
ioctl(3, FBIO_EINK_CLEAR_SCREEN, 0)     = 0
close(3)
...
```

It's good to see that program is also `mmap2`ing `240000` bytes, just as I am.

It seems that you need to clear the screen twice to really get a clean slate.

I also discovered that the main keyboard send different keycodes when `Alt` is combined with the top row of letter keys:

```
# Alt-Q
{Time:2019-02-06 02:13:56.412672 +0015 GMT-00:20 Type:KeyDown Key:KeyType(2)}
{Time:2019-02-06 02:13:56.57269 +0015 GMT-00:20 Type:KeyUp Key:KeyType(2)}

# Alt-P
{Time:2019-02-06 02:14:01.742672 +0015 GMT-00:20 Type:KeyDown Key:KeyType(11)}
{Time:2019-02-06 02:14:01.892816 +0015 GMT-00:20 Type:KeyUp Key:KeyType(11)}
```

## Day 16

* Wrote a small program that used the  `freetype` package to draw text to the screen.

## Day 17

* Improved `text` program to have it wrap text at the end of a line and draw within a defined part of the screen.

## Day 18

* Added a new `letters` program that allows large letters to be typed across the screen. It becomes sluggish when many keys are pressed quickly. I will need to add some throttling to the screen updating and move updating the buffer and telling the screen to refresh to a goroutine.
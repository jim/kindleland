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

# Day 3

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

# Day 4

Today I wanted to draw some graphics. The Kindle Keyboard has a framebuffer device available at `/dev/fb0`, so that is where I'm going to start. Once I saw that Go would run on the device, I became a lot less interested in using Java and the KDK.

A 4-bit number will hold values from 0-15, which is all that's needed to represent one of the 16 levels of grey supported by the display.

The framebuffer on the Kindle packs two 4-bit pixels into each byte. This means that each row of the display uses 300 bytes to hold the values of its 600 pixels. This is enough information to start working with the raw framebuffer data. Writing to a pixel will involve setting either the higher or lower 4 bits of the appropriate byte to a value representing the grey value to display.

* I tried using some code I had written for another experiment to draw to the framebuffer, leveraging [this framebuffer library](https://github.com/kaey/framebuffer). Unfortunately, I ran into an issue because the value of `Smem_len` is zero, so this library thinks there is nowhere to write data to. I hardcoded `483328` (which I got from running `eips -i`), which allowed something to be drawn to the screen. Looking at what I would need to do to modify that library to handle 4-bit pixels packed two to a byte instead of 4-byte RGBA pixels, I decided it would be easy enough to `Mmap` the framebuffer myself and write a few utility functions. I'm pretty sure that I only need to map 240,000 bytes to do what I need to do, and since this code is only ever going to run on this device (any maybe a DX if I find one), hardcoding these values is easier than troubleshooting why the `FBIOGET_FSCREENINFO` ioctl syscall isn't returning the expected value for `Smem_len`.
* After writing to the framebuffer, you have to tell the device to update the display. So far I have been using `echo 1 > /proc/eink/update-display` to do this. There is probably a better way.
* I was able to draw some simple graphics on the display at the end of this session.

# screenmgr

**screenmgr** is a tool for managing and monitoring a pool of networked devices.

With it, you can

* monitor uptime/reachability via ping (ICMP)
* boot (WOL), shutdown, restart, or commit other device actions (SSH)
* view the screen (VNC or platform-specific snapshot mechanisms)
* view operating system and hardware info (SSH)

![screenmgr](http://i.imgur.com/YcH2aFO.jpg)

I made this mostly so I could monitor my collection of Power Macs. It originally
only supported Mac OS 8/9/X, but it now supports other UNIX-like operating
systems, as well as Windows to some extent.

It checks for a variety of protocols on each device and tries to
find as much info as it can to show on the web interface. The availability of
each feature depends on the platform.

## License

[ISC](LICENSE)

## Author

[Mitchell Cooper](https://mitchellcooper.me), <mitchell@mitchellcooper.me>

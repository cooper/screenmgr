# screenmgr

**screenmgr** is a tool for managing and monitoring a pool of networked devices.

with it, you can

* monitor uptime/reachability via ping (ICMP)
* boot (WOL), shutdown, restart, or commit other device actions (SSH)
* view the screen (VNC or platform-specific snapshot mechanisms)
* view operating system and hardware info (SSH)

![screenmgr](http://i.imgur.com/YcH2aFO.jpg)

I made this mostly so I could monitor my collection of Power Macs. it originally
only supported Mac OS 8/9/X, but it now supports other unix-like operating
systems as well as windows to some extent.

basically it checks for a variety of protocols on each device and tries to
find as much info as it can to show on the web interface. the availability of
each feature depends on the platform.

more features/OS support planned.

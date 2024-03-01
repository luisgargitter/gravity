# Installation
## Dependencies
See [go.mod](go.mod).
The rendering is done with double-precision, therefore OpenGl version 4.1 is needed.

To install the needed glfw headers see [glfw-compilation](https://www.glfw.org/docs/3.3/compile.html).
Check supported OpenGL-version with `glxinfo | grep "OpenGL version"`.

# Running the program
Use the command `go run .` to compile and execute the program (Note: first time compilation takes about 3-5 minutes, as the used go-libraries need to be compiled.).
Program initialization may take a few moments, as the textures used are all about 8k-textures, which needs time to be generated and loaded onto the gpu. The operationg system might give a `Program is not responding` warning. Texture loading takes about 10-20 seconds (on my machine).

# Controls
By default the movement is done like in typical FPS-Games (`w`-`a`-`s`-`d`-`shift`(down)-`space`(up)), but can be ajusted using the functions `FreeMove` and `FreeLook` directly.
There exists the additional feature to have a geocentric view by pressing `tab`(hold). This locks the camera with the earth in centered on the screen and all movement relative to earth.
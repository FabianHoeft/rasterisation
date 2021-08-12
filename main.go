package main


import (
	"fmt"
	"time"
	"image/color"
	"math"

)




func main() {
	t_start:=time.Now()
	xres := 190
	yres := 100
	cores := 20
	drawdist := 20.
	blocksize :=xres*yres/cores/8
	Fov := 90.

	T0:=Newtriangle(vec3{9.,0.,0.},vec3{3.,2.,0.},vec3{3.,0.,1.},color.RGBA{20,150,10,0xff})
	P0:=Newpolygon([]vec3{vec3{4.,0.5,0.},vec3{4.,-2.,0.},vec3{4.,0.5,1.},vec3{4.,-2.,1.}},[][2]int{{1,0},{1,3},{3,2},{2,0}},color.RGBA{120,20,200,0xff} )
	T1:=Newtriangle(vec3{4.,0.5,0.},vec3{4.,0.5,1.},vec3{4.,-2.,0.},color.RGBA{120,20,200,0xff})
	T2:=Newtriangle(vec3{4.,-2.,1.},vec3{4.,-2.,0.},vec3{4.,0.5,1.},color.RGBA{120,20,200,0xff})
	T4:=Newtriangle(vec3{8.,0.,0.},vec3{3.,-2.,0.},vec3{3.,0.,-1.},color.RGBA{50,20,50,0xff})
	Pl0:=Newplayer(0,math.Pi/2.,vec3{0,0,0},Fov ,xres,yres ,drawdist, cores , blocksize)

	t_init:=time.Now()
	fmt.Println("init:", t_init.Sub(t_start))

	C0:=multirender(Pl0, []geometry{T0,T1,T2,T4,P0} )

	t_render:=time.Now()
	fmt.Println("render:", t_render.Sub(t_init))

	C0.print("view.png")
	C0.printzmap("zmap.png")

	t_done:=time.Now()
	fmt.Println("store:", t_done.Sub(t_render))


}

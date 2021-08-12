package main

import(
	"math"
	"image/color"
	"log"
	"os"
	"image"
)

type backtexture struct{
	texture [][]color.RGBA
	vres,hres int
}

func (b backtexture) color(vel vec3) color.RGBA{
	x:=int((math.Atan(vel.y/vel.x)/math.Pi+0.5)*float64(b.vres))
	y:=int(math.Acos(vel.z)/(math.Pi)*float64(b.hres))

	return b.texture[y][x]

}

type backsolid struct{
	colorset color.RGBA
}

func (b backsolid) color(vel vec3) color.RGBA{
	return b.colorset
}

func getbackground() [][]color.RGBA {
	reader, err := os.Open("skytexture0.PNG")
  if err != nil {
          log.Fatal(err)
  }
  defer reader.Close()

  m, _, err := image.Decode(reader)
  if err != nil {
          log.Fatal(err)
  }
  bounds := m.Bounds()
	out:= make([][]color.RGBA, bounds.Max.Y)
  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		out[y]=make([]color.RGBA, bounds.Max.X)
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r,g,b,a:=m.At(x, y).RGBA()
      out[y-bounds.Min.Y][x-bounds.Min.X] =color.RGBA{uint8(r),uint8(g),uint8(b),uint8(a)}
		}
 	}
	return out
}

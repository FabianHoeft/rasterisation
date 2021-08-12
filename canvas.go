package main


import (
  "image/color"
  "image"
  "os"
  "image/png"
  "math"
)


type canvas struct{
  pic [][]color.RGBA
  z [][]float64
  xres,yres int
  x0,y0 int
  res float64
}

func emptycanvas(P player) canvas {
  pic:=make([][]color.RGBA, P.yres,P.yres)
  zs:=make([][]float64, P.yres,P.yres)

  for i := 0; i < P.yres; i++ {
    fill0:=make([]color.RGBA, P.xres,P.xres)
    fill1:=make([]float64, P.xres,P.xres)
    for j:=0; j < P.xres; j++ {
      fill1[j]=P.drawdist
      fill0[j]=color.RGBA{0,0,0,255}
    }
    pic[i]=fill0
    zs[i]=fill1
  }
  return canvas{pic,zs,P.xres,P.yres,-P.xres/2,-P.yres/2,P.res}
}

func newcanvas(xres int,yres int,x0 int,y0 int,res float64,drawdist float64) canvas {
  pic:=make([][]color.RGBA, yres)
  zs:=make([][]float64, yres)
  for i := 0; i < yres; i++ {
    pici:=make([]color.RGBA, xres)
    zsi:=make([]float64, xres)
    for j := 0; j < xres; j++ {
      zsi[j]=drawdist
      pici[j]=color.RGBA{0,0,0,255}
    }
    zs[i]=zsi
    pic[i]=pici
  }
  return canvas{pic,zs,xres,yres,x0,y0,res}
}

func (C0 canvas) horizontaladd( C1 canvas) canvas {
  pic:=make([][]color.RGBA, C0.yres)
  zs:=make([][]float64, C0.yres )
  for i := 0; i < C0.yres; i++ {
    pic[i]=append(C0.pic[i],C0.pic[i]...)
    zs[i]=append(C0.z[i],C0.z[i]...)
  }
  return canvas{pic,zs,C0.xres+C1.xres,C0.yres,C0.x0,C0.y0,C0.res}
}

func (C0 canvas) verticaladd( C1 canvas) canvas {
  pic:=make([][]color.RGBA, C0.yres+C1.yres)
  zs:=make([][]float64, C0.yres+C1.yres )
  for i := 0; i < C0.yres; i++ {
    pic[i]=C0.pic[i]
    zs[i]=C0.z[i]
  }
  for i := 0; i < C1.yres; i++ {
    pic[i+C0.yres]=C1.pic[i]
    zs[i+C0.yres]=C1.z[i]
  }
  return canvas{pic,zs,C0.xres,C0.yres+C1.yres,C0.x0,C0.y0,C0.res}
}

func (C0 canvas) add( C1 canvas) canvas   {
  if C0.xres==C1.xres || C0.x0 == C1.x0 {
    if C0.y0+C0.yres==C1.y0 {
      return C0.verticaladd(C1)
    } else if C1.y0+C1.yres==C0.y0 {
      return C1.verticaladd(C0)
    }
  } else if C0.yres==C1.yres || C0.y0 == C1.y0 {
    if C0.x0+C0.xres==C1.x0 {
      return C0.horizontaladd(C1)
    } else if C1.x0+C1.xres==C0.x0 {
      return C1.horizontaladd(C0)
    }
  }

  x0:=C0.x0
  if C1.x0>x0 {
    x0=C1.x0
  }
  x1:=C0.x0+C0.xres
  if x1>C1.x0+C1.xres {
    x1=C1.x0+C1.xres
  }
  xres:=x1-x0

  y0:=C0.y0
  if C1.y0>y0 {
    y0=C1.y0
  }
  y1:=C0.y0+C0.yres
  if y1>C1.y0+C1.yres {
    y1=C1.y0+C1.yres
  }
  yres:=y1-y0

  pic:=make([][]color.RGBA, yres)
  zs:=make([][]float64, yres )
  for i := 0; i < yres; i++ {
    pic[i]=make([]color.RGBA, xres)
    zs[i]=make([]float64, xres)
    for j := 0; j < xres; j++ {
      ztemp0:=C0.z[i-y0+C0.y0][j-x0+C0.x0]
      ztemp1:=C1.z[i-y0+C1.y0][j-x0+C1.x0]
      if ztemp0 < ztemp1 {
        zs[i][j]=ztemp0
        pic[i][j]=C0.pic[i-y0+C0.y0][j-x0+C0.x0]
      } else {
        zs[i][j]=ztemp1
        pic[i][j]=C1.pic[i-y0+C1.y0][j-x0+C1.x0]
      }
    }
  }
  return  canvas{pic,zs,xres,yres,x0,y0,C0.res}
}




func (C canvas) print(loc string) {
  upLeft := image.Point{0, 0}
	lowRight := image.Point{C.xres, C.yres}
  img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
  for x := 0; x < C.xres; x++ {
      for y := 0; y < C.yres; y++ {
          img.Set(x,y,C.pic[y][x])
      }
  }

  f, _ := os.Create(loc)
  png.Encode(f, img)
}

func (C canvas) printzmap(loc string) {
  upLeft := image.Point{0, 0}
	lowRight := image.Point{C.xres, C.yres}
  img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
  min,max:= 0.,0.
  inf:=math.Inf(1)
  for x := 0; x < C.xres; x++ {
    for y := 0; y < C.yres; y++ {
      cur:=C.z[y][x]
      if cur!=inf{
        if cur>max {
          max=cur
        }
        if cur<min {
          min=cur
        }
      }
    }
  }
  for x := 0; x < C.xres; x++ {
    for y := 0; y < C.yres; y++ {
      c:=uint8(255.*(1.-(C.z[y][x]-min)/(max-min)))
      img.Set(x,y,color.RGBA{c,c,c,0xff})
    }
  }

  f, _ := os.Create(loc)
  png.Encode(f, img)
}

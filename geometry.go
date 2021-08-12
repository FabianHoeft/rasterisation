package main


import (
  "sort"
  "image/color"
  "fmt"

)

type geometry interface{
  rasterize( C canvas)
  transform(M Tmat)  geometry
  shouldrender() (float64,float64,float64 ,float64) // return (x^2+y^2)/z^2, z dist, direction approx and size

}


func render(P player,objects []geometry) canvas {
  C:=emptycanvas(P)
  m:=map[int]float64{}
  var objectstrans []geometry
  i:=0
  for _,obj := range objects {
    transformed:=obj.transform(P.gentransform())
    Fovangle,distance,direction,size:= transformed.shouldrender()
    if P.drawdist<distance-size {
      continue
    }
    if 0>distance+size {
      continue
    }
    if P.Fov<Fovangle {
      continue
    }
    if 0>direction {
      continue
    }
    objectstrans=append(objectstrans,transformed)
    m[i]=distance
    i=i+1
  }
  n := map[float64][]int{}
  var a []float64
  for k, v := range m {
    n[v] = append(n[v], k)
  }
  for k := range n {
    a = append(a, k)
  }
  sort.Sort(sort.Float64Slice(a))
  for _, k := range a {
    for _, s := range n[k] {
      objectstrans[s].rasterize(C)
    }
  }
  return C
}


func multirender(P player, objects []geometry) canvas  {
  m:=map[int]float64{}
  var objectstrans []geometry
  i:=0
  for _,obj := range objects {
    transformed:=obj.transform(P.gentransform())
    Fovangle,distance,direction,size:= transformed.shouldrender()
    if P.drawdist<distance-size {
      continue
    }
    if 0>distance+size {
      continue
    }
    if P.Fov<Fovangle {
      continue
    }
    if 0>direction {
      continue
    }
    objectstrans=append(objectstrans,transformed)
    m[i]=distance
    i=i+1
  }
  n := map[float64][]int{}
  var a []float64
  var index []int
  for k, v := range m {
    n[v] = append(n[v], k)
  }
  for k := range n {
    a = append(a, k)
  }
  sort.Sort(sort.Float64Slice(a))
  for _, k := range a {
    for _, s := range n[k] {
      index=append(index,s)
    }
  }
  xgrid:=10
  ygrid:=10
  yspacing:=P.yres/ygrid
  if yspacing*ygrid!=P.yres || (P.xres/xgrid)*xgrid!=P.xres {
    fmt.Println("grid error")

  }
  result:=make(chan canvas, xgrid*ygrid)
  xypair0:=make(chan [2]int, xgrid*ygrid)
  xypair1:=make(chan [2]int, xgrid*ygrid)
  tomerge:=make([]canvas,xgrid*ygrid )
  for wid := 0; wid < P.cores; wid++ {
    go renderthread(P,objectstrans,result,index,xgrid,ygrid,xypair0,xypair1)
  }
  for i := 0; i < xgrid; i++ {
    for j := 0; j < ygrid; j++ {
      xypair0 <- [2]int{i,j}
    }
  }
  close(xypair0)
  for i := 0; i < xgrid*ygrid; i++ {
    xy:=<-xypair1
    tomerge[xy[0]+xy[1]*xgrid]=<-result
    tomerge[xy[0]+xy[1]*xgrid].print(string(xy[0])+string(xy[1])+".png")
  }
  pic:=make([][]color.RGBA, P.yres)
  zs:=make([][]float64, P.yres)

  for i := 0; i < ygrid; i++ {
    for isub := 0; isub < yspacing; isub++ {
      pici:=make( []color.RGBA ,0)
      zsi:=make( []float64,0)
      for j := 0; j < xgrid; j++ {
        pici=append(pici,tomerge[j+i*xgrid].pic[isub]...)
        zsi=append(zsi,tomerge[j+i*xgrid].z[isub]...)
      }
      pic[i*yspacing+isub]=pici
      zs[i*yspacing+isub]=zsi
    }
  }
  return canvas{pic ,zs,P.xres,P.yres,-P.xres/2,-P.yres/2,P.res}
}

func renderthread(P player, objects []geometry, result chan<- canvas, order []int, xgrid int, ygrid int, xypair0 <-chan [2]int , xypair1 chan<- [2]int )  {
  for xy := range xypair0 {
    canvas:=newcanvas(P.xres/xgrid,P.yres/ygrid,-P.xres/2+(xy[0]*P.xres)/xgrid,-P.yres/2+(xy[1]*P.yres)/ygrid,P.res,P.drawdist)
    for _,index := range order {
      objects[index].rasterize(canvas)
    }
    result <- canvas
    xypair1 <- xy
  }
}

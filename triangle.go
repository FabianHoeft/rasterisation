package main


import (
	"image/color"
)




type triangle struct{
	p0,p1,p2,n vec3
	p0dotn float64
	size float64
	color color.RGBA
}

func Newtriangle(p0 vec3, p1 vec3, p2 vec3,c color.RGBA) triangle  {
	n:=p1.minus(p0).cross(p2.minus(p1))
	size0:=p0.minus(p1).Abs()
	size1:=p0.minus(p2).Abs()
	if size0<size1{size0=size1}
	return triangle{p0,p1,p2,n,n.dot(p0),size0,c}

}

func (T triangle) transform(M Tmat) geometry {
	p0n:=M.transform(T.p0)
	nn:=M.Rot.multvec(T.n)
	return triangle{p0n,M.transform(T.p1),M.transform(T.p2),nn,p0n.dot(nn),T.size,T.color}
}

// return (x^2+y^2)/z^2, z dist and direction approx
func (T triangle) shouldrender() (float64,float64,float64,float64){
	return (T.p0.x*T.p0.x+T.p0.y*T.p0.y-T.size)/(T.p0.z*T.p0.z),T.p0.z,T.p0dotn/T.p0.Abs()/T.n.Abs(),T.size
}

// res = pixel per unit

func (T triangle) rasterize( C canvas) {
	mask:=make([][]bool,C.yres,C.yres)
	for i := 0; i < C.yres; i++ {
		maskpre:=make([]bool,C.xres,C.xres )
		mask[i]=maskpre
	}
	p0:=[2]float64{T.p0.x/T.p0.z*C.res,T.p0.y/T.p0.z*C.res}
	p1:=[2]float64{T.p1.x/T.p1.z*C.res,T.p1.y/T.p1.z*C.res}
	p2:=[2]float64{T.p2.x/T.p2.z*C.res,T.p2.y/T.p2.z*C.res}
	edges:=[3][2][2]float64{[2][2]float64{p0,p1},[2][2]float64{p1,p2},[2][2]float64{p2,p0}}

	miny,maxy:=2147483647,0
	for _,edge := range edges {
		if edge[0][1]>edge[1][1] {
			edge[0],edge[1]=edge[1],edge[0]
		}
		y0:=int( edge[0][1]-0.5)-C.y0
		y1:=int( edge[1][1]-0.5)-C.y0
		if y0 < 0 {
			y0=0
		} else if y0 >= C.yres {
			y0=C.yres-1
		}
		if y1 <0 {
			y1=0
		} else if y1 >= C.yres {
			y1=C.yres-1
		}
		if y0<miny {
			miny=y0
		}
		if y1>maxy {
			maxy=y1
		}
		for i := y0; i <= y1; i++ {
			floati:=float64(i+C.y0)+0.5
	 		slopeinv:= (edge[0][0]-edge[1][0])/(edge[0][1]-edge[1][1])
			crosx:=edge[0][0]+slopeinv*(floati-edge[0][1])
			xindex:=int(crosx+0.5)-C.x0
			if xindex<C.xres {
				if xindex>=0 {
					mask[i][xindex]= !mask[i][xindex]
				} else{
					mask[i][0]= !mask[i][0]

				}
			}
		}
	}
	for i :=miny ; i < maxy ; i++ {
		inside:=false
		for j := 0; j < C.xres; j++ {
			if mask[i][j]{
				inside= !inside
				if !inside{
					break
				}
			}
			if inside {
				z:=T.p0dotn/T.n.dot(vec3{float64(j+C.x0)/C.res,float64(i+C.y0)/C.res,1.})
				if z<C.z[i][j] {
					C.z[i][j]=z
					C.pic[i][j]=T.color
				}
			}
		}
	}
}

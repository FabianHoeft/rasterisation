package main


import (
	"image/color"
	"math"
)



type polygon struct{
	p []vec3
	edges [][2]int
	n vec3
	p0dotn float64
	size float64
	color color.RGBA
}

func Newpolygon(p []vec3, edges [][2]int, c color.RGBA) polygon  {
	n:=p[1].minus(p[0]).cross(p[2].minus(p[1]))
	size:= 0.
	for i := 1; i < len(edges); i++ {
		temp:=p[0].minus(p[1]).Abs()
		if temp> size {size=temp}
	}
	return polygon{p,edges,n,n.dot(p[0]),size,c}

}

func (P polygon) transform(M Tmat) geometry {
	pnew:=make([]vec3, len(P.edges))
	for i,ps := range P.p {
		pnew[i]=M.transform(ps)
	}
	nn:=M.Rot.multvec(P.n)
	return polygon{pnew,P.edges,nn,nn.dot(pnew[0]),P.size,P.color}
}

// return (x^2+y^2)/z^2, z dist and direction approx
func (P polygon) shouldrender() (float64,float64,float64,float64){
	return (P.p[0].x*P.p[0].x+P.p[0].y*P.p[0].y-P.size)/(P.p[0].z*P.p[0].z),P.p[0].z,P.p0dotn/P.p[0].Abs()/P.n.Abs(),P.size
}

// res = pixel per unit

func (P polygon) rasterize(C canvas) {
	mask:=make([][]bool,C.yres,C.yres)
	for i := 0; i < C.yres; i++ {
		maskpre:=make([]bool,C.xres,C.xres )
		mask[i]=maskpre
	}

	miny,maxy:=2147483647,0
	for _,edind := range P.edges {
		edge:=[2][2]float64{{P.p[edind[0]].x/P.p[edind[0]].z*C.res,P.p[edind[0]].y/P.p[edind[0]].z*C.res},{P.p[edind[1]].x/P.p[edind[1]].z*C.res,P.p[edind[1]].y/P.p[edind[1]].z*C.res}}

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
		if y1 >= C.yres {
			y1=C.yres-1
		}
		if y0<miny {
			miny=y0
		}
		if y1>maxy {
			maxy=y1
		}
		for i := y0; i < y1; i++ {
			floati:=float64(i+C.y0)+0.5
	 		slopeinv:= (edge[0][0]-edge[1][0])/(edge[0][1]-edge[1][1])
			crosx:=edge[0][0]+slopeinv*(floati-edge[0][1])
			x:=math.Floor(crosx+0.5)
			if (x<=crosx-0.5) {
				x=x+1
			}
			xindex:=int(x)-C.x0
			if xindex<C.xres {
				if xindex>=0 {
					mask[i][xindex]= !mask[i][xindex]
				} else{
					mask[i][0]= !mask[i][0]
				}
			}
		}
	}
	for i :=miny ; i <maxy ; i++ {
		inside:=false
		for j := 0; j < C.xres; j++ {
			if mask[i][j]{
				inside= !inside
			}
			if inside {
				z:=P.p0dotn/P.n.dot(vec3{float64(j+C.x0)/C.res,float64(i+C.y0)/C.res,1.})
				if z<C.z[i][j] {
					C.z[i][j]=z
					C.pic[i][j]=P.color
				}
			}
		}
	}
}

package main


import (
	"math"

)

type player struct {
  phi,theta float64
  pos vec3
	Fov float64
	xres,yres int
	res float64
	drawdist float64
	cores, blocksize int
}

func Newplayer(phi float64,theta float64, pos vec3, Fov float64,xres int,yres int, drawdist float64, cores int, blocksize int) player {
	return player{phi,theta,pos,Fov, xres,yres,float64(xres)/math.Tan(Fov*0.5*math.Pi/180.)*0.5,drawdist,cores,blocksize}
}


func (P player) gentransform() Tmat  {
	thetarot:=mat3{vec3{math.Cos(math.Pi/2.-P.theta),0.,-math.Sin(math.Pi/2.-P.theta)},
								vec3{0.,1.,0.},
								vec3{math.Sin(math.Pi/2.-P.theta),0.,math.Cos(math.Pi/2.-P.theta)}}
	phirot:=mat3{vec3{math.Cos(P.phi),-math.Sin(P.phi),0.},
							vec3{math.Sin(P.phi),math.Cos(P.phi),0},
							vec3{0.,0.,1.}}
	screentrans:=mat3{vec3{0,-1.,0},
								vec3{0,0,1},
								vec3{1,0,0}}
	totrot:=screentrans.multmat(thetarot.multmat(phirot))
	return Tmat{totrot,P.pos.inv()}
}

func (P player) gentransform2() Tmat  {
	screentrans:=mat3{vec3{1,0,0},
								vec3{0,1,0},
								vec3{0,0,1}}
	return Tmat{screentrans,P.pos.inv()}
}

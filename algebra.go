package main


import (
	"math"
)


type vec3 struct {
	x,y,z float64
}

func (v0 vec3) minus(v1 vec3) vec3 {
	return vec3{v0.x-v1.x,v0.y-v1.y,v0.z-v1.z}
}

func (v0 vec3) plus(v1 vec3) vec3 {
	return vec3{v0.x+v1.x,v0.y+v1.y,v0.z+v1.z}
}

func (v0 vec3) dot(v1 vec3) float64 {
  return v0.x*v1.x + v0.y*v1.y+ v0.z*v1.z
}

func (v0 vec3) cross(v1 vec3) vec3 {
  return vec3{v0.y*v1.z-v0.z*v1.y, v0.z*v1.x-v0.x*v1.z, v0.x*v1.y-v0.y*v1.x}
}

func  (v vec3) Abs() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y+ v.z*v.z)
}

func (v vec3) normalize() vec3 {
  temp:=1./v.Abs()
  return vec3{v.x*temp,v.y*temp,v.z*temp}
}

func (v vec3) inv() vec3 {
  return vec3{-v.x,-v.y,-v.z}
}



type mat3 struct {
  r0,r1,r2 vec3
}

func ( M mat3)T() mat3{
  return mat3{vec3{M.r0.x,M.r1.x,M.r2.x},vec3{M.r0.y,M.r1.y,M.r2.y},vec3{M.r0.z,M.r1.z,M.r2.z}}
}

func (M mat3) multvec(v vec3) vec3 {
  return vec3{M.r0.dot(v),M.r1.dot(v),M.r2.dot(v)}
}

func (M0 mat3) multmat( M1 mat3) mat3 {
  M1T:=M1.T()
  return mat3{vec3{M0.r0.dot(M1T.r0),M0.r0.dot(M1T.r1),M0.r0.dot(M1T.r2)},
              vec3{M0.r1.dot(M1T.r0),M0.r1.dot(M1T.r1),M0.r1.dot(M1T.r2)},
              vec3{M0.r2.dot(M1T.r0),M0.r2.dot(M1T.r1),M0.r2.dot(M1T.r2)},}
  //return T(mat3{mult(M0,M1T.r0),mult(M0,M1T.r1),mult(M0,M1T.r2)})
}


type Tmat struct{
 	Rot mat3
	trans vec3
}

func (T Tmat) transform( v vec3) vec3{
	return T.Rot.multvec(v).plus(T.trans)
}

func (T0 Tmat) concat(T1 Tmat) Tmat{
	return Tmat{T0.Rot.multmat(T1.Rot),T0.Rot.multvec(T1.trans).plus(T0.trans)}
}

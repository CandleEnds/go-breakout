#version 120

#define M_PI 3.1415926535897932384626433832795

attribute vec3 position;
attribute vec2 texCoord;

uniform vec2 offset;
uniform mat4 VP;

//Radius of output cylinder
uniform float cylinderRadius;
//Height of output cylinder
uniform float cylinderHeight;
//Level width, mapped to cylinder circumference
uniform float levelWidth;
//Level height, mapped to cylinder height
uniform float levelHeight;

varying vec2 TexCoordOut;
varying float normPosOut;

void main()
{
   vec4 vCoord4 = vec4(position, 1.0);
   vec4 vOffset4 = vec4(offset, 0.0, 0.0);
   vec4 levelPosition = vCoord4 + vOffset4;

   float twopi = 2 * M_PI;
   float angleNorm = levelPosition.x / levelWidth;
   float angleRad = angleNorm * twopi;

   float xOut = cylinderRadius * sin(angleRad);
   float yOut = levelPosition.y * levelHeight / cylinderHeight;
   float zOut = cylinderRadius * cos(angleRad);

   gl_Position = VP * vec4(xOut, yOut, zOut, 1.0);
   TexCoordOut = texCoord;
   normPosOut = angleNorm;
}
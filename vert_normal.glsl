#version 120
  
attribute vec3 position;
attribute vec2 texCoord;

uniform vec2 offset;
uniform mat4 VP;

varying vec2 TexCoordOut;

void main()
{
  vec4 vCoord4 = vec4(position, 1.0);
  vec4 vOffset4 = vec4(offset, 0.0, 0.0);
  gl_Position = VP * (vCoord4 + vOffset4);
  TexCoordOut = texCoord;
}
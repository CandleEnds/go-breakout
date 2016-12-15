#version 330

uniform mat4 VP;
uniform vec2 offset;

in vec2 texCoord;
in vec3 position;

out vec2 fragTexCoord;

void main() {
  vec4 vCoord4 = vec4(position, 1.0);
  vec4 vOffset4 = vec4(offset, 0.0, 0.0);
  gl_Position = VP * (vCoord4 + vOffset4);
  fragTexCoord = texCoord;
}

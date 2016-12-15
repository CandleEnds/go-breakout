#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in float normPosOut;

out vec4 outputColor;

void main() {
   vec4 additive = vec4(normPosOut, 0.0, 0.0, 0.0);
   outputColor = texture(tex, fragTexCoord); // + additive
}

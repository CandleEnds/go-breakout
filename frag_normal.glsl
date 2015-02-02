#version 120
varying vec2 TexCoordOut;

uniform sampler2D Sampler;

void main()
{
   //gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
   gl_FragColor = texture2D(Sampler, TexCoordOut);
}
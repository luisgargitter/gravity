#version 330 core

out vec4 outputColor;

void main()
{
    outputColor = vec4(vec3(1/gl_FragCoord.z), 1.0);
}
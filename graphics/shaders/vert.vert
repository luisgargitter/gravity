#version 330 core

uniform mat4 view;

in vec3 vert;

void main() {
    gl_Position = view * vec4(vert, 1.0f);
}
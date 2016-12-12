#version 400 core

in vec2 pass_textureCoords;

out vec4 outColour;

uniform sampler2D textureSampler;

void main(void){
    outColour = texture(textureSampler, pass_textureCoords);
}
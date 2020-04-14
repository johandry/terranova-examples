package data

// DockerComposeTmpl is the Docker Compose code in form of Go template to be
// executed by Docker Compose after rendered.
const DockerComposeTmpl = `
version: '3'
services:
  app:
    image: sdelements/lets-chat:latest
    links:
      - mongo
    ports:
      - ${ letschat_port }:8080
      - 5222:5222

  mongo:
    image: mongo:latest
`

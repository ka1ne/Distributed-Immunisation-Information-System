FROM node:latest

# LABEL maintainer="emailkaine@gmail.com"

# Create app directory
WORKDIR /usr/src/app

#COPY test-network ./

#COPY chaincode ./

COPY package*.json ./

RUN npm install

# Bundle app source
COPY . .

EXPOSE 3000



ENTRYPOINT [ "node", "app.js" ]
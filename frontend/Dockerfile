FROM node:22-alpine

WORKDIR /app

ARG VITE_API_BASE_URL

ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

COPY package.json .

RUN npm i

COPY . .

RUN npm run build

RUN npm i serve -g

CMD [ "serve", "-l", "80", "-s", "dist" ]

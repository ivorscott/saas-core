FROM node:14.18.2

WORKDIR /app

COPY . ./

EXPOSE 3000

RUN npm i

CMD ["npm", "start"]

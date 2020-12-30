FROM node:alpine3.12 as base
ENV NODE_ENV=production
EXPOSE 4000
WORKDIR /application
COPY package*.json ./
RUN npm ci && npm cache clean --force

FROM base as dev
ENV NODE_ENV=development
ENV PATH /application/node_modules/.bin:$PATH
RUN npm i --only=development \
    && npm cache clean --force \
    && npm config ls -l
COPY . .
CMD ["./node_modules/.bin/ts-node-dev", \
    "--rs", "--notify", "false", "--transpile-only", "--ignore-watch", "node_modules", "index.ts"]
FROM dev as test
RUN npm audit

FROM test as build-stage
RUN npm run build

FROM base as prod
COPY --from=build-stage /application/dist /application/dist
CMD ["node", "./dist/index.js"]
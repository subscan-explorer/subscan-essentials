FROM node:10.17.0-jessie as builder

WORKDIR /app
COPY package.json /app
COPY yarn.lock /app
RUN yarn install
ADD . /app
RUN yarn build

FROM nginx:latest
COPY --from=builder /app/dist /usr/share/nginx/html
COPY --from=builder /app/default.conf /etc/nginx/conf.d/default.conf
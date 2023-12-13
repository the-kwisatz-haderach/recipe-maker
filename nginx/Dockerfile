# syntax=docker/dockerfile:1
FROM nginx:1.24.0-alpine-slim

ADD index.html /etc/nginx/html/index.html
COPY nginx.conf /etc/nginx/nginx.conf
COPY nginx.conf /etc/nginx/conf.d/

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-interval=5s \
  CMD service nginx status

CMD ["nginx", "-g", "daemon off;"]